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
)

func TestCheckLocalDryRunYAMLAndCLIOverrides(t *testing.T) {
	rootDir := t.TempDir()
	existingFile := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(existingFile), 0755))
	require.NoError(t, os.WriteFile(existingFile, []byte("original"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	reportFile := filepath.Join(t.TempDir(), "report.json")
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir:  filepath.ToSlash(configDir),
		DataType:       e2eDemoDataType,
		CheckInterval:  "1h",
		DelayTime:      "1h",
		MaxCheckCount:  1,
		ForecastTimes:  []string{"24h"},
		ExecuteCommand: doneCommand,
		ReportFile:     reportFile,
	})

	output, err := runCheckLocalDryRun(t, checkerConfigPath, "48h",
		"--check-mode=batched", "--arrival-order=unordered", "--location-level=runtime")
	require.NoError(t, err, output)
	assert.Contains(t, output, "check mode: batched")
	assert.Contains(t, output, "effective arrival order: unordered")
	assert.Contains(t, output, "forecast time count: 1")
	assert.Contains(t, output, existingFile)
	assert.NotContains(t, output, "data.2026101200048.grb2", "stdin must not override YAML forecast times")
	assert.NotContains(t, output, filepath.Join(rootDir, "archive"), "CLI location level was not applied")

	content, err := os.ReadFile(existingFile)
	require.NoError(t, err)
	assert.Equal(t, "original", string(content))
	_, err = os.Stat(existingFile + ".done")
	assert.True(t, os.IsNotExist(err), "dry-run executed the command")
	_, err = os.Stat(reportFile)
	assert.True(t, os.IsNotExist(err), "dry-run created a report")
}

func TestCheckLocalDryRunReadsForecastTimesFromStdin(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "1h",
		MaxCheckCount: 1,
		ForecastTimes: []string{},
	})

	output, err := runCheckLocalDryRun(t, checkerConfigPath, "24h 48h")
	require.NoError(t, err, output)
	assert.Contains(t, output, "forecast time count: 2")
	assert.Contains(t, output, "data.2026101200024.grb2")
	assert.Contains(t, output, "data.2026101200048.grb2")
	_, err = os.Stat(filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2"))
	assert.True(t, os.IsNotExist(err), "dry-run created a candidate data file")
}

func TestCheckLocalDryRunErrors(t *testing.T) {
	t.Run("missing data type does not fall back to embedded", func(t *testing.T) {
		checkerConfigPath := writeCheckerConfig(t, checkerConfig{
			DataConfigDir: filepath.ToSlash(t.TempDir()),
			DataType:      "cma_gfs_gmf/current/grib2/orig",
			ForecastTimes: []string{"24h"},
		})
		output, err := runCheckLocalDryRun(t, checkerConfigPath, "")
		require.Error(t, err)
		assert.Contains(t, output, "model data type config is not found")
	})

	t.Run("data config template", func(t *testing.T) {
		dataConfigPath := filepath.Join(t.TempDir(), "invalid.yaml")
		require.NoError(t, os.WriteFile(dataConfigPath, []byte("file_name: '{{'\npaths: []\n"), 0644))
		checkerConfigPath := writeCheckerConfig(t, checkerConfig{
			DataConfigFile: filepath.ToSlash(dataConfigPath),
			ForecastTimes:  []string{"24h"},
		})
		output, err := runCheckLocalDryRun(t, checkerConfigPath, "")
		require.Error(t, err)
		assert.Contains(t, output, "parse data config template")
	})

	t.Run("command template variable", func(t *testing.T) {
		configDir := writeConfigWithRoot(t, e2eDemoConfigPath, t.TempDir())
		checkerConfigPath := writeCheckerConfig(t, checkerConfig{
			DataConfigDir:  filepath.ToSlash(configDir),
			DataType:       e2eDemoDataType,
			ForecastTimes:  []string{"24h"},
			ExecuteCommand: "echo {{.Unknown}}",
		})
		output, err := runCheckLocalDryRun(t, checkerConfigPath, "")
		require.Error(t, err)
		assert.Contains(t, output, "execute command template")
		assert.Contains(t, output, "Unknown")
	})

	t.Run("no forecast times", func(t *testing.T) {
		configDir := writeConfigWithRoot(t, e2eDemoConfigPath, t.TempDir())
		checkerConfigPath := writeCheckerConfig(t, checkerConfig{
			DataConfigDir: filepath.ToSlash(configDir),
			DataType:      e2eDemoDataType,
			ForecastTimes: []string{},
		})
		output, err := runCheckLocalDryRun(t, checkerConfigPath, "")
		require.Error(t, err)
		assert.Contains(t, output, "no forecast times")
	})
}

func runCheckLocalDryRun(t *testing.T, checkerConfigPath, stdin string, extraArgs ...string) (string, error) {
	t.Helper()
	binPath := getBinaryPath(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := []string{"check", "local", "2026101200", "--checker-config=" + checkerConfigPath, "--dry-run"}
	args = append(args, extraArgs...)
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdin = strings.NewReader(stdin)
	output, err := cmd.CombinedOutput()
	require.NotEqual(t, context.DeadlineExceeded, ctx.Err(), "dry-run blocked or slept: %s", string(output))
	return string(output), err
}
