package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/common"
	"github.com/commander-cli/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	checkLocalCmd.Flags().SortFlags = false

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DataConfigDir, "data-config-dir", "",
		"Data config dir, same as nwpc_data_client local command.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DataConfigFile, "data-config-file", "",
		"Data config file path. If set, --data-config-dir and --data-type are ignored.")

	checkLocalCmd.Flags().StringVar(&checkConfigFile, "checker-config", "",
		"Checker runtime config file path. CLI flags override values in this file.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.LocationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	checkLocalCmd.Flags().IntVar(&checkFlagConfig.MaxCheckCount, "max-check-count", 2880,
		"max check count for one forecast time.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.CheckInterval, "check-interval", "5s",
		"check interval, time duration, such as 30s, 1min and so on.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.ExecuteCommand, "execute-command", "",
		"command template to be executed when file is available")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DelayTime, "delay-time", "10s",
		"delay time for each forecast time.")

	checkLocalCmd.Flags().BoolVar(&checkFlagConfig.Debug, "debug", false, "debug mode")
}

const checkLocalCommandName = "local"

const checkLocalCommandDocString = `nwpc_data_client check local
Check local data path using config files in config dir.

Args:
    start_time: YYYYMMDDHH, such as 2018080100`

var checkLocalCmd = &cobra.Command{
	Use:   checkLocalCommandName,
	Short: "Check local data.",
	Long:  checkLocalCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires one arguments")
		}
		var err error
		checkStartTime, err = common.ParseStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check startTime failed: %s", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Load checker runtime config from YAML if provided, otherwise start empty.
		fileConfig := CheckerConfig{}
		if checkConfigFile != "" {
			var err error
			fileConfig, err = LoadCheckerConfig(checkConfigFile)
			if err != nil {
				log.Fatalf("load checker config failed: %v", err)
			}
		}

		// Merge YAML config with CLI flags into a single config object.
		config, err := MergeCheckerConfig(fileConfig, checkFlagConfig, cmd)
		if err != nil {
			log.Fatalf("invalid checker config: %v", err)
		}

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
			config.ForecastTimes = readForecastTimesFromStdin(os.Stdin)
		}
		if len(config.ForecastTimes) == 0 {
			log.Fatalf("no forecast times provided: set forecast_times in checker config or pipe forecast times to stdin")
		}

		forecastTimeList, err := config.ParseForecastTimes()
		if err != nil {
			log.Fatalf("parse forecast times failed: %v", err)
		}
		logForecastTimeList(forecastTimeList)

		// data config dir, data type
		var configContent string
		if config.DataConfigFile != "" {
			configContent, err = common.LoadConfigContentFromFile(config.DataConfigFile)
		} else {
			dataType := config.DataType
			if len(config.DataConfigDir) == 0 {
				dataType = checkLocalCommandName + "/" + dataType
			}
			configContent, err = common.LoadConfigContent(config.DataConfigDir, dataType)
		}
		if err != nil {
			log.Fatalf("load config content failed: %v\n", err)
			return
		}

		checkDataFile(
			config,
			configContent,
			levels,
			checkDuration,
			commandTemplates,
			delayTime,
			forecastTimeList)

		log.Infof("exiting")
	},
}

func checkDataFile(
	config CheckerConfig,
	configContent string,
	levels []string,
	checkDuration time.Duration,
	commandTemplates []*template.Template,
	delayTime time.Duration,
	forecastTimeList []time.Duration) {
	ch := make(chan CheckResult)

	member := ""

	for index, oneTime := range forecastTimeList {
		go func(currentIndex int, forecastTime time.Duration) {
			sleepTime := delayTime * time.Duration(currentIndex)
			forecastTimeString := common.FormatForecastTimeShort(forecastTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("sleeping before check...%v", sleepTime)
			time.Sleep(sleepTime)
			log.WithFields(log.Fields{"forecast_time": forecastTimeString}).
				Infof("checking begin...")
			checkForOneTime(config, ch, configContent, levels, forecastTime, member, checkDuration)
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
			forecastTimeString := common.FormatForecastTimeShort(result.ForecastTime)
			currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})
			if result.Error != nil {
				currentLog.Fatalf("check failed: %v", result.Error)
			} else {
				currentLog.Infof("file is available, run command...")

				if len(commandTemplates) == 0 {
					currentLog.Debugf("ignore execute command because of empty command")
					continue
				}

				err := runCommands(commandTemplates, checkStartTime, result.ForecastTime, member, result.FilePath)
				if err != nil {
					currentLog.Fatalf("run command failed: %v", err)
				} else {
					currentLog.Infof("run command success")
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
	config CheckerConfig,
	ch chan CheckResult,
	configContent string,
	levels []string,
	forecastTime time.Duration,
	member string,
	checkDuration time.Duration) {
	foundData := false
	roundNumber := 0
	forecastTimeString := common.FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecast_time": forecastTimeString})

	dataConfig, err := common.ParseConfigContent(configContent, checkStartTime, forecastTime, member)
	if err != nil {
		currentLog.Fatalf("parse config content failed: %v", err)
		return
	}

	filePath := dataConfig.Default

	for roundNumber < config.MaxCheckCount {
		currentLog.Infof("checking... %d/%d", roundNumber, config.MaxCheckCount)
		filePath = findLocalFileForCheck(dataConfig, levels, forecastTime)
		if filePath == dataConfig.Default {
			currentLog.Infof("checking exist...not found")
		} else {
			currentLog.Infof("checking exist...success: %s", filePath)
			currentLog.Infof("checking size... %d/%d", roundNumber, config.MaxCheckCount)

			var lastSize int64 = -1
			for roundNumber < config.MaxCheckCount {
				currentSize, _ := getFileSize(filePath)
				if currentSize == 0 {
					currentLog.Warnf("get size 0, continue to wait for a valid data size %d/%d", roundNumber, config.MaxCheckCount)
					time.Sleep(checkDuration)
					lastSize = currentSize
				} else if currentSize == lastSize {
					currentLog.Infof("checking size...success %d/%d", roundNumber, config.MaxCheckCount)
					foundData = true
					break
				} else {
					currentLog.Infof("checking size...changed %d/%d", roundNumber, config.MaxCheckCount)
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

func findLocalFileForCheck(config common.DataConfig, levels []string, forecastTime time.Duration) string {
	pathItem := common.FindLocalFile(config, levels, checkStartTime, forecastTime)
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
	common.ConfigTemplateVariable
	FilePath string
}

func runCommand(commandTemplate *template.Template, startTime time.Time, forecastTime time.Duration, member string, filePath string) error {
	tpVar := common.GenerateConfigTemplateVariable(startTime, forecastTime, member)
	var checkerVar CheckerTemplateVariable
	checkerVar.ConfigTemplateVariable = tpVar
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
	}

	return nil
}

func runCommands(commandTemplates []*template.Template, startTime time.Time, forecastTime time.Duration, member string, filePath string) error {
	for _, commandTemplate := range commandTemplates {
		if err := runCommand(commandTemplate, startTime, forecastTime, member, filePath); err != nil {
			return err
		}
	}
	return nil
}

func buildCommandTemplates(command string, commands []string) ([]*template.Template, error) {
	if command != "" && len(commands) > 0 {
		return nil, fmt.Errorf("execute-command and execute-commands are mutually exclusive")
	}

	funcMap := template.FuncMap{
		"generateStartTime":    common.GenerateStartTime,
		"getYear":              common.GetYear,
		"getMonth":             common.GetMonth,
		"getDay":               common.GetDay,
		"getHour":              common.GetHour,
		"generateForecastTime": common.GenerateForecastTime,
		"getForecastHour":      common.GetForecastHour,
		"getForecastMinute":    common.GetForecastMinute,
	}

	var sources []string
	if command != "" {
		sources = []string{command}
	} else {
		sources = commands
	}

	var templates []*template.Template
	for i, s := range sources {
		t, err := template.New(fmt.Sprintf("command-%d", i)).Funcs(funcMap).Delims("{", "}").Parse(s)
		if err != nil {
			return nil, fmt.Errorf("parse execute command failed: %v", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func logForecastTimeList(forecastTimeList []time.Duration) {
	startTimeString := checkStartTime.Format("2006010215")
	for _, forecastTime := range forecastTimeList {
		forecastTimeString := common.FormatForecastTimeShort(forecastTime)
		log.Infof("got check task for %s + %s", startTimeString, forecastTimeString)
	}
}
