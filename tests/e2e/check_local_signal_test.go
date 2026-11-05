//go:build !windows

package e2e

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, getBinaryPath(t), "check", "local", "2026101200",
		"--checker-config="+checkerConfigPath)
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	require.NoError(t, cmd.Start())

	// Wait until signal.NotifyContext is installed and Run has reached its
	// first cancellable wait, so the test does not race process startup.
	var output strings.Builder
	scanner := bufio.NewScanner(stderr)
	signalled := false
	for scanner.Scan() {
		line := scanner.Text()
		output.WriteString(line + "\n")
		if !signalled && strings.Contains(line, "startup jitter") {
			signalled = true
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
	require.True(t, signalled, "checker never reached startup wait: %s", output.String())

	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "check local should exit non-zero after SIGTERM: %s", output.String())
	assert.Equal(t, 3, exitErr.ExitCode())
	assert.Contains(t, output.String(), "cancelled")
	assert.NotContains(t, output.String(), "Usage:")
}
