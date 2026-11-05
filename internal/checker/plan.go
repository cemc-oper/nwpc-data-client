package checker

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/finder"
)

// checkTarget is the expanded, location-filtered plan for one forecast time.
// candidatePaths is informational; polling uses dataConfig and the filtered
// fields so real runs and dry-run share exactly the same expansion.
type checkTarget struct {
	forecastTime   time.Duration
	dataConfig     dataconfig.DataConfig
	fileNames      []string
	paths          []dataconfig.PathItem
	candidatePaths []string
	defaultPath    string
	arrivalOrder   string

	// Batched polling state. These fields are unused during plan construction.
	rounds    int
	observing bool
}

// batchedTarget retains the polling implementation's existing terminology.
type batchedTarget = checkTarget

// runPlan contains all deterministic work that can be completed before any
// polling, sleeping, filesystem discovery, or command execution begins.
type runPlan struct {
	config           CheckerConfig
	startTime        time.Time
	locationLevels   []string
	checkInterval    time.Duration
	delayTime        time.Duration
	forecastTimeList []time.Duration
	targets          []*checkTarget
	commandTemplates []*template.Template
	arrivalOrder     string
}

// buildRunPlan is shared by real runs and dry-run. It reads checker/data
// configs and optionally forecast times from stdin, expands every forecast
// time, and validates command templates with a dummy file path. It never stats
// or reads candidate data paths and never starts polling or an executor.
func buildRunPlan(config CheckerConfig, startTime time.Time, stdin io.Reader) (*runPlan, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid checker config: %w", err)
	}

	checkInterval, err := time.ParseDuration(config.CheckInterval)
	if err != nil {
		return nil, fmt.Errorf("parse check_interval failed: %w", err)
	}
	delayTime, err := time.ParseDuration(config.DelayTime)
	if err != nil {
		return nil, fmt.Errorf("parse delay_time failed: %w", err)
	}
	commandTemplates, err := buildCommandTemplates(config.ExecuteCommand, config.ExecuteCommands)
	if err != nil {
		return nil, fmt.Errorf("build command templates failed: %w", err)
	}

	// YAML forecast times take precedence over stdin, matching check local.
	if len(config.ForecastTimes) == 0 {
		forecastTimes, err := readForecastTimesFromStdin(stdin)
		if err != nil {
			return nil, fmt.Errorf("read forecast times from stdin failed: %w", err)
		}
		config.ForecastTimes = forecastTimes
	}
	if len(config.ForecastTimes) == 0 {
		return nil, fmt.Errorf("no forecast times provided: set forecast_times in checker config or pipe forecast times to stdin")
	}
	forecastTimeList, err := config.ParseForecastTimes()
	if err != nil {
		return nil, fmt.Errorf("parse forecast times failed: %w", err)
	}

	configContent, err := loadCheckerDataConfig(config)
	if err != nil {
		return nil, fmt.Errorf("load config content failed: %w", err)
	}
	locationLevels := strings.Split(config.LocationLevels, ",")
	targets, err := expandCheckTargets(configContent, startTime, forecastTimeList, locationLevels)
	if err != nil {
		return nil, err
	}

	for _, target := range targets {
		if err := validateCommandTemplates(commandTemplates, startTime, target.forecastTime); err != nil {
			return nil, err
		}
	}

	arrivalOrder := dataconfig.ArrivalOrderUnordered
	if len(targets) > 0 {
		arrivalOrder = resolveArrivalOrder(config.ArrivalOrder, targets[0].arrivalOrder)
	}
	return &runPlan{
		config:           config,
		startTime:        startTime,
		locationLevels:   locationLevels,
		checkInterval:    checkInterval,
		delayTime:        delayTime,
		forecastTimeList: forecastTimeList,
		targets:          targets,
		commandTemplates: commandTemplates,
		arrivalOrder:     arrivalOrder,
	}, nil
}

func loadCheckerDataConfig(config CheckerConfig) (string, error) {
	if config.DataConfigFile != "" {
		return dataconfig.LoadConfigContentFromFile(config.DataConfigFile)
	}
	dataType := config.DataType
	if config.DataConfigDir == "" {
		dataType = "local/" + dataType
	}
	return dataconfig.LoadConfigContent(config.DataConfigDir, dataType)
}

// expandCheckTargets renders every forecast time and applies location-level
// filtering in the same order used by both polling modes.
func expandCheckTargets(
	configContent string,
	startTime time.Time,
	forecastTimeList []time.Duration,
	levels []string,
) ([]*checkTarget, error) {
	var targets []*checkTarget
	for _, forecastTime := range forecastTimeList {
		dataConfig, err := dataconfig.ParseConfigContent(configContent, startTime, forecastTime, "")
		if err != nil {
			return nil, fmt.Errorf("parse config content for forecast time %s failed: %w",
				dataconfig.FormatForecastTimeShort(forecastTime), err)
		}

		var fileNames []string
		fileNames = append(fileNames, dataConfig.FileNames...)
		if dataConfig.FileName != "" {
			fileNames = append(fileNames, dataConfig.FileName)
		}

		var paths []dataconfig.PathItem
		var candidatePaths []string
		for _, item := range dataConfig.Paths {
			if !finder.InLocationLevelTypes(levels, item.LevelType) {
				continue
			}
			paths = append(paths, item)
			for _, fileName := range fileNames {
				candidatePaths = append(candidatePaths, filepath.Join(item.Path, fileName))
			}
		}

		targets = append(targets, &checkTarget{
			forecastTime:   forecastTime,
			dataConfig:     dataConfig,
			fileNames:      fileNames,
			paths:          paths,
			candidatePaths: candidatePaths,
			defaultPath:    dataConfig.Default,
			arrivalOrder:   dataConfig.ArrivalOrder,
		})
	}
	return targets, nil
}

// validateCommandTemplates executes the Go templates with a dummy path to
// catch unknown variables and template function errors. It does not run the
// rendered strings as shell commands.
func validateCommandTemplates(commandTemplates []*template.Template, startTime time.Time, forecastTime time.Duration) error {
	checkerVar := CheckerTemplateVariable{
		ConfigTemplateVariable: dataconfig.GenerateConfigTemplateVariable(startTime, forecastTime, ""),
		FilePath:               "dry-run-file",
	}
	for i, commandTemplate := range commandTemplates {
		if err := commandTemplate.Execute(io.Discard, checkerVar); err != nil {
			return fmt.Errorf("execute command template %d for forecast time %s failed: %w",
				i, dataconfig.FormatForecastTimeShort(forecastTime), err)
		}
	}
	return nil
}

// DryRun builds the real run plan and writes it without polling, sleeping,
// accessing candidate data paths, starting an executor, or writing a report.
func DryRun(config CheckerConfig, startTime time.Time, stdin io.Reader, output io.Writer) error {
	plan, err := buildRunPlan(config, startTime, stdin)
	if err != nil {
		return err
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "check mode: %s\n", plan.config.CheckMode)
	fmt.Fprintf(&builder, "effective arrival order: %s\n", plan.arrivalOrder)
	fmt.Fprintf(&builder, "forecast time count: %d\n", len(plan.forecastTimeList))
	for _, target := range plan.targets {
		fmt.Fprintf(&builder, "forecast time: %s\n", dataconfig.FormatForecastTimeShort(target.forecastTime))
		builder.WriteString("  candidate paths (planned only; not checked for existence or selection):\n")
		if len(target.candidatePaths) == 0 {
			builder.WriteString("    (none)\n")
			continue
		}
		for _, candidatePath := range target.candidatePaths {
			fmt.Fprintf(&builder, "    %s\n", candidatePath)
		}
	}
	builder.WriteString("command templates:\n")
	if len(plan.commandTemplates) == 0 {
		builder.WriteString("  (none)\n")
	} else if plan.config.ExecuteCommand != "" {
		fmt.Fprintf(&builder, "  %s\n", plan.config.ExecuteCommand)
	} else {
		for _, command := range plan.config.ExecuteCommands {
			fmt.Fprintf(&builder, "  %s\n", command)
		}
	}
	_, err = io.WriteString(output, builder.String())
	return err
}
