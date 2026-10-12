package checker

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CheckerConfig holds runtime options for the data_checker local command.
// It is intentionally separate from the data path config (DataConfig).
//
// The `flag` tag maps a struct field to its CLI flag name. It is used by
// MergeCheckerConfig to decide whether a CLI flag should override the YAML
// value.
type CheckerConfig struct {
	DataConfigDir   string   `yaml:"data_config_dir" flag:"data-config-dir"`
	DataConfigFile  string   `yaml:"data_config_file" flag:"data-config-file"`
	DataType        string   `yaml:"data_type" flag:"data-type"`
	LocationLevels  string   `yaml:"location_levels" flag:"location-level"`
	MaxCheckCount   int      `yaml:"max_check_count" flag:"max-check-count"`
	CheckInterval   string   `yaml:"check_interval" flag:"check-interval"`
	DelayTime       string   `yaml:"delay_time" flag:"delay-time"`
	Debug           bool     `yaml:"debug" flag:"debug"`
	ForecastTimes   []string `yaml:"forecast_times"`
	ExecuteCommand  string   `yaml:"execute_command" flag:"execute-command"`
	ExecuteCommands []string `yaml:"execute_commands"`

	// CheckIntervalJitter is the per-round jitter ratio applied to the
	// polling interval: each round sleeps interval×(1±j). 0 disables jitter.
	CheckIntervalJitter float64 `yaml:"check_interval_jitter" flag:"check-interval-jitter"`
	// MaxConcurrentStableChecks caps how many files may be in size-stability
	// observation at once. It only applies to batched mode; per-time
	// goroutines are inherently independent.
	MaxConcurrentStableChecks int `yaml:"max_concurrent_stable_checks" flag:"max-concurrent-stable-checks"`

	// CheckMode selects the polling loop: "per_time" (one goroutine per
	// forecast time) or "batched" (a single loop checks all pending forecast
	// times via directory readdir).
	CheckMode string `yaml:"check_mode" flag:"check-mode"`

	// ArrivalOrder overrides the data config's _arrival_order for batched
	// mode gating: "auto" (default) follows the data config, "sequential"
	// gates polling to the earliest pending forecast time, "unordered"
	// disables gating.
	ArrivalOrder string `yaml:"arrival_order" flag:"arrival-order"`
}

// Check modes accepted by CheckerConfig.CheckMode.
const (
	CheckModePerTime = "per_time"
	CheckModeBatched = "batched"
)

// ArrivalOrderAuto makes the checker follow the data config's
// _arrival_order. The sequential/unordered values themselves are defined in
// internal/dataconfig.
const ArrivalOrderAuto = "auto"

// resolveArrivalOrder applies the layered precedence for arrival gating:
// CLI/YAML (any value other than auto) > the data config's _arrival_order >
// unordered (safe default). Overrides work in both directions: sequential →
// unordered is the emergency switch, unordered → sequential is for trials.
func resolveArrivalOrder(checkerOrder, dataConfigOrder string) string {
	if checkerOrder == dataconfig.ArrivalOrderSequential ||
		checkerOrder == dataconfig.ArrivalOrderUnordered {
		return checkerOrder
	}
	if dataConfigOrder == dataconfig.ArrivalOrderSequential {
		return dataconfig.ArrivalOrderSequential
	}
	return dataconfig.ArrivalOrderUnordered
}

// DefaultCheckerConfig returns the default runtime configuration.
// These values match the CLI flag defaults registered in local.go.
func DefaultCheckerConfig() CheckerConfig {
	return CheckerConfig{
		MaxCheckCount:             2880,
		CheckInterval:             "5s",
		DelayTime:                 "10s",
		Debug:                     false,
		CheckIntervalJitter:       0.1,
		MaxConcurrentStableChecks: 4,
		CheckMode:                 CheckModePerTime,
		ArrivalOrder:              ArrivalOrderAuto,
	}
}

// LoadCheckerConfig loads a runtime checker config from a YAML file.
func LoadCheckerConfig(filePath string) (CheckerConfig, error) {
	var config CheckerConfig
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("read checker config file failed: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("parse checker config file failed: %v", err)
	}
	return config, nil
}

// Validate returns an error if the config is missing required fields or contains
// incompatible settings.
func (c CheckerConfig) Validate() error {
	if c.DataType == "" && c.DataConfigFile == "" {
		return fmt.Errorf(
			"requires data_type (--data-type) or data_config_file (--data-config-file); " +
				"data_config_dir (--data-config-dir) is optional and cannot be used alone")
	}
	if c.ExecuteCommand != "" && len(c.ExecuteCommands) > 0 {
		return fmt.Errorf("execute_command and execute_commands are mutually exclusive")
	}
	if c.CheckIntervalJitter < 0 || c.CheckIntervalJitter > 0.5 {
		return fmt.Errorf("check_interval_jitter must be in [0, 0.5], got %v", c.CheckIntervalJitter)
	}
	// An empty check_mode is accepted here: it means the field never went
	// through MergeCheckerConfig, which fills in the default.
	if c.CheckMode != "" && c.CheckMode != CheckModePerTime && c.CheckMode != CheckModeBatched {
		return fmt.Errorf("check_mode must be %q or %q, got %q", CheckModePerTime, CheckModeBatched, c.CheckMode)
	}
	// Same for arrival_order: empty means "not merged yet".
	if c.ArrivalOrder != "" &&
		c.ArrivalOrder != ArrivalOrderAuto &&
		c.ArrivalOrder != dataconfig.ArrivalOrderSequential &&
		c.ArrivalOrder != dataconfig.ArrivalOrderUnordered {
		return fmt.Errorf("arrival_order must be %q, %q or %q, got %q",
			ArrivalOrderAuto, dataconfig.ArrivalOrderSequential, dataconfig.ArrivalOrderUnordered, c.ArrivalOrder)
	}
	return nil
}

// ParseForecastTimes converts string forecast times (e.g. "6h", "001h10m")
// into time.Duration values using the same parser as stdin input.
func (c CheckerConfig) ParseForecastTimes() ([]time.Duration, error) {
	if len(c.ForecastTimes) == 0 {
		return nil, nil
	}
	var result []time.Duration
	for _, s := range c.ForecastTimes {
		d, err := dataconfig.ParseForecastTime(s)
		if err != nil {
			return nil, fmt.Errorf("parse forecast time %q failed: %v", s, err)
		}
		result = append(result, d)
	}
	return result, nil
}

// MergeCheckerConfig combines a YAML file config with CLI flag values.
// CLI flags always override values from the YAML file. The result is validated
// before being returned.
func MergeCheckerConfig(fileConfig CheckerConfig, flagConfig CheckerConfig, cmd *cobra.Command) (CheckerConfig, error) {
	config := DefaultCheckerConfig()

	// For fields that have a `flag` tag, use the CLI value if the flag was
	// explicitly set; otherwise fall back to the YAML value if it is non-zero.
	mergeFlaggedFields(fileConfig, flagConfig, cmd, &config)

	// Forecast times have no CLI flag, so they come from YAML only.
	if len(fileConfig.ForecastTimes) > 0 {
		config.ForecastTimes = fileConfig.ForecastTimes
	}

	// If --execute-command was set on the CLI, discard any execute_commands
	// from YAML because the two fields are mutually exclusive.
	if cmd.Flags().Changed("execute-command") {
		config.ExecuteCommands = nil
	} else if len(fileConfig.ExecuteCommands) > 0 {
		config.ExecuteCommands = fileConfig.ExecuteCommands
	}

	if err := config.Validate(); err != nil {
		return CheckerConfig{}, err
	}
	return config, nil
}

// mergeFlaggedFields copies values from fileConfig/flagConfig into config for
// fields that carry a `flag` struct tag. CLI flags take precedence over YAML
// values.
func mergeFlaggedFields(fileConfig CheckerConfig, flagConfig CheckerConfig, cmd *cobra.Command, config *CheckerConfig) {
	fileValue := reflect.ValueOf(fileConfig)
	flagValue := reflect.ValueOf(flagConfig)
	configValue := reflect.ValueOf(config).Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		field := configType.Field(i)
		flagName := field.Tag.Get("flag")
		if flagName == "" {
			continue
		}

		configField := configValue.Field(i)
		fileField := fileValue.Field(i)
		flagField := flagValue.Field(i)

		if cmd.Flags().Changed(flagName) {
			configField.Set(flagField)
		} else if !fileField.IsZero() {
			configField.Set(fileField)
		}
	}
}
