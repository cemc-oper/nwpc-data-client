package e2e

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// rootPlaceholder is the placeholder in testdata config templates that is
// replaced with the absolute path of the test's temporary data root.
const rootPlaceholder = "%ROOT%"

// getBinaryPath returns the path to the nwpc_data_client binary. It honors
// the NWPC_DATA_CLIENT_PROGRAM environment variable and otherwise defaults to
// <repo>/bin/nwpc_data_client (.exe on Windows). The test is skipped when the
// default binary does not exist, so plain "go test ./..." keeps working in
// environments where the binary has not been built.
func getBinaryPath(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("NWPC_DATA_CLIENT_PROGRAM"); p != "" {
		return p
	}
	repoRoot := getRepoRoot(t)
	binPath := filepath.Join(repoRoot, "bin", "nwpc_data_client")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}
	if _, err := os.Stat(binPath); err != nil {
		t.Skipf("nwpc_data_client binary not found: %s (build it first, e.g. make)", binPath)
	}
	return binPath
}

// getConfigDir returns the absolute path of the repo's conf directory. It
// honors the NWPC_DATA_CLIENT_CONFIG_DIR environment variable.
func getConfigDir(t *testing.T) string {
	t.Helper()
	if d := os.Getenv("NWPC_DATA_CLIENT_CONFIG_DIR"); d != "" {
		return d
	}
	return filepath.Join(getRepoRoot(t), "conf")
}

// getRepoRoot returns the repository root derived from this file's location.
func getRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "failed to get current file path")
	return filepath.Join(filepath.Dir(file), "..", "..")
}

// writeConfigWithRoot reads the config template at the given slash-separated
// relative path under testdata/configs, replaces the %ROOT% placeholder with
// rootDir, and writes the result into a fresh temporary config dir preserving
// the relative layout. It returns the config dir to pass to
// --data-config-dir.
func writeConfigWithRoot(t *testing.T, configRelPath string, rootDir string) string {
	t.Helper()
	templatePath := filepath.Join(getRepoRoot(t), "tests", "e2e", "testdata", "configs",
		filepath.FromSlash(configRelPath))
	templateContent, err := os.ReadFile(templatePath)
	require.NoError(t, err, "read config template failed: %s", templatePath)

	// Use forward slashes so the value stays a valid YAML plain scalar on
	// every platform (backslashes would be escape characters in YAML).
	content := strings.ReplaceAll(string(templateContent), rootPlaceholder, filepath.ToSlash(rootDir))

	configDir := t.TempDir()
	configPath := filepath.Join(configDir, filepath.FromSlash(configRelPath))
	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0644))
	return configDir
}
