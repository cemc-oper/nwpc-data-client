package checker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/finder"
	log "github.com/sirupsen/logrus"
)

// Run executes the per-forecast-time polling loop: for every forecast time it
// waits for the data file to appear and become stable, then runs the
// configured command(s). Real polling and DryRun share buildRunPlan.
//
// It terminates the process via log.Fatal on unrecoverable errors, matching
// the CLI's historical behavior.
func Run(config CheckerConfig, startTime time.Time, stdin io.Reader) {
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	plan, err := buildRunPlan(config, startTime, stdin)
	if err != nil {
		log.Fatalf("build run plan failed: %v", err)
	}
	logForecastTimeList(startTime, plan.forecastTimeList)

	// Startup jitter: shift this process's polling phase once before entering
	// the loop, so processes launched at the same moment do not stay aligned.
	jitterSleep := startupJitter(plan.checkInterval)
	log.Infof("startup jitter...sleep %v", jitterSleep)
	time.Sleep(jitterSleep)

	// Exit code 3 on termination signals, matching historical behavior.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		log.Infof("catching signal: %v\n", sig)
		os.Exit(3)
	}()

	if plan.config.CheckMode == CheckModeBatched {
		if plan.config.DelayTime != "" {
			log.Infof("delay-time is ignored in batched check mode")
		}
		exec := newCommandExecutor(len(plan.targets), plan.commandTemplates, startTime, "", nil)
		checkBatched(plan, exec)
	} else {
		checkDataFile(plan)
	}

	log.Infof("exiting")
}

func checkDataFile(plan *runPlan) {
	member := ""
	exec := newCommandExecutor(len(plan.targets), plan.commandTemplates, plan.startTime, member, nil)

	var wg sync.WaitGroup
	for index, target := range plan.targets {
		wg.Add(1)
		go func(currentIndex int, target *checkTarget) {
			defer wg.Done()
			sleepTime := plan.delayTime * time.Duration(currentIndex)
			forecastTimeString := dataconfig.FormatForecastTimeShort(target.forecastTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("sleeping before check...%v", sleepTime)
			time.Sleep(sleepTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("checking begin...")
			checkForOneTime(plan, exec, target)
		}(index, target)
	}

	// All producers have submitted (the queue is buffered, so submissions
	// never block); close the queue and wait until it is drained.
	wg.Wait()
	exec.finish()
	exec.wait()
}

type CheckResult struct {
	ForecastTime time.Duration
	FilePath     string
	Error        error
}

func checkForOneTime(plan *runPlan, exec *commandExecutor, target *checkTarget) {
	config := plan.config
	forecastTime := target.forecastTime
	dataConfig := target.dataConfig
	foundData := false
	roundNumber := 0
	forecastTimeString := dataconfig.FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})
	filePath := dataConfig.Default

	for roundNumber < config.MaxCheckCount {
		currentLog.Infof("checking... %d/%d", roundNumber, config.MaxCheckCount)
		filePath = findLocalFileForCheck(dataConfig, plan.startTime, plan.locationLevels, forecastTime)
		if filePath == dataConfig.Default {
			currentLog.Infof("checking exist...not found")
		} else {
			currentLog.Infof("checking exist...success: %s", filePath)
			// The size checks share the same budget of remaining checks.
			_, err := waitForStableFile(filePath, plan.checkInterval, config.MaxCheckCount-roundNumber)
			if err != nil {
				currentLog.Warnf("checking size...failed: %v", err)
			} else {
				foundData = true
			}
			break
		}
		sleepWithJitter(plan.checkInterval, config.CheckIntervalJitter)
		roundNumber += 1
	}

	result := CheckResult{
		ForecastTime: forecastTime,
		FilePath:     filePath,
		Error:        nil,
	}

	if !foundData {
		result.Error = fmt.Errorf("too many times")
	}
	exec.submit(result)
}

// readForecastTimesFromStdin reads whitespace-separated forecast time tokens
// from r (usually os.Stdin) and returns them as strings, such as "000h" or
// "003h10m". The caller is responsible for parsing and logging.
func readForecastTimesFromStdin(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var forecastTimeStrings []string
	for scanner.Scan() {
		forecastTimeStrings = append(forecastTimeStrings, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return forecastTimeStrings
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
