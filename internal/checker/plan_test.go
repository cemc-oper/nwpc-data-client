package checker

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRunPlanForecastTimeSources(t *testing.T) {
	dataConfigPath := writeRunPlanDataConfig(t, `
default: NOTFOUND
file_name: data.{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: /data
`)
	startTime := time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC)

	t.Run("YAML takes precedence over stdin", func(t *testing.T) {
		config := newRunPlanTestConfig(dataConfigPath)
		config.ForecastTimes = []string{"6h"}
		plan, err := buildRunPlan(config, startTime, strings.NewReader("0h 24h"))
		require.NoError(t, err)
		require.Len(t, plan.targets, 1)
		assert.Equal(t, 6*time.Hour, plan.targets[0].forecastTime)
	})

	t.Run("stdin is used when YAML is empty", func(t *testing.T) {
		config := newRunPlanTestConfig(dataConfigPath)
		config.ForecastTimes = nil
		plan, err := buildRunPlan(config, startTime, strings.NewReader("0h 24h"))
		require.NoError(t, err)
		require.Len(t, plan.targets, 2)
		assert.Equal(t, time.Duration(0), plan.targets[0].forecastTime)
		assert.Equal(t, 24*time.Hour, plan.targets[1].forecastTime)
	})

	t.Run("no forecast times", func(t *testing.T) {
		config := newRunPlanTestConfig(dataConfigPath)
		config.ForecastTimes = nil
		_, err := buildRunPlan(config, startTime, strings.NewReader(""))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no forecast times")
	})
}

func TestDryRunOutputsResolvedPlan(t *testing.T) {
	dataConfigPath := writeRunPlanDataConfig(t, `
_arrival_order: sequential
default: NOTFOUND
file_names:
  - first.{{.ForecastHour}}.grb2
file_name: second.{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: /data/runtime
  - type: local
    level: archive
    path: /data/archive
`)
	config := newRunPlanTestConfig(dataConfigPath)
	config.CheckMode = CheckModeBatched
	config.LocationLevels = "runtime"
	config.ExecuteCommands = []string{"echo {{.FilePath}}", "echo {{.ForecastHour}}"}

	var output bytes.Buffer
	err := DryRun(config, time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC), nil, &output)
	require.NoError(t, err)
	assert.Contains(t, output.String(), "check mode: batched")
	assert.Contains(t, output.String(), "effective arrival order: sequential")
	assert.Contains(t, output.String(), "forecast time count: 1")
	assert.Contains(t, output.String(), "not checked for existence or selection")
	assert.Contains(t, output.String(), filepath.Join("/data/runtime", "first.000.grb2"))
	assert.Contains(t, output.String(), filepath.Join("/data/runtime", "second.000.grb2"))
	assert.NotContains(t, output.String(), filepath.Join("/data/archive", "first.000.grb2"))
	assert.Contains(t, output.String(), "echo {{.FilePath}}")
	assert.Contains(t, output.String(), "echo {{.ForecastHour}}")
}

func TestBuildRunPlanMissingDataTypeDoesNotFallBack(t *testing.T) {
	config := newRunPlanTestConfig("")
	config.DataConfigDir = t.TempDir()
	// This type exists in the embedded baseline. A non-empty config directory
	// is a full replacement, so its absence must remain an error.
	config.DataType = "cma_gfs_gmf/current/grib2/orig"

	_, err := buildRunPlan(config, time.Now(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "model data type config is not found")
}

func TestBuildRunPlanTemplateErrors(t *testing.T) {
	t.Run("data config template syntax", func(t *testing.T) {
		dataConfigPath := writeRunPlanDataConfig(t, `
default: NOTFOUND
file_name: data.{{.ForecastHour}.grb2
paths: []
`)
		config := newRunPlanTestConfig(dataConfigPath)
		_, err := buildRunPlan(config, time.Now(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parse data config template")
	})

	t.Run("data config template variable", func(t *testing.T) {
		dataConfigPath := writeRunPlanDataConfig(t, `
default: NOTFOUND
file_name: data.{{.Unknown}}.grb2
paths: []
`)
		config := newRunPlanTestConfig(dataConfigPath)
		_, err := buildRunPlan(config, time.Now(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unknown")
	})

	t.Run("command template variable", func(t *testing.T) {
		dataConfigPath := writeRunPlanDataConfig(t, testRunPlanDataConfig)
		config := newRunPlanTestConfig(dataConfigPath)
		config.ExecuteCommand = "echo {{.Unknown}}"
		_, err := buildRunPlan(config, time.Now(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "execute command template")
		assert.Contains(t, err.Error(), "Unknown")
	})
}

func TestDryRunHasNoFilesystemPollingOrCommandSideEffects(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	require.NoError(t, os.MkdirAll(dataDir, 0755))
	existingDataFile := filepath.Join(dataDir, "data.000.grb2")
	require.NoError(t, os.WriteFile(existingDataFile, []byte("original"), 0644))
	missingDataFile := filepath.Join(dataDir, "data.024.grb2")
	commandOutput := filepath.Join(dir, "command-output")
	reportFile := filepath.Join(dir, "report.json")
	dataConfigPath := writeRunPlanDataConfig(t, fmt.Sprintf(`
default: NOTFOUND
file_name: data.{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: %s
  - type: local
    level: archive
    path: %s
`, filepath.ToSlash(dataDir), filepath.ToSlash(filepath.Join(dir, "missing-archive"))))

	config := newRunPlanTestConfig(dataConfigPath)
	config.ForecastTimes = []string{"0h", "24h"}
	config.CheckInterval = "1h"
	config.DelayTime = "1h"
	config.ExecuteCommand = "echo done > " + filepath.ToSlash(commandOutput)
	config.ReportFile = reportFile

	var output bytes.Buffer
	require.NoError(t, DryRun(config, time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC), nil, &output))
	assert.Contains(t, output.String(), existingDataFile)
	assert.Contains(t, output.String(), missingDataFile)
	assert.Contains(t, output.String(), filepath.Join(dir, "missing-archive", "data.000.grb2"))

	content, err := os.ReadFile(existingDataFile)
	require.NoError(t, err)
	assert.Equal(t, "original", string(content))
	_, err = os.Stat(missingDataFile)
	assert.True(t, os.IsNotExist(err), "dry-run created a candidate data file")
	_, err = os.Stat(commandOutput)
	assert.True(t, os.IsNotExist(err), "dry-run executed the command")
	_, err = os.Stat(reportFile)
	assert.True(t, os.IsNotExist(err), "dry-run created a report")
}

const testRunPlanDataConfig = `
default: NOTFOUND
file_name: data.{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: /data
`

func newRunPlanTestConfig(dataConfigPath string) CheckerConfig {
	config := DefaultCheckerConfig()
	config.DataConfigFile = dataConfigPath
	config.ForecastTimes = []string{"0h"}
	return config
}

func writeRunPlanDataConfig(t *testing.T, content string) string {
	t.Helper()
	filePath := filepath.Join(t.TempDir(), "data-config.yaml")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
	return filePath
}
