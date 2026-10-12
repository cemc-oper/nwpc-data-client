// Package testutil provides helpers shared by the e2e and integration test
// suites: locating the nwpc_data_client binary and the config directory, and
// running the CLI binary. Suite-specific skip logic (missing binary, missing
// operational data) stays in the suites themselves.
package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// RepoRoot returns the repository root derived from this file's location.
func RepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "failed to get current file path")
	return filepath.Join(filepath.Dir(file), "..", "..")
}

// BinaryPath returns the path to the nwpc_data_client binary. It honors the
// NWPC_DATA_CLIENT_PROGRAM environment variable and otherwise defaults to
// <repo>/bin/nwpc_data_client (.exe on Windows). It does not check that the
// binary exists; callers decide how to handle a missing binary.
func BinaryPath(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("NWPC_DATA_CLIENT_PROGRAM"); p != "" {
		return p
	}
	binPath := filepath.Join(RepoRoot(t), "bin", "nwpc_data_client")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}
	return binPath
}

// ConfigDir returns the absolute path of a config directory. It honors the
// NWPC_DATA_CLIENT_CONFIG_DIR environment variable and otherwise defaults to
// the repo's conf directory; subdirs are appended (e.g. "local").
func ConfigDir(t *testing.T, subdirs ...string) string {
	t.Helper()
	base := os.Getenv("NWPC_DATA_CLIENT_CONFIG_DIR")
	if base == "" {
		base = filepath.Join(RepoRoot(t), "conf")
	}
	parts := append([]string{base}, subdirs...)
	return filepath.Join(parts...)
}

// RunCLI runs the nwpc_data_client binary with the given arguments, requires
// a zero exit code, and returns the trimmed combined output.
func RunCLI(t *testing.T, args ...string) string {
	t.Helper()
	binPath := BinaryPath(t)
	cmd := exec.Command(binPath, args...)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "CLI failed: %s", string(output))
	return strings.TrimSpace(string(output))
}
