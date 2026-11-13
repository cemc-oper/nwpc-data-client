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

func TestFetchSIGTERMExitsThreeAndPublishesNothing(t *testing.T) {
	rootDir := t.TempDir()
	source := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(source), 0755))
	require.NoError(t, os.WriteFile(source, []byte("data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	targetDir := filepath.Join(rootDir, "input")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, getBinaryPath(t), "fetch",
		"--data-type="+e2eDemoDataType,
		"--data-config-dir="+configDir,
		"--start-time=2026101200",
		"--forecast-time=24h",
		"--target-dir="+targetDir,
		"--stable-interval=1h",
	)
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	require.NoError(t, cmd.Start())

	var output strings.Builder
	scanner := bufio.NewScanner(stderr)
	signalled := false
	for scanner.Scan() {
		line := scanner.Text()
		output.WriteString(line + "\n")
		if !signalled && strings.Contains(line, "checking local source stability") {
			signalled = true
			require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
		}
	}
	require.NoError(t, scanner.Err())
	err = cmd.Wait()
	require.NotEqual(t, context.DeadlineExceeded, ctx.Err(), "fetch did not exit after SIGTERM: %s", output.String())
	require.True(t, signalled, "fetch never reached stability observation: %s", output.String())

	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 3, exitErr.ExitCode())
	assert.Contains(t, output.String(), "fetch cancelled")
	assert.Empty(t, stdout.String())
	_, statErr := os.Stat(filepath.Join(targetDir, filepath.Base(source)))
	assert.True(t, os.IsNotExist(statErr))
	matches, globErr := filepath.Glob(filepath.Join(targetDir, "*.nwpc-fetch-partial"))
	require.NoError(t, globErr)
	assert.Empty(t, matches)
}
