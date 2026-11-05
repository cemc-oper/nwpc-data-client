package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLintCheckerFile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantError string
	}{
		{
			name: "valid legacy config uses defaults",
			content: `
data_type: some/type
location_levels: runtime
max_check_count: 100
execute_command: echo ready
`,
		},
		{
			name: "empty forecast times allowed for runtime stdin",
			content: `
data_type: some/type
forecast_times: []
`,
		},
		{
			name: "explicit zero values",
			content: `
data_type: some/type
check_interval_jitter: 0
command_timeout: 0
`,
		},
		{
			name: "unknown field",
			content: `
data_type: some/type
check_intervel: 10s
`,
			wantError: "check_intervel",
		},
		{
			name: "invalid forecast time",
			content: `
data_type: some/type
forecast_times: [0h, invalid]
`,
			wantError: "forecast time",
		},
		{
			name: "invalid command template",
			content: `
data_type: some/type
execute_command: echo {{.FilePath
`,
			wantError: "parse execute command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := writeCheckerLintFile(t, t.TempDir(), "checker.yaml", tt.content)
			err := LintCheckerFile(filePath)
			if tt.wantError == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError)
				assert.Contains(t, err.Error(), filePath)
			}
		})
	}
}

func TestLintCheckerPathsAggregatesFileErrors(t *testing.T) {
	dir := t.TempDir()
	writeCheckerLintFile(t, dir, "valid.yaml", "data_type: some/type\n")
	writeCheckerLintFile(t, dir, "unknown.yaml", "data_type: some/type\nunknown_field: true\n")
	writeCheckerLintFile(t, dir, "forecast.yaml", "data_type: some/type\nforecast_times: [invalid]\n")
	writeCheckerLintFile(t, dir, "command.yaml", "data_type: some/type\nexecute_command: echo {{\n")

	fileCount, err := LintCheckerPaths([]string{dir})
	require.Error(t, err)
	assert.Equal(t, 4, fileCount)
	assert.Contains(t, err.Error(), "unknown.yaml")
	assert.Contains(t, err.Error(), "unknown_field")
	assert.Contains(t, err.Error(), "forecast.yaml")
	assert.Contains(t, err.Error(), "forecast time")
	assert.Contains(t, err.Error(), "command.yaml")
	assert.Contains(t, err.Error(), "parse execute command")
}

func TestLintCheckerPathsSupportsFilesAndMultipleDirectories(t *testing.T) {
	firstDir := t.TempDir()
	nestedDir := filepath.Join(firstDir, "nested")
	writeCheckerLintFile(t, nestedDir, "nested.yaml", "data_type: some/type\n")
	writeCheckerLintFile(t, nestedDir, "ignored.yml", "not: a checker config\n")
	writeCheckerLintFile(t, nestedDir, "ignored.txt", "not: a checker config\n")

	secondDir := t.TempDir()
	writeCheckerLintFile(t, secondDir, "second.yaml", "data_type: some/type\n")
	directFile := writeCheckerLintFile(t, t.TempDir(), "direct.conf", "data_type: some/type\n")

	fileCount, err := LintCheckerPaths([]string{firstDir, secondDir, directFile, firstDir})
	require.NoError(t, err)
	assert.Equal(t, 3, fileCount)
}

func TestLintCheckerPathsEmptyDirectory(t *testing.T) {
	fileCount, err := LintCheckerPaths([]string{t.TempDir()})
	require.NoError(t, err)
	assert.Zero(t, fileCount)
}

func TestLintCheckerPathsAggregatesTraversalAndFileErrors(t *testing.T) {
	dir := t.TempDir()
	writeCheckerLintFile(t, dir, "invalid.yaml", "data_type: some/type\nmax_check_count: 0\n")
	missingPath := filepath.Join(t.TempDir(), "missing")

	fileCount, err := LintCheckerPaths([]string{missingPath, dir})
	require.Error(t, err)
	assert.Equal(t, 1, fileCount)
	assert.Contains(t, err.Error(), missingPath)
	assert.Contains(t, err.Error(), "invalid.yaml")
	assert.Contains(t, err.Error(), "max_check_count")
}

func TestLintCheckerFileHasNoRuntimeSideEffects(t *testing.T) {
	dir := t.TempDir()
	missingDataConfigDir := filepath.Join(dir, "missing-data-config")
	commandOutput := filepath.Join(dir, "command-output")
	reportFile := filepath.Join(dir, "report.json")
	content := fmt.Sprintf(`
data_type: some/type
data_config_dir: %q
forecast_times: []
execute_command: %q
report_file: %q
`, filepath.ToSlash(missingDataConfigDir),
		"echo should-not-run > "+filepath.ToSlash(commandOutput), filepath.ToSlash(reportFile))
	filePath := writeCheckerLintFile(t, dir, "checker.yaml", content)

	require.NoError(t, LintCheckerFile(filePath))
	_, err := os.Stat(commandOutput)
	assert.True(t, os.IsNotExist(err), "lint executed the configured command")
	_, err = os.Stat(reportFile)
	assert.True(t, os.IsNotExist(err), "lint created a report file")
}

func writeCheckerLintFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0755))
	filePath := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
	return filePath
}
