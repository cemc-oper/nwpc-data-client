package e2e

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// checkerConfig mirrors the YAML structure of the checker's runtime config.
type checkerConfig struct {
	DataConfigDir  string   `yaml:"data_config_dir,omitempty"`
	DataConfigFile string   `yaml:"data_config_file,omitempty"`
	DataType       string   `yaml:"data_type,omitempty"`
	CheckInterval  string   `yaml:"check_interval,omitempty"`
	MaxCheckCount  int      `yaml:"max_check_count,omitempty"`
	DelayTime      string   `yaml:"delay_time,omitempty"`
	CheckMode      string   `yaml:"check_mode,omitempty"`
	ArrivalOrder   string   `yaml:"arrival_order,omitempty"`
	ForecastTimes  []string `yaml:"forecast_times"`
	ExecuteCommand string   `yaml:"execute_command,omitempty"`
	CommandTimeout string   `yaml:"command_timeout,omitempty"`
	ReportFile     string   `yaml:"report_file,omitempty"`
}

type runReport struct {
	ReportVersion int    `json:"report_version"`
	Status        string `json:"status"`
	Summary       struct {
		Total     int `json:"total"`
		Succeeded int `json:"succeeded"`
		Failed    int `json:"failed"`
		Cancelled int `json:"cancelled"`
		Pending   int `json:"pending"`
	} `json:"summary"`
	Results []struct {
		ForecastTime   string  `json:"forecast_time"`
		Status         string  `json:"status"`
		FilePath       string  `json:"file_path"`
		CheckCount     int     `json:"check_count"`
		WaitSeconds    float64 `json:"wait_seconds"`
		CommandSeconds float64 `json:"command_seconds"`
		Error          string  `json:"error"`
	} `json:"results"`
}

func readRunReport(t *testing.T, path string) runReport {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err, "read run report")
	var report runReport
	require.NoError(t, json.Unmarshal(data, &report), "parse run report: %s", data)
	return report
}

// writeCheckerConfig marshals cfg into a temporary checker config YAML file
// and returns its path.
func writeCheckerConfig(t *testing.T, cfg checkerConfig) string {
	t.Helper()
	data, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	configPath := filepath.Join(t.TempDir(), "checker.yaml")
	require.NoError(t, os.WriteFile(configPath, data, 0644))
	return configPath
}

// doneCommand is a shell command template that creates "<file>.done" next to
// the found data file. It works in both /bin/sh and cmd.exe, which the CLI
// uses to run commands on Unix and Windows respectively. The path is left
// unquoted because cmd.exe's /C parsing mangles quoted redirect targets;
// t.TempDir() paths are expected to be space-free.
const doneCommand = `echo done > {{.FilePath}}.done`

// runCheckLocal runs "nwpc_data_client check local" with the given checker
// config, allowing a non-zero exit code. A context timeout keeps a stuck
// polling loop from hanging the test suite.
func runCheckLocal(t *testing.T, checkerConfigPath string) (string, error) {
	t.Helper()
	binPath := getBinaryPath(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binPath, "check", "local", "2026101200",
		"--checker-config="+checkerConfigPath)
	output, err := cmd.CombinedOutput()
	require.NotEqual(t, context.DeadlineExceeded, ctx.Err(), "check local timed out: %s", string(output))
	return string(output), err
}

// TestCheckLocalFileStable is scenario E7: the data file already exists with
// a stable size, so the CLI runs the execute command and exits 0.
func TestCheckLocalFileStable(t *testing.T) {
	rootDir := t.TempDir()
	reportFile := filepath.Join(rootDir, "success-report.json")
	dataFile := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(dataFile), 0755))
	require.NoError(t, os.WriteFile(dataFile, []byte("data"), 0644))

	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir:  filepath.ToSlash(configDir),
		DataType:       e2eDemoDataType,
		CheckInterval:  "50ms",
		MaxCheckCount:  20,
		DelayTime:      "1ms",
		ForecastTimes:  []string{"24h"},
		ExecuteCommand: doneCommand,
		ReportFile:     reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.NoError(t, err, "check local failed: %s", output)
	_, statErr := os.Stat(dataFile + ".done")
	assert.NoError(t, statErr, "execute command did not create the .done file")
	report := readRunReport(t, reportFile)
	assert.Equal(t, 1, report.ReportVersion)
	assert.Equal(t, "success", report.Status)
	assert.Equal(t, 1, report.Summary.Total)
	assert.Equal(t, 1, report.Summary.Succeeded)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "success", report.Results[0].Status)
	assert.Equal(t, dataFile, report.Results[0].FilePath)
	assert.Greater(t, report.Results[0].CheckCount, 0)
}

func TestCheckLocalCommandFailureReturnsError(t *testing.T) {
	rootDir := t.TempDir()
	reportFile := filepath.Join(rootDir, "command-failure-report.json")
	dataFile := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(dataFile), 0755))
	require.NoError(t, os.WriteFile(dataFile, []byte("data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir:  filepath.ToSlash(configDir),
		DataType:       e2eDemoDataType,
		CheckInterval:  "20ms",
		MaxCheckCount:  20,
		DelayTime:      "1ms",
		ForecastTimes:  []string{"24h"},
		ExecuteCommand: "exit 7",
		ReportFile:     reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "check local should return the command failure: %s", output)
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Contains(t, output, "exit code is not 0: 7")
	assert.NotContains(t, output, "Usage:")
	report := readRunReport(t, reportFile)
	assert.Equal(t, "failed", report.Status)
	assert.Equal(t, 1, report.Summary.Failed)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "command_failed", report.Results[0].Status)
	assert.Contains(t, report.Results[0].Error, "exit code is not 0: 7")
}

func TestCheckLocalCommandTimeout(t *testing.T) {
	rootDir := t.TempDir()
	reportFile := filepath.Join(rootDir, "timeout-report.json")
	dataFile := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(dataFile), 0755))
	require.NoError(t, os.WriteFile(dataFile, []byte("data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir:  filepath.ToSlash(configDir),
		DataType:       e2eDemoDataType,
		CheckInterval:  "20ms",
		MaxCheckCount:  20,
		DelayTime:      "1ms",
		ForecastTimes:  []string{"24h"},
		ExecuteCommand: e2eCommandHelperCommand(t, "wait"),
		CommandTimeout: "50ms",
		ReportFile:     reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "check local should return the command timeout: %s", output)
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Contains(t, output, "command timed out after 50ms")
	assert.NotContains(t, output, "Usage:")
	report := readRunReport(t, reportFile)
	assert.Equal(t, "failed", report.Status)
	assert.Equal(t, 1, report.Summary.Failed)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "command_failed", report.Results[0].Status)
	assert.Contains(t, report.Results[0].Error, "command timed out after 50ms")
}

// TestCheckLocalFileNeverAppears is scenario E8: the data file never shows
// up, so the CLI gives up after max_check_count checks and exits non-zero.
func TestCheckLocalFileNeverAppears(t *testing.T) {
	rootDir := t.TempDir()
	reportFile := filepath.Join(rootDir, "check-failure-report.json")
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "20ms",
		MaxCheckCount: 5,
		DelayTime:     "1ms",
		ForecastTimes: []string{"24h"},
		ReportFile:    reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "check local should fail when the file never appears: %s", output)
	assert.Contains(t, output, "too many times")
	assert.NotContains(t, output, "Usage:")
	assert.False(t, strings.Contains(output, "run command success"))
	report := readRunReport(t, reportFile)
	assert.Equal(t, "failed", report.Status)
	assert.Equal(t, 1, report.Summary.Failed)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "check_failed", report.Results[0].Status)
	assert.Equal(t, 5, report.Results[0].CheckCount)
	assert.Contains(t, report.Results[0].Error, "too many times")
}

// TestCheckLocalBatchedStable is scenario E9: in batched mode, both data
// files already exist with stable sizes, so the CLI runs the execute command
// for each and exits 0.
func TestCheckLocalBatchedStable(t *testing.T) {
	rootDir := t.TempDir()
	reportFile := filepath.Join(rootDir, "batched-success-report.json")
	var dataFiles []string
	for _, forecastHour := range []string{"024", "048"} {
		dataFile := filepath.Join(rootDir, "runtime", "2026101200",
			"data.2026101200"+forecastHour+".grb2")
		require.NoError(t, os.MkdirAll(filepath.Dir(dataFile), 0755))
		require.NoError(t, os.WriteFile(dataFile, []byte("data"), 0644))
		dataFiles = append(dataFiles, dataFile)
	}

	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir:  filepath.ToSlash(configDir),
		DataType:       e2eDemoDataType,
		CheckInterval:  "50ms",
		MaxCheckCount:  20,
		CheckMode:      "batched",
		ForecastTimes:  []string{"24h", "48h"},
		ExecuteCommand: doneCommand,
		ReportFile:     reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.NoError(t, err, "batched check local failed: %s", output)
	for _, dataFile := range dataFiles {
		_, statErr := os.Stat(dataFile + ".done")
		assert.NoError(t, statErr, "execute command did not create the .done file for %s", dataFile)
	}
	report := readRunReport(t, reportFile)
	assert.Equal(t, "success", report.Status)
	assert.Equal(t, 2, report.Summary.Total)
	assert.Equal(t, 2, report.Summary.Succeeded)
	require.Len(t, report.Results, 2)
	for _, result := range report.Results {
		assert.Equal(t, "success", result.Status)
		assert.Greater(t, result.CheckCount, 0)
	}
}

// TestCheckLocalBatchedNeverAppears is scenario E10: in batched mode the file
// never shows up, so the CLI exhausts max_check_count and exits non-zero.
func TestCheckLocalBatchedNeverAppears(t *testing.T) {
	rootDir := t.TempDir()
	reportFile := filepath.Join(rootDir, "batched-check-failure-report.json")
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "20ms",
		MaxCheckCount: 5,
		CheckMode:     "batched",
		ForecastTimes: []string{"24h"},
		ReportFile:    reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "batched check local should fail when the file never appears: %s", output)
	assert.Contains(t, output, "too many times")
	assert.NotContains(t, output, "Usage:")
	assert.False(t, strings.Contains(output, "run command success"))
	report := readRunReport(t, reportFile)
	assert.Equal(t, "failed", report.Status)
	assert.Equal(t, 1, report.Summary.Failed)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "check_failed", report.Results[0].Status)
	assert.Equal(t, 5, report.Results[0].CheckCount)
}

// TestCheckLocalBatchedPathPriority is scenario E11: in batched mode the file
// exists in both the runtime and the archive directory, and the first path in
// the config (runtime) wins.
func TestCheckLocalBatchedPathPriority(t *testing.T) {
	rootDir := t.TempDir()
	runtimeFile := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	archiveFile := filepath.Join(rootDir, "archive", "2026101200", "data.2026101200024.grb2")
	for _, f := range []string{runtimeFile, archiveFile} {
		require.NoError(t, os.MkdirAll(filepath.Dir(f), 0755))
		require.NoError(t, os.WriteFile(f, []byte("data"), 0644))
	}

	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir:  filepath.ToSlash(configDir),
		DataType:       e2eDemoDataType,
		CheckInterval:  "50ms",
		MaxCheckCount:  20,
		CheckMode:      "batched",
		ForecastTimes:  []string{"24h"},
		ExecuteCommand: doneCommand,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.NoError(t, err, "batched check local failed: %s", output)
	_, statErr := os.Stat(runtimeFile + ".done")
	assert.NoError(t, statErr, "execute command did not run for the runtime file")
	_, statErr = os.Stat(archiveFile + ".done")
	assert.True(t, os.IsNotExist(statErr), "archive file must not be used when runtime hits")
}

func TestCheckLocalFailFastReportKeepsPending(t *testing.T) {
	rootDir := t.TempDir()
	reportFile := filepath.Join(rootDir, "fail-fast-report.json")
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "1ms",
		MaxCheckCount: 1,
		DelayTime:     "1h",
		ForecastTimes: []string{"24h", "48h"},
		ReportFile:    reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "check local should fail fast: %s", output)
	report := readRunReport(t, reportFile)
	assert.Equal(t, "failed", report.Status)
	assert.Equal(t, 2, report.Summary.Total)
	assert.Equal(t, 1, report.Summary.Failed)
	assert.Equal(t, 1, report.Summary.Pending)
	require.Len(t, report.Results, 2)
	assert.Equal(t, "check_failed", report.Results[0].Status)
	assert.Equal(t, "pending", report.Results[1].Status)
	assert.Zero(t, report.Results[1].CheckCount)
}

func TestCheckLocalReportWriteFailure(t *testing.T) {
	rootDir := t.TempDir()
	dataFile := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(dataFile), 0755))
	require.NoError(t, os.WriteFile(dataFile, []byte("data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	reportFile := filepath.Join(rootDir, "missing", "report.json")
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "20ms",
		MaxCheckCount: 20,
		DelayTime:     "0s",
		ForecastTimes: []string{"24h"},
		ReportFile:    reportFile,
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "report write failure must make the CLI fail: %s", output)
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Contains(t, output, "create temporary run report")
	assert.NotContains(t, output, "Usage:")
}
