package e2e

import (
	"context"
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
	DataConfigDir  string   `yaml:"data_config_dir"`
	DataType       string   `yaml:"data_type"`
	CheckInterval  string   `yaml:"check_interval"`
	MaxCheckCount  int      `yaml:"max_check_count"`
	DelayTime      string   `yaml:"delay_time"`
	CheckMode      string   `yaml:"check_mode,omitempty"`
	ForecastTimes  []string `yaml:"forecast_times"`
	ExecuteCommand string   `yaml:"execute_command,omitempty"`
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
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.NoError(t, err, "check local failed: %s", output)
	_, statErr := os.Stat(dataFile + ".done")
	assert.NoError(t, statErr, "execute command did not create the .done file")
}

// TestCheckLocalFileNeverAppears is scenario E8: the data file never shows
// up, so the CLI gives up after max_check_count checks and exits non-zero.
func TestCheckLocalFileNeverAppears(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "20ms",
		MaxCheckCount: 5,
		DelayTime:     "1ms",
		ForecastTimes: []string{"24h"},
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "check local should fail when the file never appears: %s", output)
	assert.Contains(t, output, "too many times")
	assert.False(t, strings.Contains(output, "run command success"))
}

// TestCheckLocalBatchedStable is scenario E9: in batched mode, both data
// files already exist with stable sizes, so the CLI runs the execute command
// for each and exits 0.
func TestCheckLocalBatchedStable(t *testing.T) {
	rootDir := t.TempDir()
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
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.NoError(t, err, "batched check local failed: %s", output)
	for _, dataFile := range dataFiles {
		_, statErr := os.Stat(dataFile + ".done")
		assert.NoError(t, statErr, "execute command did not create the .done file for %s", dataFile)
	}
}

// TestCheckLocalBatchedNeverAppears is scenario E10: in batched mode the file
// never shows up, so the CLI exhausts max_check_count and exits non-zero.
func TestCheckLocalBatchedNeverAppears(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "20ms",
		MaxCheckCount: 5,
		CheckMode:     "batched",
		ForecastTimes: []string{"24h"},
	})

	output, err := runCheckLocal(t, checkerConfigPath)

	require.Error(t, err, "batched check local should fail when the file never appears: %s", output)
	assert.Contains(t, output, "too many times")
	assert.False(t, strings.Contains(output, "run command success"))
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
