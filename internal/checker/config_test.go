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
	assert.Equal(t, "30m", cfg.CommandTimeout)
	assert.Empty(t, cfg.ReportFile)
}

func TestLoadCheckerConfig(t *testing.T) {
	content := `
data_type: gmfsgrapes_global
location_levels: runtime,storage
storage_base: /mnt/cmadaas
max_check_count: 100
check_interval: 10s
delay_time: 5s
command_timeout: 45m
report_file: /tmp/checker-report.json
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
	assert.Equal(t, "/mnt/cmadaas", cfg.StorageBase)
	assert.Equal(t, 100, cfg.MaxCheckCount)
	assert.Equal(t, "10s", cfg.CheckInterval)
	assert.Equal(t, "5s", cfg.DelayTime)
	assert.Equal(t, "45m", cfg.CommandTimeout)
	assert.Equal(t, "/tmp/checker-report.json", cfg.ReportFile)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 0.2, cfg.CheckIntervalJitter)
	assert.Equal(t, 8, cfg.MaxConcurrentStableChecks)
	assert.Equal(t, "batched", cfg.CheckMode)
	assert.Equal(t, []string{"0h", "6h", "12h"}, cfg.ForecastTimes)
	assert.Equal(t, "", cfg.ExecuteCommand)
	assert.Equal(t, []string{`echo "found {{.FilePath}}"`, `/app/postprocess.sh {{.FilePath}}`}, cfg.ExecuteCommands)
}

func TestLoadCheckerConfigRejectsUnknownField(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "unknown.yaml")
	content := "data_type: some/type\ncheck_intervel: 10s\n"
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0644))

	_, err := LoadCheckerConfig(configPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), configPath)
	assert.Contains(t, err.Error(), "check_intervel")
	assert.Contains(t, err.Error(), "not found")
}

func TestLoadCheckerConfigReadErrorIncludesPath(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "missing.yaml")
	_, err := LoadCheckerConfig(configPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), configPath)
}

func TestLoadCheckerConfigAllowsEmptyFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "empty.yaml")
	require.NoError(t, os.WriteFile(configPath, nil, 0644))

	config, err := LoadCheckerConfig(configPath)
	require.NoError(t, err)
	assert.Empty(t, config.DataType)
	assert.NotNil(t, config.yamlFields)
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
		StorageBase:     "/mnt/yaml-cmadaas",
		MaxCheckCount:   100,
		CheckInterval:   "10s",
		DelayTime:       "5s",
		Debug:           true,
		ForecastTimes:   []string{"0h", "6h"},
		ExecuteCommands: []string{"echo ok"},
		CommandTimeout:  "45m",
		ReportFile:      "/tmp/report.json",
	}

	cfg, err := MergeCheckerConfig(fileConfig, CheckerConfig{DataType: "ignored"}, cmd)
	require.NoError(t, err)

	assert.Equal(t, "cma_gfs_gmf/current/grib2/orig", cfg.DataType)
	assert.Equal(t, "runtime", cfg.LocationLevels)
	assert.Equal(t, "/mnt/yaml-cmadaas", cfg.StorageBase)
	assert.Equal(t, 100, cfg.MaxCheckCount)
	assert.Equal(t, "10s", cfg.CheckInterval)
	assert.Equal(t, "5s", cfg.DelayTime)
	assert.True(t, cfg.Debug)
	assert.Equal(t, []string{"0h", "6h"}, cfg.ForecastTimes)
	assert.Equal(t, []string{"echo ok"}, cfg.ExecuteCommands)
	assert.Equal(t, "45m", cfg.CommandTimeout)
	assert.Equal(t, "/tmp/report.json", cfg.ReportFile)
}

func TestMergeCheckerConfigCLIOverridesYAML(t *testing.T) {
	cmd := newTestLocalCommand(t)
	require.NoError(t, cmd.Flags().Set("data-type", "override/type"))
	require.NoError(t, cmd.Flags().Set("storage-base", "/mnt/cli-cmadaas"))
	require.NoError(t, cmd.Flags().Set("max-check-count", "50"))
	require.NoError(t, cmd.Flags().Set("check-interval", "20s"))
	require.NoError(t, cmd.Flags().Set("check-interval-jitter", "0.3"))
	require.NoError(t, cmd.Flags().Set("max-concurrent-stable-checks", "2"))
	require.NoError(t, cmd.Flags().Set("command-timeout", "1h"))
	require.NoError(t, cmd.Flags().Set("report-file", "/tmp/cli-report.json"))

	fileConfig := CheckerConfig{
		DataType:                  "original/type",
		StorageBase:               "/mnt/yaml-cmadaas",
		MaxCheckCount:             100,
		CheckInterval:             "10s",
		CheckIntervalJitter:       0.2,
		MaxConcurrentStableChecks: 8,
		CommandTimeout:            "45m",
		ReportFile:                "/tmp/yaml-report.json",
	}
	flagConfig := CheckerConfig{
		DataType:                  "override/type",
		StorageBase:               "/mnt/cli-cmadaas",
		MaxCheckCount:             50,
		CheckInterval:             "20s",
		CheckIntervalJitter:       0.3,
		MaxConcurrentStableChecks: 2,
		CommandTimeout:            "1h",
		ReportFile:                "/tmp/cli-report.json",
	}

	cfg, err := MergeCheckerConfig(fileConfig, flagConfig, cmd)
	require.NoError(t, err)

	assert.Equal(t, "override/type", cfg.DataType)
	assert.Equal(t, "/mnt/cli-cmadaas", cfg.StorageBase)
	assert.Equal(t, 50, cfg.MaxCheckCount)
	assert.Equal(t, "20s", cfg.CheckInterval)
	assert.Equal(t, 0.3, cfg.CheckIntervalJitter)
	assert.Equal(t, 2, cfg.MaxConcurrentStableChecks)
	assert.Equal(t, "1h", cfg.CommandTimeout)
	assert.Equal(t, "/tmp/cli-report.json", cfg.ReportFile)
}

func TestMergeCheckerConfigYAMLExplicitZeroValues(t *testing.T) {
	fileConfig := loadCheckerConfigForTest(t, `
data_type: some/type
check_interval_jitter: 0
command_timeout: 0
debug: false
report_file: ""
forecast_times: []
`)

	config, err := MergeCheckerConfig(fileConfig, CheckerConfig{}, newTestLocalCommand(t))
	require.NoError(t, err)
	assert.Zero(t, config.CheckIntervalJitter)
	assert.Equal(t, "0", config.CommandTimeout)
	assert.False(t, config.Debug)
	assert.Empty(t, config.ReportFile)
	assert.Empty(t, config.ForecastTimes)
}

func TestMergeCheckerConfigRejectsInvalidExplicitYAMLValues(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		wantError string
	}{
		{
			name:      "relative storage base",
			yaml:      "storage_base: relative/mount",
			wantError: "storage base",
		},
		{
			name:      "zero max check count",
			yaml:      "max_check_count: 0",
			wantError: "max_check_count",
		},
		{
			name:      "zero stable check concurrency",
			yaml:      "max_concurrent_stable_checks: 0",
			wantError: "max_concurrent_stable_checks",
		},
		{
			name:      "empty check interval",
			yaml:      `check_interval: ""`,
			wantError: "check_interval",
		},
		{
			name:      "negative command timeout",
			yaml:      "command_timeout: -1s",
			wantError: "command_timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileConfig := loadCheckerConfigForTest(t, "data_type: some/type\n"+tt.yaml+"\n")
			_, err := MergeCheckerConfig(fileConfig, CheckerConfig{}, newTestLocalCommand(t))
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}

func TestMergeCheckerConfigCLIExplicitZeroFalseAndEmpty(t *testing.T) {
	cmd := newTestLocalCommand(t)
	require.NoError(t, cmd.Flags().Set("check-interval-jitter", "0"))
	require.NoError(t, cmd.Flags().Set("command-timeout", "0"))
	require.NoError(t, cmd.Flags().Set("debug", "false"))
	require.NoError(t, cmd.Flags().Set("report-file", ""))

	fileConfig := loadCheckerConfigForTest(t, `
data_type: some/type
check_interval_jitter: 0.3
command_timeout: 45m
debug: true
report_file: /tmp/yaml-report.json
`)
	flagConfig := CheckerConfig{
		CheckIntervalJitter: 0,
		CommandTimeout:      "0",
		Debug:               false,
		ReportFile:          "",
	}

	config, err := MergeCheckerConfig(fileConfig, flagConfig, cmd)
	require.NoError(t, err)
	assert.Zero(t, config.CheckIntervalJitter)
	assert.Equal(t, "0", config.CommandTimeout)
	assert.False(t, config.Debug)
	assert.Empty(t, config.ReportFile)
}

func TestMergeCheckerConfigCLIExplicitEmptyDataType(t *testing.T) {
	cmd := newTestLocalCommand(t)
	require.NoError(t, cmd.Flags().Set("data-type", ""))
	fileConfig := loadCheckerConfigForTest(t, "data_type: some/type\n")

	_, err := MergeCheckerConfig(fileConfig, CheckerConfig{}, cmd)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data_type")
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

func loadCheckerConfigForTest(t *testing.T, content string) CheckerConfig {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "checker.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0644))
	config, err := LoadCheckerConfig(configPath)
	require.NoError(t, err)
	return config
}

func newTestLocalCommand(t *testing.T) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "check-local"}
	cmd.Flags().String("data-config-dir", "", "")
	cmd.Flags().String("data-config-file", "", "")
	cmd.Flags().String("checker-config", "", "")
	cmd.Flags().String("data-type", "", "")
	cmd.Flags().String("location-level", "", "")
	cmd.Flags().String("storage-base", "", "")
	cmd.Flags().Int("max-check-count", 2880, "")
	cmd.Flags().String("check-interval", "5s", "")
	cmd.Flags().String("execute-command", "", "")
	cmd.Flags().String("command-timeout", "30m", "")
	cmd.Flags().String("report-file", "", "")
	cmd.Flags().String("delay-time", "10s", "")
	cmd.Flags().Float64("check-interval-jitter", 0.1, "")
	cmd.Flags().Int("max-concurrent-stable-checks", 4, "")
	cmd.Flags().String("check-mode", "per_time", "")
	cmd.Flags().String("arrival-order", "auto", "")
	cmd.Flags().Bool("debug", false, "")
	return cmd
}

func TestCheckerConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*CheckerConfig)
		wantError string
	}{
		{name: "valid data type"},
		{
			name: "valid data config file",
			mutate: func(c *CheckerConfig) {
				c.DataType = ""
				c.DataConfigFile = "/path/to/config.yaml"
			},
		},
		{
			name: "data config dir with data type",
			mutate: func(c *CheckerConfig) {
				c.DataConfigDir = "/path/to/config"
			},
		},
		{
			name: "missing data source",
			mutate: func(c *CheckerConfig) {
				c.DataType = ""
			},
			wantError: "data_type",
		},
		{
			name: "data config dir alone",
			mutate: func(c *CheckerConfig) {
				c.DataType = ""
				c.DataConfigDir = "/path/to/config"
			},
			wantError: "cannot be used alone",
		},
		{
			name: "invalid check interval",
			mutate: func(c *CheckerConfig) {
				c.CheckInterval = "often"
			},
			wantError: "check_interval",
		},
		{
			name: "empty check interval",
			mutate: func(c *CheckerConfig) {
				c.CheckInterval = ""
			},
			wantError: "check_interval",
		},
		{
			name: "zero check interval",
			mutate: func(c *CheckerConfig) {
				c.CheckInterval = "0s"
			},
			wantError: "greater than 0",
		},
		{
			name: "negative check interval",
			mutate: func(c *CheckerConfig) {
				c.CheckInterval = "-1s"
			},
			wantError: "greater than 0",
		},
		{
			name: "invalid delay time",
			mutate: func(c *CheckerConfig) {
				c.DelayTime = "later"
			},
			wantError: "delay_time",
		},
		{
			name: "invalid command timeout",
			mutate: func(c *CheckerConfig) {
				c.CommandTimeout = "eventually"
			},
			wantError: "command_timeout",
		},
		{
			name: "negative command timeout",
			mutate: func(c *CheckerConfig) {
				c.CommandTimeout = "-1s"
			},
			wantError: "command_timeout",
		},
		{
			name: "zero disables command timeout",
			mutate: func(c *CheckerConfig) {
				c.CommandTimeout = "0"
			},
		},
		{
			name: "zero max check count",
			mutate: func(c *CheckerConfig) {
				c.MaxCheckCount = 0
			},
			wantError: "max_check_count",
		},
		{
			name: "negative max check count",
			mutate: func(c *CheckerConfig) {
				c.MaxCheckCount = -1
			},
			wantError: "max_check_count",
		},
		{
			name: "zero stable check concurrency",
			mutate: func(c *CheckerConfig) {
				c.MaxConcurrentStableChecks = 0
			},
			wantError: "max_concurrent_stable_checks",
		},
		{
			name: "negative stable check concurrency",
			mutate: func(c *CheckerConfig) {
				c.MaxConcurrentStableChecks = -1
			},
			wantError: "max_concurrent_stable_checks",
		},
		{
			name: "negative jitter",
			mutate: func(c *CheckerConfig) {
				c.CheckIntervalJitter = -0.1
			},
			wantError: "check_interval_jitter",
		},
		{
			name: "too large jitter",
			mutate: func(c *CheckerConfig) {
				c.CheckIntervalJitter = 0.6
			},
			wantError: "check_interval_jitter",
		},
		{
			name: "zero jitter",
			mutate: func(c *CheckerConfig) {
				c.CheckIntervalJitter = 0
			},
		},
		{
			name: "boundary jitter",
			mutate: func(c *CheckerConfig) {
				c.CheckIntervalJitter = 0.5
			},
		},
		{
			name: "invalid check mode",
			mutate: func(c *CheckerConfig) {
				c.CheckMode = "nonsense"
			},
			wantError: "check_mode",
		},
		{
			name: "batched check mode",
			mutate: func(c *CheckerConfig) {
				c.CheckMode = CheckModeBatched
			},
		},
		{
			name: "invalid arrival order",
			mutate: func(c *CheckerConfig) {
				c.ArrivalOrder = "weekly"
			},
			wantError: "arrival_order",
		},
		{
			name: "invalid forecast time",
			mutate: func(c *CheckerConfig) {
				c.ForecastTimes = []string{"0h", "not-a-duration"}
			},
			wantError: "forecast time",
		},
		{
			name: "mutually exclusive commands",
			mutate: func(c *CheckerConfig) {
				c.ExecuteCommands = []string{"echo ok"}
			},
			wantError: "mutually exclusive",
		},
		{
			name: "invalid command template",
			mutate: func(c *CheckerConfig) {
				c.ExecuteCommand = "echo {{"
			},
			wantError: "parse execute command",
		},
		{
			name: "invalid command list template",
			mutate: func(c *CheckerConfig) {
				c.ExecuteCommand = ""
				c.ExecuteCommands = []string{"echo ok", "echo {{.FilePath"}
			},
			wantError: "parse execute command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultCheckerConfig()
			config.DataType = "some/type"
			config.ExecuteCommand = "echo {{.FilePath}}"
			if tt.mutate != nil {
				tt.mutate(&config)
			}

			err := config.Validate()
			if tt.wantError == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError)
			}
		})
	}
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
