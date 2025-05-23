package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/SimonBaeumer/cmd"
	"github.com/cemc-oper/nwpc-data-client/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"
)

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.Flags().SortFlags = false

	localCmd.Flags().StringVar(&configDir, "data-config-dir", "",
		"Data config dir, same as nwpc_data_client local command.")

	localCmd.Flags().StringVar(&dataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	localCmd.Flags().StringVar(&locationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	localCmd.Flags().IntVar(&maxCheckCount, "max-check-count", 2880,
		"max check count for one forecast time.")

	localCmd.Flags().StringVar(&checkInterval, "check-interval", "5s",
		"check interval, time duration, such as 30s, 1min and so on.")

	localCmd.Flags().StringVar(&executeCommand, "execute-command", "",
		"command template to be executed when file is available")

	localCmd.Flags().StringVar(&delayTimeForEachForecastTime, "delay-time", "10s",
		"delay time for each forecast time.")

	localCmd.Flags().BoolVar(&debugMode, "debug", false, "debug mode")

}

const localCommandName = "local"

const localCommandDocString = `nwpc_data_checker local
Check local data path using config files in config dir.

Args:
    start_time: YYYYMMDDHH, such as 2018080100`

var localCmd = &cobra.Command{
	Use:   localCommandName,
	Short: "Check local data.",
	Long:  localCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires one arguments")
		}
		var err error
		startTime, err = common.ParseStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check startTime failed: %s", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// parse options

		// debug
		if debugMode {
			log.SetLevel(log.DebugLevel)
		}

		// location levels
		levels := strings.Split(locationLevels, ",")

		// check duration
		checkDuration, err := time.ParseDuration(checkInterval)
		if err != nil {
			log.Fatalf("parse check-interval failed: %v", err)
		}

		// execute command
		var commandTemplate *template.Template = nil
		if executeCommand != "" {
			commandTemplate = template.Must(template.New("command").Funcs(template.FuncMap{
				"generateStartTime":    common.GenerateStartTime,
				"getYear":              common.GetYear,
				"getMonth":             common.GetMonth,
				"getDay":               common.GetDay,
				"getHour":              common.GetHour,
				"generateForecastTime": common.GenerateForecastTime,
				"getForecastHour":      common.GetForecastHour,
				"getForecastMinute":    common.GetForecastMinute,
			}).Delims("{", "}").Parse(executeCommand))
		}

		// delay time
		delayTime, err := time.ParseDuration(delayTimeForEachForecastTime)
		if err != nil {
			log.Fatalf("parse delay-time failed: %v", err)
		}

		// data config dir, data type
		if len(configDir) == 0 {
			dataType = localCommandName + "/" + dataType
		}
		config, err := common.LoadConfig(configDir, dataType)
		if err != nil {
			log.Fatalf("load config failed: %v\n", err)
			return
		}
		fmt.Printf("%v\n", config)

		checkDataFile(
			config,
			levels,
			checkDuration,
			commandTemplate,
			delayTime)

		log.Infof("exiting")
	},
}

func checkDataFile(
	config common.DataConfig,
	levels []string,
	checkDuration time.Duration,
	commandTemplate *template.Template,
	delayTime time.Duration) {
	ch := make(chan CheckResult)

	forecastTimeList := parseForecastTimeInput(os.Stdin)
	for index, oneTime := range forecastTimeList {
		go func(currentIndex int, forecastTime time.Duration) {
			sleepTime := delayTime * time.Duration(currentIndex)
			forecastTimeString := common.FormatForecastTimeShort(forecastTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("sleeping before check...%v", sleepTime)
			time.Sleep(sleepTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("checking begin...")
			checkForOneTime(ch, config, levels, forecastTime, checkDuration)
		}(index, oneTime)
	}

	done := make(chan bool, 1)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigs
		log.Infof("catching signal: %v\n", sig)
		os.Exit(3)
	}()

	go func() {
		for _ = range forecastTimeList {
			result := <-ch
			if result.Error != nil {
				log.WithFields(log.Fields{
					"forecast_time": fmt.Sprintf("%03dh%02dm", int(result.ForecastTime.Hours()), int(result.ForecastTime.Hours())),
				}).Fatalf("check failed: %v", result.Error)
			} else {
				log.WithFields(log.Fields{
					"forecast_time": fmt.Sprintf("%03dh%02dm", int(result.ForecastTime.Hours()), int(result.ForecastTime.Hours())),
				}).Infof("file is available, run command...")

				if executeCommand == "" {
					continue
				}

				err := runCommand(commandTemplate, startTime, result.ForecastTime, result.FilePath)
				if err != nil {
					log.WithFields(log.Fields{
						"forecast_time": fmt.Sprintf("%03dh%02dm", int(result.ForecastTime.Hours()), int(result.ForecastTime.Hours())),
					}).Fatalf("run command failed: %v", err)
				} else {
					log.WithFields(log.Fields{
						"forecast_time": fmt.Sprintf("%03dh%02dm", int(result.ForecastTime.Hours()), int(result.ForecastTime.Hours())),
					}).Infof("run command success")
				}
			}
		}
		done <- true
	}()

	<-done
}

type CheckResult struct {
	ForecastTime time.Duration
	FilePath     string
	Error        error
}

func checkForOneTime(
	ch chan CheckResult,
	config common.DataConfig,
	levels []string,
	forecastTime time.Duration,
	checkDuration time.Duration) {
	foundData := false
	roundNumber := 0
	filePath := config.Default

	forecastTimeString := common.FormatForecastTimeShort(forecastTime)

	currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})

	for roundNumber < maxCheckCount {
		currentLog.Infof("checking... %d/%d", roundNumber, maxCheckCount)
		filePath = findLocalFile(config, levels, forecastTime)
		if filePath == config.Default {
			currentLog.Infof("checking exist...not found")
		} else {
			currentLog.Infof("checking exist...success: %s", filePath)
			currentLog.Infof("checking size... %d/%d", roundNumber, maxCheckCount)

			var lastSize int64 = -1
			for roundNumber < maxCheckCount {
				currentSize, _ := getFileSize(filePath)
				if currentSize == lastSize {
					currentLog.Infof("checking size...success %d/%d", roundNumber, maxCheckCount)
					foundData = true
					break
				} else {
					currentLog.Infof("checking size...changed %d/%d", roundNumber, maxCheckCount)
					time.Sleep(checkDuration)
					lastSize = currentSize
				}
				roundNumber += 1
			}
			break
		}
		time.Sleep(checkDuration)
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
	ch <- result
}

// parse input string to a forecast time list.
// input string is a list of forecast time, each line or each token is a forecast time,
// the format is "000h00m" or "000h", such as "000h 000h10m 001h00m 001h10m"
// return a list of forecast time.
func parseForecastTimeInput(r io.Reader) []time.Duration {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	startTimeString := startTime.Format("2006010215")
	var forecastTimeList []time.Duration
	for scanner.Scan() {
		forecastTimeString := scanner.Text()
		forecastTime, err := time.ParseDuration(forecastTimeString)
		if err != nil {
			log.Fatalf("parse input has error: %v", err)
		}
		forecastTimeList = append(forecastTimeList, forecastTime)

		forecastTimeStringForLog := common.FormatForecastTimeShort(forecastTime)
		log.Infof("got check task for %s + %s", startTimeString, forecastTimeStringForLog)
	}
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return forecastTimeList
}

func findLocalFile(config common.DataConfig, levels []string, forecastTime time.Duration) string {
	pathItem := common.FindLocalFile(config, levels, startTime, forecastTime)
	return pathItem.Path
}

func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return -100, fmt.Errorf("get file info has error:%v", err)
	}
	return fileInfo.Size(), nil
}

type CheckerTemplateVariable struct {
	common.TimeTemplateVariable
	FilePath string
}

func runCommand(commandTemplate *template.Template, startTime time.Time, forecastTime time.Duration, filePath string) error {
	tpVar := common.GenerateTimeTemplateVariable(startTime, forecastTime)
	var checkerVar CheckerTemplateVariable
	checkerVar.TimeTemplateVariable = tpVar
	checkerVar.FilePath = filePath

	var commandBuilder strings.Builder
	err := commandTemplate.Execute(&commandBuilder, checkerVar)
	if err != nil {
		return fmt.Errorf("command template execute has error: %v", err)
	}

	commandString := commandBuilder.String()

	forecastTimeString := common.FormatForecastTimeShort(forecastTime)
	log.WithFields(log.Fields{
		"forecast_time": forecastTimeString,
	}).Infof("running command <%s> ...", commandString)

	c := cmd.NewCommand(
		commandString,
		cmd.WithStandardStreams,
		cmd.WithInheritedEnvironment(nil))
	err = c.Execute()
	if err != nil {
		return fmt.Errorf("run command <%s> has error: %v", commandString, err)
	} else if c.ExitCode() != 0 {
		return fmt.Errorf("run command <%s> exit code is not 0: %d", commandString, c.ExitCode())
	} else {
		return nil
	}
}
