//go:build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fallback pairs an alternative start time with the path expected when that start time is used.
type Fallback struct {
	StartTime string
	Path      string
}

// TestCase defines one integration test case for nwpc_data_client local.
type TestCase struct {
	Name          string
	DataType      string
	LocationLevel string
	StartTime     string // YYYYMMDDHH
	ForecastTime  string // e.g. "3h", "0h"
	ConfigDir     string // "local" 或空字符串（空表示用 embedded config）
	ExpectedPath  string
	Fallbacks     []Fallback // 部分 runtime 测试先查昨天再查今天
	AlwaysSkip    bool
}

// runTestCase executes a single integration test case via the CLI binary.
func runTestCase(t *testing.T, tc TestCase) {
	t.Helper()

	if tc.AlwaysSkip {
		t.Skipf("permanently skipped: %s", tc.Name)
	}

	startTime := tc.StartTime
	expectedPath := tc.ExpectedPath
	if tc.Fallbacks != nil {
		found := false
		if _, err := os.Stat(expectedPath); err == nil {
			found = true
		} else {
			for _, fb := range tc.Fallbacks {
				if _, err := os.Stat(fb.Path); err == nil {
					expectedPath = fb.Path
					startTime = fb.StartTime
					found = true
					break
				}
			}
		}
		if !found {
			fallbacks := make([]string, 0, len(tc.Fallbacks))
			for _, fb := range tc.Fallbacks {
				fallbacks = append(fallbacks, fb.Path)
			}
			t.Skipf("data not available: %s (fallbacks: %s)", expectedPath, strings.Join(fallbacks, ", "))
		}
	} else if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Skipf("data not available: %s", expectedPath)
	}

	binPath := getBinaryPath(t)
	args := []string{
		"local",
		"--location-level=" + tc.LocationLevel,
		"--data-type=" + tc.DataType,
		"--start-time=" + startTime,
		"--forecast-time=" + tc.ForecastTime,
	}
	if tc.ConfigDir != "" {
		args = append(args, "--data-config-dir="+getConfigDir(t, tc.ConfigDir))
	}

	cmd := exec.Command(binPath, args...)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "CLI failed: %s", string(output))

	got := strings.TrimSpace(string(output))
	assert.Equal(t, expectedPath, got)
}

// getBinaryPath returns the path to the nwpc_data_client binary.
// It honors the NWPC_DATA_CLIENT_PROGRAM environment variable.
func getBinaryPath(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("NWPC_DATA_CLIENT_PROGRAM"); p != "" {
		return p
	}
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "failed to get current file path")
	repoRoot := filepath.Join(filepath.Dir(file), "..", "..")
	return filepath.Join(repoRoot, "bin", "nwpc_data_client")
}

// getConfigDir returns the absolute config directory path for a subdir (e.g. "local").
// It honors the NWPC_DATA_CLIENT_CONFIG_DIR environment variable.
func getConfigDir(t *testing.T, subdir string) string {
	t.Helper()
	base := ""
	if d := os.Getenv("NWPC_DATA_CLIENT_CONFIG_DIR"); d != "" {
		base = d
	} else {
		_, file, _, ok := runtime.Caller(0)
		require.True(t, ok, "failed to get current file path")
		repoRoot := filepath.Join(filepath.Dir(file), "..", "..")
		base = filepath.Join(repoRoot, "conf")
	}
	return filepath.Join(base, subdir)
}

// yesterday returns YYYYMMDD for the day offsetDays before today.
func yesterday(offsetDays int) string {
	return time.Now().AddDate(0, 0, -offsetDays).Format("20060102")
}

// today returns YYYYMMDD for today.
func today() string {
	return time.Now().Format("20060102")
}

// dateHour formats a YYYYMMDD date string with the given hour.
func dateHour(date string, hour string) string {
	return fmt.Sprintf("%s%s", date, hour)
}
