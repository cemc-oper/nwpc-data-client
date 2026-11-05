package checker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
)

// collectDirectories returns the union of directories referenced by the
// targets, deduplicated. Templates containing {{.ForecastHour}} place
// different forecast times in different directories, so grouping happens on
// the expanded directory strings.
func collectDirectories(targets []*batchedTarget) []string {
	seen := make(map[string]struct{})
	var dirs []string
	for _, t := range targets {
		for _, item := range t.paths {
			if _, ok := seen[item.Path]; ok {
				continue
			}
			seen[item.Path] = struct{}{}
			dirs = append(dirs, item.Path)
		}
	}
	return dirs
}

// readDirSet returns the set of file names currently in dir. A missing
// directory is normal (model output directories are created late), so ENOENT
// yields an empty set. Any other error is reported so the caller can skip the
// directory for this round.
func readDirSet(dir string) (map[string]struct{}, error) {
	// Open + handle-based Stat instead of os.ReadDir: on Windows ReadDir on a
	// regular file silently succeeds, while callers rely on ENOTDIR. The stat
	// rides on the open handle, so it costs no extra path lookup.
	f, err := os.Open(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]struct{}{}, nil
		}
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, &os.PathError{Op: "stat", Path: dir, Err: err}
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "readdir", Path: dir, Err: syscall.ENOTDIR}
	}
	entries, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		set[entry.Name()] = struct{}{}
	}
	return set, nil
}

// gateParticipants enforces sequential arrival gating: while a previous
// forecast time is being observed or its command is still running (gateHeld),
// nothing participates; otherwise only the earliest pending forecast time
// does.
func gateParticipants(participants []*batchedTarget, gateHeld bool) []*batchedTarget {
	if gateHeld || len(participants) == 0 {
		return nil
	}
	earliest := participants[0]
	for _, t := range participants[1:] {
		if t.forecastTime < earliest.forecastTime {
			earliest = t
		}
	}
	return []*batchedTarget{earliest}
}

// matchTarget walks the target's paths in config order and returns the first
// file present in this round's directory sets, preserving FindLocalFile's
// "first existing path wins" priority (runtime beats archive). Directories
// absent from dirSets (skipped after a read error) never match.
func matchTarget(t *batchedTarget, dirSets map[string]map[string]struct{}) (string, bool) {
	for _, item := range t.paths {
		set, ok := dirSets[item.Path]
		if !ok {
			continue
		}
		for _, name := range t.fileNames {
			if _, found := set[name]; found {
				return filepath.Join(item.Path, name), true
			}
		}
	}
	return "", false
}

// checkBatched is the batched-mode polling loop: a single loop checks all
// pending forecast times with one readdir per directory per round, instead of
// one goroutine per forecast time. Stability observers are tracked so failure
// and cancellation cannot leave submissions running after the queue closes.
func checkBatched(ctx context.Context, plan *runPlan, exec *commandExecutor) error {
	if plan.tracker == nil {
		plan.tracker = newRunTracker(plan)
	}
	config := plan.config
	targets := plan.targets

	var mu sync.Mutex // guards pending, gateHeld, and the targets' rounds/observing fields

	pending := make(map[time.Duration]*batchedTarget, len(targets))
	for _, t := range targets {
		pending[t.forecastTime] = t
	}

	// Sequential arrival gating: only the earliest pending forecast time
	// takes part in each round, and the next one is released only after its
	// command has finished (via the executor's onTaskDone hook).
	gated := false
	if len(targets) > 0 {
		gated = resolveArrivalOrder(config.ArrivalOrder, targets[0].arrivalOrder) == dataconfig.ArrivalOrderSequential
	}
	var gateHeld bool
	if gated {
		log.Infof("arrival order is sequential: gating to the earliest pending forecast time")
		exec.setOnTaskDone(func(CheckResult) {
			mu.Lock()
			gateHeld = false
			mu.Unlock()
		})
	}

	maxConcurrent := config.MaxConcurrentStableChecks
	if maxConcurrent < 1 {
		log.Warnf("max_concurrent_stable_checks %d is invalid, falling back to 1", maxConcurrent)
		maxConcurrent = 1
	}
	sem := make(chan struct{}, maxConcurrent)
	var observations sync.WaitGroup

polling:
	for {
		if ctx.Err() != nil {
			break
		}
		mu.Lock()
		if len(pending) == 0 {
			mu.Unlock()
			break
		}

		// Participants: pending targets that still have check budget and are
		// not in stability observation. Iteration follows the original
		// forecast time order for deterministic behavior.
		var participants []*batchedTarget
		for _, t := range targets {
			if _, ok := pending[t.forecastTime]; !ok || t.observing {
				continue
			}
			if t.rounds >= config.MaxCheckCount {
				delete(pending, t.forecastTime)
				result := CheckResult{
					ForecastTime: t.forecastTime,
					FilePath:     t.defaultPath,
					Error:        fmt.Errorf("too many times"),
					Status:       CheckStatusWaiting,
					CheckCount:   t.checks,
					targetIndex:  t.index,
				}
				plan.tracker.completeCheck(result)
				exec.submit(result)
				continue
			}
			participants = append(participants, t)
		}

		// Sequential gating: only the earliest pending forecast time
		// participates, and only while no previous one is being observed or
		// executing.
		if gated {
			participants = gateParticipants(participants, gateHeld)
		}

		// Discovery phase: one readdir per directory per round.
		dirSets := make(map[string]map[string]struct{})
		for _, dir := range collectDirectories(participants) {
			set, err := readDirSet(dir)
			if err != nil {
				log.Warnf("read directory %s failed, skip this round: %v", dir, err)
				continue
			}
			dirSets[dir] = set
		}

		// Match phase: in-memory priority match per target.
		for _, t := range participants {
			plan.tracker.start(t.index)
			forecastTimeString := dataconfig.FormatForecastTimeShort(t.forecastTime)
			currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})
			currentLog.Infof("checking... %d/%d", t.rounds, config.MaxCheckCount)

			t.checks++
			plan.tracker.progress(t.index, t.defaultPath, t.checks)
			filePath, found := matchTarget(t, dirSets)
			if !found {
				currentLog.Infof("checking exist...not found")
				t.rounds++
				continue
			}
			plan.tracker.progress(t.index, filePath, t.checks)
			currentLog.Infof("checking exist...success: %s", filePath)

			// Stability observation runs asynchronously, capped by the
			// semaphore; targets that miss a slot retry next round.
			select {
			case sem <- struct{}{}:
			default:
				currentLog.Infof("stability observation slots full, retry next round")
				continue
			}

			if ctx.Err() != nil {
				<-sem
				mu.Unlock()
				break polling
			}
			t.observing = true
			if gated {
				// The gate is released by the executor's onTaskDone hook once
				// this target's command has finished.
				gateHeld = true
			}
			budget := config.MaxCheckCount - t.rounds
			baseChecks := t.checks
			observations.Add(1)
			go func(t *batchedTarget, filePath string, budget int, baseChecks int) {
				defer observations.Done()
				defer func() { <-sem }()
				_, stableChecks, err := waitForStableFile(ctx, filePath, plan.checkInterval, budget)
				checkCount := baseChecks + stableChecks
				plan.tracker.progress(t.index, filePath, checkCount)
				if ctx.Err() != nil {
					return
				}

				result := CheckResult{
					ForecastTime: t.forecastTime,
					FilePath:     filePath,
					Status:       CheckStatusWaiting,
					CheckCount:   checkCount,
					targetIndex:  t.index,
				}
				if err != nil {
					currentLog.Warnf("checking size...failed: %v", err)
					result.Error = fmt.Errorf("too many times")
				}
				plan.tracker.completeCheck(result)

				// Submit before removing the target from pending: the main
				// loop closes the queue once pending is empty, so a removal
				// must imply the result is already queued.
				exec.submit(result)

				mu.Lock()
				delete(pending, t.forecastTime)
				mu.Unlock()
			}(t, filePath, budget, baseChecks)
		}
		mu.Unlock()

		if err := sleepWithJitter(ctx, plan.checkInterval, config.CheckIntervalJitter); err != nil {
			break
		}
	}

	// No producer remains before the queue is closed, including stability
	// observers that were active when failure or cancellation stopped the loop.
	observations.Wait()
	exec.finish()
	return exec.wait()
}
