package checker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/finder"
	log "github.com/sirupsen/logrus"
)

// Run executes the configured polling loop: for every forecast time it waits
// for the data file to appear and become stable, then runs the configured
// command(s). Real polling and DryRun share buildRunPlan. Cancellation stops
// all checker-owned waits and producers; errors are returned to the caller.
func Run(ctx context.Context, config CheckerConfig, startTime time.Time, stdin io.Reader) (result *RunResult, err error) {
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	plan, err := buildRunPlan(config, startTime, stdin)
	if err != nil {
		return nil, fmt.Errorf("build run plan failed: %w", err)
	}
	plan.tracker = newRunTracker(plan)
	defer func() {
		result = plan.tracker.finish(err)
	}()
	logForecastTimeList(startTime, plan.forecastTimeList)

	// Startup jitter: shift this process's polling phase once before entering
	// the loop, so processes launched at the same moment do not stay aligned.
	jitterSleep := startupJitter(plan.checkInterval)
	log.Infof("startup jitter...sleep %v", jitterSleep)
	if err := sleepContext(ctx, jitterSleep); err != nil {
		return nil, err
	}

	exec := newCommandExecutor(ctx, plan.commandTimeout, len(plan.targets), plan.commandTemplates, startTime, "", nil, plan.tracker)
	if plan.config.CheckMode == CheckModeBatched {
		if plan.config.DelayTime != "" {
			log.Infof("delay-time is ignored in batched check mode")
		}
		err = checkBatched(exec.context(), plan, exec)
	} else {
		err = checkDataFile(exec.context(), plan, exec)
	}
	if err != nil {
		return nil, err
	}

	log.Infof("exiting")
	return nil, nil
}

func checkDataFile(ctx context.Context, plan *runPlan, exec *commandExecutor) error {
	if plan.tracker == nil {
		plan.tracker = newRunTracker(plan)
	}
	var wg sync.WaitGroup
	for index, target := range plan.targets {
		if ctx.Err() != nil {
			break
		}
		wg.Add(1)
		go func(currentIndex int, target *checkTarget) {
			defer wg.Done()
			sleepTime := plan.delayTime * time.Duration(currentIndex)
			forecastTimeString := dataconfig.FormatForecastTimeShort(target.forecastTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("sleeping before check...%v", sleepTime)
			if err := sleepContext(ctx, sleepTime); err != nil {
				return
			}
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("checking begin...")
			checkForOneTime(ctx, plan, exec, target)
		}(index, target)
	}

	// Close only after every producer has returned, then wait until the
	// consumer has drained the queue or stopped on the first error.
	wg.Wait()
	exec.finish()
	return exec.wait()
}

func checkForOneTime(ctx context.Context, plan *runPlan, exec *commandExecutor, target *checkTarget) {
	config := plan.config
	forecastTime := target.forecastTime
	dataConfig := target.dataConfig
	foundData := false
	roundNumber := 0
	checkCount := 0
	forecastTimeString := dataconfig.FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})
	filePath := dataConfig.Default
	plan.tracker.start(target.index)

	for roundNumber < config.MaxCheckCount {
		if ctx.Err() != nil {
			return
		}
		currentLog.Infof("checking... %d/%d", roundNumber, config.MaxCheckCount)
		checkCount++
		filePath = findLocalFileForCheck(dataConfig, plan.startTime, plan.locationLevels, forecastTime)
		plan.tracker.progress(target.index, filePath, checkCount)
		if filePath == dataConfig.Default {
			currentLog.Infof("checking exist...not found")
		} else {
			currentLog.Infof("checking exist...success: %s", filePath)
			// The size checks share the same budget of remaining checks.
			_, stableChecks, err := waitForStableFile(ctx, filePath, plan.checkInterval, config.MaxCheckCount-roundNumber)
			checkCount += stableChecks
			plan.tracker.progress(target.index, filePath, checkCount)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				currentLog.Warnf("checking size...failed: %v", err)
			} else {
				foundData = true
			}
			break
		}
		if err := sleepWithJitter(ctx, plan.checkInterval, config.CheckIntervalJitter); err != nil {
			return
		}
		roundNumber += 1
	}

	result := CheckResult{
		ForecastTime: forecastTime,
		FilePath:     filePath,
		Status:       CheckStatusWaiting,
		CheckCount:   checkCount,
		targetIndex:  target.index,
	}

	if !foundData {
		result.Error = fmt.Errorf("too many times")
	}
	plan.tracker.completeCheck(result)
	exec.submit(result)
}

// readForecastTimesFromStdin reads whitespace-separated forecast time tokens
// from r (usually os.Stdin) and returns them as strings, such as "000h" or
// "003h10m". The caller is responsible for parsing and logging.
func readForecastTimesFromStdin(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var forecastTimeStrings []string
	for scanner.Scan() {
		forecastTimeStrings = append(forecastTimeStrings, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return forecastTimeStrings, nil
}

func findLocalFileForCheck(config dataconfig.DataConfig, startTime time.Time, levels []string, forecastTime time.Duration) string {
	pathItem := finder.FindLocalFile(config, levels, startTime, forecastTime)
	return pathItem.Path
}

func logForecastTimeList(startTime time.Time, forecastTimeList []time.Duration) {
	startTimeString := startTime.Format("2006010215")
	for _, forecastTime := range forecastTimeList {
		forecastTimeString := dataconfig.FormatForecastTimeShort(forecastTime)
		log.Infof("got check task for %s + %s", startTimeString, forecastTimeString)
	}
}
