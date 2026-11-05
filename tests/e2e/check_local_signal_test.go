//go:build !windows

package e2e

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runCheckLocalWithSIGTERM(t *testing.T, checkerConfigPath string, readyLog string) (string, error, time.Duration) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, getBinaryPath(t), "check", "local", "2026101200",
		"--checker-config="+checkerConfigPath)
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	require.NoError(t, cmd.Start())

	// Wait until signal.NotifyContext is installed and the requested run stage
	// has been reached, so the test does not race process startup.
	var output strings.Builder
	scanner := bufio.NewScanner(stderr)
	signalled := false
	var signalledAt time.Time
	for scanner.Scan() {
		line := scanner.Text()
		output.WriteString(line + "\n")
		if !signalled && strings.Contains(line, readyLog) {
			signalled = true
			signalledAt = time.Now()
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				_ = cmd.Process.Kill()
				t.Fatalf("send SIGTERM: %v", err)
			}
		}
	}
	require.NoError(t, scanner.Err())
	err = cmd.Wait()
	output.WriteString(stdout.String())
	require.NotEqual(t, context.DeadlineExceeded, ctx.Err(), "check local did not exit after SIGTERM: %s", output.String())
	require.True(t, signalled, "checker never reached %q: %s", readyLog, output.String())
	return output.String(), err, time.Since(signalledAt)
}

func TestCheckLocalSIGTERMExitsThree(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir: filepath.ToSlash(configDir),
		DataType:      e2eDemoDataType,
		CheckInterval: "1h",
		MaxCheckCount: 10,
		DelayTime:     "0s",
		ForecastTimes: []string{"24h"},
	})

	output, err, _ := runCheckLocalWithSIGTERM(t, checkerConfigPath, "startup jitter")

	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "check local should exit non-zero after SIGTERM: %s", output)
	assert.Equal(t, 3, exitErr.ExitCode())
	assert.Contains(t, output, "cancelled")
	assert.NotContains(t, output, "Usage:")
}

func TestCheckLocalSIGTERMCancelsRunningCommand(t *testing.T) {
	rootDir := t.TempDir()
	dataFile := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(dataFile), 0755))
	require.NoError(t, os.WriteFile(dataFile, []byte("data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	checkerConfigPath := writeCheckerConfig(t, checkerConfig{
		DataConfigDir:  filepath.ToSlash(configDir),
		DataType:       e2eDemoDataType,
		CheckInterval:  "20ms",
		MaxCheckCount:  20,
		DelayTime:      "0s",
		ForecastTimes:  []string{"24h"},
		ExecuteCommand: e2eCommandHelperCommand("wait"),
		CommandTimeout: "30m",
	})

	output, err, elapsed := runCheckLocalWithSIGTERM(t, checkerConfigPath, "running command <")

	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "check local should exit non-zero after SIGTERM: %s", output)
	assert.Equal(t, 3, exitErr.ExitCode())
	assert.Less(t, elapsed, time.Second, "running command was not cancelled promptly: %s", output)
	assert.Contains(t, output, "cancelled")
	assert.NotContains(t, output, "command timed out")
	assert.NotContains(t, output, "Usage:")
}
