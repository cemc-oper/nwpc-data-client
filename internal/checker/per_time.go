package checker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/finder"
	log "github.com/sirupsen/logrus"
)

// Run executes the per-forecast-time polling loop: for every forecast time it
// waits for the data file to appear and become stable, then runs the
// configured command(s). It resolves the data config content itself from the
// checker config (embedded configs, config dir, or a single config file).
//
// It terminates the process via log.Fatal on unrecoverable errors, matching
// the CLI's historical behavior.
func Run(config CheckerConfig, startTime time.Time, stdin io.Reader) {
	// debug
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// location levels
	levels := strings.Split(config.LocationLevels, ",")

	// check duration
	checkDuration, err := time.ParseDuration(config.CheckInterval)
	if err != nil {
		log.Fatalf("parse check-interval failed: %v", err)
	}

	// delay time
	delayTime, err := time.ParseDuration(config.DelayTime)
	if err != nil {
		log.Fatalf("parse delay-time failed: %v", err)
	}

	// command templates
	commandTemplates, err := buildCommandTemplates(config.ExecuteCommand, config.ExecuteCommands)
	if err != nil {
		log.Fatalf("build command templates failed: %v", err)
	}

	// forecast times: YAML config takes precedence, otherwise stdin.
	// The final string list is stored in config.ForecastTimes so the config
	// object remains the single source of truth.
	if len(config.ForecastTimes) == 0 {
		config.ForecastTimes = readForecastTimesFromStdin(stdin)
	}
	if len(config.ForecastTimes) == 0 {
		log.Fatalf("no forecast times provided: set forecast_times in checker config or pipe forecast times to stdin")
	}

	forecastTimeList, err := config.ParseForecastTimes()
	if err != nil {
		log.Fatalf("parse forecast times failed: %v", err)
	}
	logForecastTimeList(startTime, forecastTimeList)

	// data config dir, data type
	var configContent string
	if config.DataConfigFile != "" {
		configContent, err = dataconfig.LoadConfigContentFromFile(config.DataConfigFile)
	} else {
		dataType := config.DataType
		if len(config.DataConfigDir) == 0 {
			dataType = "local/" + dataType
		}
		configContent, err = dataconfig.LoadConfigContent(config.DataConfigDir, dataType)
	}
	if err != nil {
		log.Fatalf("load config content failed: %v\n", err)
		return
	}

	// Startup jitter: shift this process's polling phase once before entering
	// the loop, so processes launched at the same moment do not stay aligned.
	jitterSleep := startupJitter(checkDuration)
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

	if config.CheckMode == CheckModeBatched {
		if config.DelayTime != "" {
			log.Infof("delay-time is ignored in batched check mode")
		}
		exec := newCommandExecutor(len(forecastTimeList), commandTemplates, startTime, "", nil)
		checkBatched(
			config,
			startTime,
			configContent,
			levels,
			checkDuration,
			forecastTimeList,
			exec)
	} else {
		checkDataFile(
			config,
			startTime,
			configContent,
			levels,
			checkDuration,
			commandTemplates,
			delayTime,
			forecastTimeList)
	}

	log.Infof("exiting")
}

func checkDataFile(
	config CheckerConfig,
	startTime time.Time,
	configContent string,
	levels []string,
	checkDuration time.Duration,
	commandTemplates []*template.Template,
	delayTime time.Duration,
	forecastTimeList []time.Duration) {
	member := ""
	exec := newCommandExecutor(len(forecastTimeList), commandTemplates, startTime, member, nil)

	var wg sync.WaitGroup
	for index, oneTime := range forecastTimeList {
		wg.Add(1)
		go func(currentIndex int, forecastTime time.Duration) {
			defer wg.Done()
			sleepTime := delayTime * time.Duration(currentIndex)
			forecastTimeString := dataconfig.FormatForecastTimeShort(forecastTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("sleeping before check...%v", sleepTime)
			time.Sleep(sleepTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("checking begin...")
			checkForOneTime(config, startTime, exec, configContent, levels, forecastTime, member, checkDuration)
		}(index, oneTime)
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

func checkForOneTime(
	config CheckerConfig,
	startTime time.Time,
	exec *commandExecutor,
	configContent string,
	levels []string,
	forecastTime time.Duration,
	member string,
	checkDuration time.Duration) {
	foundData := false
	roundNumber := 0
	forecastTimeString := dataconfig.FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})

	dataConfig, err := dataconfig.ParseConfigContent(configContent, startTime, forecastTime, member)
	if err != nil {
		exec.submit(CheckResult{
			ForecastTime: forecastTime,
			Error:        fmt.Errorf("parse config content failed: %v", err),
		})
		return
	}

	filePath := dataConfig.Default

	for roundNumber < config.MaxCheckCount {
		currentLog.Infof("checking... %d/%d", roundNumber, config.MaxCheckCount)
		filePath = findLocalFileForCheck(dataConfig, startTime, levels, forecastTime)
		if filePath == dataConfig.Default {
			currentLog.Infof("checking exist...not found")
		} else {
			currentLog.Infof("checking exist...success: %s", filePath)
			// The size checks share the same budget of remaining checks.
			_, err := waitForStableFile(filePath, checkDuration, config.MaxCheckCount-roundNumber)
			if err != nil {
				currentLog.Warnf("checking size...failed: %v", err)
			} else {
				foundData = true
			}
			break
		}
		sleepWithJitter(checkDuration, config.CheckIntervalJitter)
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
