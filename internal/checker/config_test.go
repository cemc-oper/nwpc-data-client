package checker

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCheckerConfig(t *testing.T) {
	cfg := DefaultCheckerConfig()
	assert.Equal(t, 2880, cfg.MaxCheckCount)
	assert.Equal(t, "5s", cfg.CheckInterval)
	assert.Equal(t, "10s", cfg.DelayTime)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 0.1, cfg.CheckIntervalJitter)
	assert.Equal(t, 4, cfg.MaxConcurrentStableChecks)
	assert.Equal(t, CheckModePerTime, cfg.CheckMode)
	assert.Equal(t, ArrivalOrderAuto, cfg.ArrivalOrder)
}

func TestLoadCheckerConfig(t *testing.T) {
	content := `
data_type: gmfsgrapes_global
location_levels: runtime,storage
max_check_count: 100
check_interval: 10s
delay_time: 5s
debug: true
check_interval_jitter: 0.2
max_concurrent_stable_checks: 8
check_mode: batched
forecast_times:
  - 0h
  - 6h
  - 12h
execute_commands:
  - echo "found {{.FilePath}}"
  - /app/postprocess.sh {{.FilePath}}
`
	dir := t.TempDir()
	configPath := filepath.Join(dir, "checker.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0644))

	cfg, err := LoadCheckerConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, "gmfsgrapes_global", cfg.DataType)
	assert.Equal(t, "runtime,storage", cfg.LocationLevels)
	assert.Equal(t, 100, cfg.MaxCheckCount)
	assert.Equal(t, "10s", cfg.CheckInterval)
	assert.Equal(t, "5s", cfg.DelayTime)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 0.2, cfg.CheckIntervalJitter)
	assert.Equal(t, 8, cfg.MaxConcurrentStableChecks)
	assert.Equal(t, "batched", cfg.CheckMode)
	assert.Equal(t, []string{"0h", "6h", "12h"}, cfg.ForecastTimes)
	assert.Equal(t, "", cfg.ExecuteCommand)
	assert.Equal(t, []string{`echo "found {{.FilePath}}"`, `/app/postprocess.sh {{.FilePath}}`}, cfg.ExecuteCommands)
}

func TestMergeCheckerConfigDefaults(t *testing.T) {
	cmd := newTestLocalCommand(t)
	fileConfig := CheckerConfig{DataType: "some/type"}
	cfg, err := MergeCheckerConfig(fileConfig, CheckerConfig{}, cmd)
	require.NoError(t, err)

	assert.Equal(t, "some/type", cfg.DataType)
	assert.Equal(t, DefaultCheckerConfig().MaxCheckCount, cfg.MaxCheckCount)
	assert.Equal(t, DefaultCheckerConfig().CheckInterval, cfg.CheckInterval)
	assert.Equal(t, DefaultCheckerConfig().DelayTime, cfg.DelayTime)
	assert.False(t, cfg.Debug)
}

func TestMergeCheckerConfigYAMLValues(t *testing.T) {
	cmd := newTestLocalCommand(t)
	fileConfig := CheckerConfig{
		DataType:        "cma_gfs_gmf/current/grib2/orig",
		LocationLevels:  "runtime",
		MaxCheckCount:   100,
		CheckInterval:   "10s",
		DelayTime:       "5s",
		Debug:           true,
		ForecastTimes:   []string{"0h", "6h"},
		ExecuteCommands: []string{"echo ok"},
	}

	cfg, err := MergeCheckerConfig(fileConfig, CheckerConfig{DataType: "ignored"}, cmd)
	require.NoError(t, err)

	assert.Equal(t, "cma_gfs_gmf/current/grib2/orig", cfg.DataType)
	assert.Equal(t, "runtime", cfg.LocationLevels)
	assert.Equal(t, 100, cfg.MaxCheckCount)
	assert.Equal(t, "10s", cfg.CheckInterval)
	assert.Equal(t, "5s", cfg.DelayTime)
	assert.True(t, cfg.Debug)
	assert.Equal(t, []string{"0h", "6h"}, cfg.ForecastTimes)
	assert.Equal(t, []string{"echo ok"}, cfg.ExecuteCommands)
}

func TestMergeCheckerConfigCLIOverridesYAML(t *testing.T) {
	cmd := newTestLocalCommand(t)
	require.NoError(t, cmd.Flags().Set("data-type", "override/type"))
	require.NoError(t, cmd.Flags().Set("max-check-count", "50"))
	require.NoError(t, cmd.Flags().Set("check-interval", "20s"))
	require.NoError(t, cmd.Flags().Set("check-interval-jitter", "0.3"))
	require.NoError(t, cmd.Flags().Set("max-concurrent-stable-checks", "2"))

	fileConfig := CheckerConfig{
		DataType:                  "original/type",
		MaxCheckCount:             100,
		CheckInterval:             "10s",
		CheckIntervalJitter:       0.2,
		MaxConcurrentStableChecks: 8,
	}
	flagConfig := CheckerConfig{
		DataType:                  "override/type",
		MaxCheckCount:             50,
		CheckInterval:             "20s",
		CheckIntervalJitter:       0.3,
		MaxConcurrentStableChecks: 2,
	}

	cfg, err := MergeCheckerConfig(fileConfig, flagConfig, cmd)
	require.NoError(t, err)

	assert.Equal(t, "override/type", cfg.DataType)
	assert.Equal(t, 50, cfg.MaxCheckCount)
	assert.Equal(t, "20s", cfg.CheckInterval)
	assert.Equal(t, 0.3, cfg.CheckIntervalJitter)
	assert.Equal(t, 2, cfg.MaxConcurrentStableChecks)
}

func TestMergeCheckerConfigExecuteCommandOverridesExecuteCommands(t *testing.T) {
	cmd := newTestLocalCommand(t)
	require.NoError(t, cmd.Flags().Set("execute-command", "echo override"))

	fileConfig := CheckerConfig{
		DataType:        "cma_gfs_gmf/current/grib2/orig",
		ExecuteCommands: []string{"echo from yaml"},
	}
	flagConfig := CheckerConfig{
		DataType:       "cma_gfs_gmf/current/grib2/orig",
		ExecuteCommand: "echo override",
	}

	cfg, err := MergeCheckerConfig(fileConfig, flagConfig, cmd)
	require.NoError(t, err)

	assert.Equal(t, "echo override", cfg.ExecuteCommand)
	assert.Empty(t, cfg.ExecuteCommands)
}

func TestMergeCheckerConfigMissingDataSource(t *testing.T) {
	cmd := newTestLocalCommand(t)
	_, err := MergeCheckerConfig(CheckerConfig{}, CheckerConfig{}, cmd)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data_type")
	assert.Contains(t, err.Error(), "data_config_file")
}

func newTestLocalCommand(t *testing.T) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "check-local"}
	cmd.Flags().String("data-config-dir", "", "")
	cmd.Flags().String("data-config-file", "", "")
	cmd.Flags().String("checker-config", "", "")
	cmd.Flags().String("data-type", "", "")
	cmd.Flags().String("location-level", "", "")
	cmd.Flags().Int("max-check-count", 2880, "")
	cmd.Flags().String("check-interval", "5s", "")
	cmd.Flags().String("execute-command", "", "")
	cmd.Flags().String("delay-time", "10s", "")
	cmd.Flags().Float64("check-interval-jitter", 0.1, "")
	cmd.Flags().Int("max-concurrent-stable-checks", 4, "")
	cmd.Flags().String("check-mode", "per_time", "")
	cmd.Flags().String("arrival-order", "auto", "")
	cmd.Flags().Bool("debug", false, "")
	return cmd
}

func TestCheckerConfigValidate(t *testing.T) {
	validDataType := CheckerConfig{DataType: "some/type", ExecuteCommand: "echo ok"}
	assert.NoError(t, validDataType.Validate())

	validDataConfigFile := CheckerConfig{DataConfigFile: "/path/to/config.yaml", ExecuteCommands: []string{"echo ok"}}
	assert.NoError(t, validDataConfigFile.Validate())

	missingDataSource := CheckerConfig{ExecuteCommand: "echo ok"}
	assert.Error(t, missingDataSource.Validate())

	invalid := CheckerConfig{
		DataType:        "some/type",
		ExecuteCommand:  "echo ok",
		ExecuteCommands: []string{"echo ok"},
	}
	assert.Error(t, invalid.Validate())

	negativeJitter := CheckerConfig{DataType: "some/type", CheckIntervalJitter: -0.1}
	assert.Error(t, negativeJitter.Validate())

	tooLargeJitter := CheckerConfig{DataType: "some/type", CheckIntervalJitter: 0.6}
	assert.Error(t, tooLargeJitter.Validate())

	boundaryJitter := CheckerConfig{DataType: "some/type", CheckIntervalJitter: 0.5}
	assert.NoError(t, boundaryJitter.Validate())

	invalidCheckMode := CheckerConfig{DataType: "some/type", CheckMode: "nonsense"}
	assert.Error(t, invalidCheckMode.Validate())

	assert.NoError(t, CheckerConfig{DataType: "some/type", CheckMode: CheckModePerTime}.Validate())
	assert.NoError(t, CheckerConfig{DataType: "some/type", CheckMode: CheckModeBatched}.Validate())
	// Empty means "not merged yet" and falls back to the default.
	assert.NoError(t, CheckerConfig{DataType: "some/type"}.Validate())

	assert.NoError(t, CheckerConfig{DataType: "some/type", ArrivalOrder: ArrivalOrderAuto}.Validate())
	assert.NoError(t, CheckerConfig{DataType: "some/type", ArrivalOrder: "sequential"}.Validate())
	assert.NoError(t, CheckerConfig{DataType: "some/type", ArrivalOrder: "unordered"}.Validate())
	assert.Error(t, CheckerConfig{DataType: "some/type", ArrivalOrder: "weekly"}.Validate())
}

func TestCheckerConfigParseForecastTimes(t *testing.T) {
	cfg := CheckerConfig{ForecastTimes: []string{"0h", "001h10m", "120h"}}
	durations, err := cfg.ParseForecastTimes()
	require.NoError(t, err)

	expected := []time.Duration{
		0 * time.Hour,
		1*time.Hour + 10*time.Minute,
		120 * time.Hour,
	}
	assert.Equal(t, expected, durations)
}

func TestCheckerConfigParseForecastTimesEmpty(t *testing.T) {
	cfg := CheckerConfig{}
	durations, err := cfg.ParseForecastTimes()
	require.NoError(t, err)
	assert.Empty(t, durations)
}

func TestCheckerConfigParseForecastTimesInvalid(t *testing.T) {
	cfg := CheckerConfig{ForecastTimes: []string{"not-a-duration"}}
	_, err := cfg.ParseForecastTimes()
	assert.Error(t, err)
}
