package checker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/finder"
	log "github.com/sirupsen/logrus"
)

// batchedTarget is one pending forecast time in batched check mode. The
// config template is expanded once up front (it is deterministic for a given
// start/forecast time), while rounds and observing mutate across rounds.
type batchedTarget struct {
	forecastTime time.Duration
	fileNames    []string
	paths        []dataconfig.PathItem // location-level filtered, config order
	defaultPath  string                // the config's default (not-found placeholder)
	rounds       int                   // existence checks consumed
	observing    bool                  // in size-stability observation
}

// expandBatchedTargets renders the config template for every forecast time
// and applies the location-level filter. A parse failure for one forecast
// time submits an error result for it (fail-fast via the executor) while the
// remaining forecast times proceed.
func expandBatchedTargets(
	configContent string,
	startTime time.Time,
	forecastTimeList []time.Duration,
	member string,
	levels []string,
	exec *commandExecutor,
) []*batchedTarget {
	var targets []*batchedTarget
	for _, forecastTime := range forecastTimeList {
		dataConfig, err := dataconfig.ParseConfigContent(configContent, startTime, forecastTime, member)
		if err != nil {
			exec.submit(CheckResult{
				ForecastTime: forecastTime,
				Error:        fmt.Errorf("parse config content failed: %v", err),
			})
			continue
		}

		var fileNames []string
		fileNames = append(fileNames, dataConfig.FileNames...)
		if len(dataConfig.FileName) > 0 {
			fileNames = append(fileNames, dataConfig.FileName)
		}

		var paths []dataconfig.PathItem
		for _, item := range dataConfig.Paths {
			if !finder.InLocationLevelTypes(levels, item.LevelType) {
				continue
			}
			paths = append(paths, item)
		}

		targets = append(targets, &batchedTarget{
			forecastTime: forecastTime,
			fileNames:    fileNames,
			paths:        paths,
			defaultPath:  dataConfig.Default,
		})
	}
	return targets
}

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
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]struct{}{}, nil
		}
		return nil, err
	}
	set := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		set[entry.Name()] = struct{}{}
	}
	return set, nil
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
// one goroutine per forecast time. Results are submitted to the caller's
// executor; when every forecast time has submitted, the loop closes the queue
// and waits for it to drain.
func checkBatched(
	config CheckerConfig,
	startTime time.Time,
	configContent string,
	levels []string,
	checkDuration time.Duration,
	forecastTimeList []time.Duration,
	exec *commandExecutor,
) {
	member := ""

	targets := expandBatchedTargets(configContent, startTime, forecastTimeList, member, levels, exec)

	pending := make(map[time.Duration]*batchedTarget, len(targets))
	for _, t := range targets {
		pending[t.forecastTime] = t
	}

	maxConcurrent := config.MaxConcurrentStableChecks
	if maxConcurrent < 1 {
		log.Warnf("max_concurrent_stable_checks %d is invalid, falling back to 1", maxConcurrent)
		maxConcurrent = 1
	}
	sem := make(chan struct{}, maxConcurrent)

	var mu sync.Mutex // guards pending and the targets' rounds/observing fields

	for {
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
				exec.submit(CheckResult{
					ForecastTime: t.forecastTime,
					FilePath:     t.defaultPath,
					Error:        fmt.Errorf("too many times"),
				})
				continue
			}
			participants = append(participants, t)
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
			forecastTimeString := dataconfig.FormatForecastTimeShort(t.forecastTime)
			currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})
			currentLog.Infof("checking... %d/%d", t.rounds, config.MaxCheckCount)

			filePath, found := matchTarget(t, dirSets)
			if !found {
				currentLog.Infof("checking exist...not found")
				t.rounds++
				continue
			}
			currentLog.Infof("checking exist...success: %s", filePath)

			// Stability observation runs asynchronously, capped by the
			// semaphore; targets that miss a slot retry next round.
			select {
			case sem <- struct{}{}:
			default:
				currentLog.Infof("stability observation slots full, retry next round")
				continue
			}

			t.observing = true
			budget := config.MaxCheckCount - t.rounds
			go func(t *batchedTarget, filePath string, budget int) {
				defer func() { <-sem }()
				_, err := waitForStableFile(filePath, checkDuration, budget)

				// Submit before removing the target from pending: the main
				// loop closes the queue once pending is empty, so a removal
				// must imply the result is already queued.
				if err != nil {
					currentLog.Warnf("checking size...failed: %v", err)
					exec.submit(CheckResult{
						ForecastTime: t.forecastTime,
						FilePath:     filePath,
						Error:        fmt.Errorf("too many times"),
					})
				} else {
					exec.submit(CheckResult{
						ForecastTime: t.forecastTime,
						FilePath:     filePath,
					})
				}

				mu.Lock()
				delete(pending, t.forecastTime)
				mu.Unlock()
			}(t, filePath, budget)
		}
		mu.Unlock()

		sleepWithJitter(checkDuration, config.CheckIntervalJitter)
	}

	// All forecast times have submitted their results (the queue is buffered,
	// so submissions never block); close the queue and wait until drained.
	exec.finish()
	exec.wait()
}
