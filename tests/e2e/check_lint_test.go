package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckLintValidFile(t *testing.T) {
	dir := t.TempDir()
	configPath := writeCheckLintConfig(t, dir, "checker.yaml", `
data_type: some/type
forecast_times: []
check_interval_jitter: 0
command_timeout: 0
`)

	output := runCLI(t, "check", "lint", configPath)
	assert.Contains(t, output, "check lint passed: 1 file(s)")
}

func TestCheckLintEmptyDirectory(t *testing.T) {
	output := runCLI(t, "check", "lint", t.TempDir())
	assert.Contains(t, output, "check lint passed: 0 file(s)")
}

func TestCheckLintAggregatesErrors(t *testing.T) {
	dir := t.TempDir()
	writeCheckLintConfig(t, dir, "valid.yaml", "data_type: some/type\n")
	writeCheckLintConfig(t, dir, "unknown.yaml", "data_type: some/type\nunknown_field: true\n")
	writeCheckLintConfig(t, dir, "forecast.yaml", "data_type: some/type\nforecast_times: [invalid]\n")
	writeCheckLintConfig(t, dir, "command.yaml", "data_type: some/type\nexecute_command: echo {{\n")

	output, err := runCheckLint(t, dir)
	require.Error(t, err)
	assert.Contains(t, output, "unknown.yaml")
	assert.Contains(t, output, "unknown_field")
	assert.Contains(t, output, "forecast.yaml")
	assert.Contains(t, output, "forecast time")
	assert.Contains(t, output, "command.yaml")
	assert.Contains(t, output, "parse execute command")
}

func TestCheckLintMissingPath(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing")
	output, err := runCheckLint(t, missingPath)
	require.Error(t, err)
	assert.Contains(t, output, "access checker path")
	assert.Contains(t, output, filepath.Base(missingPath))
}

func TestCheckLintHasNoRuntimeSideEffects(t *testing.T) {
	dir := t.TempDir()
	commandOutput := filepath.Join(dir, "command-output")
	reportFile := filepath.Join(dir, "report.json")
	configPath := writeCheckLintConfig(t, dir, "checker.yaml", fmt.Sprintf(`
data_type: some/type
data_config_dir: %q
forecast_times: []
execute_command: %q
report_file: %q
`, filepath.ToSlash(filepath.Join(dir, "missing-data-config")),
		"echo should-not-run > "+filepath.ToSlash(commandOutput), filepath.ToSlash(reportFile)))

	runCLI(t, "check", "lint", configPath)
	_, err := os.Stat(commandOutput)
	assert.True(t, os.IsNotExist(err), "check lint executed the configured command")
	_, err = os.Stat(reportFile)
	assert.True(t, os.IsNotExist(err), "check lint created a report file")
}

func runCheckLint(t *testing.T, args ...string) (string, error) {
	t.Helper()
	binPath := getBinaryPath(t)
	cliArgs := append([]string{"check", "lint"}, args...)
	cmd := exec.Command(binPath, cliArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func writeCheckLintConfig(t *testing.T, dir, name, content string) string {
	t.Helper()
	filePath := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
	return filePath
}
