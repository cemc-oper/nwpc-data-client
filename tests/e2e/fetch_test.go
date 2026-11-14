package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fetchArgs(configDir string, targetDir string, extra ...string) []string {
	args := []string{
		"fetch",
		"--data-type=" + e2eDemoDataType,
		"--data-config-dir=" + configDir,
		"--start-time=2026101200",
		"--forecast-time=24h",
		"--target-dir=" + targetDir,
		"--stable-interval=1ms",
	}
	return append(args, extra...)
}

func TestFetchCopyCLI(t *testing.T) {
	rootDir := t.TempDir()
	source := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(source), 0755))
	require.NoError(t, os.WriteFile(source, []byte("fetch data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	targetDir := filepath.Join(rootDir, "work", "input")

	cmd := exec.Command(getBinaryPath(t), fetchArgs(configDir, targetDir)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	require.NoError(t, cmd.Run(), "fetch failed: %s", stderr.String())

	target := filepath.Join(targetDir, filepath.Base(source))
	assert.Equal(t, target, strings.TrimSpace(stdout.String()))
	assert.Contains(t, stderr.String(), "fetch copying")
	content, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, []byte("fetch data"), content)
}

func TestFetchNotFoundCLI(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	targetDir := filepath.Join(rootDir, "not-created")

	cmd := exec.Command(getBinaryPath(t), fetchArgs(configDir, targetDir)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "data not found")
	_, statErr := os.Stat(targetDir)
	assert.True(t, os.IsNotExist(statErr))
}

func TestFetchTargetConflictCLI(t *testing.T) {
	rootDir := t.TempDir()
	source := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(source), 0755))
	require.NoError(t, os.WriteFile(source, []byte("source"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	targetDir := t.TempDir()
	target := filepath.Join(targetDir, filepath.Base(source))
	require.NoError(t, os.WriteFile(target, []byte("do not overwrite"), 0644))

	cmd := exec.Command(getBinaryPath(t), fetchArgs(configDir, targetDir)...)
	output, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "output: %s", output)
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Contains(t, string(output), "target conflict")
	content, readErr := os.ReadFile(target)
	require.NoError(t, readErr)
	assert.Equal(t, []byte("do not overwrite"), content)
}

func TestFetchTargetNameCLI(t *testing.T) {
	rootDir := t.TempDir()
	source := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(source), 0755))
	require.NoError(t, os.WriteFile(source, []byte("fetch data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	targetDir := filepath.Join(rootDir, "work", "input")

	args := fetchArgs(configDir, targetDir,
		"--target-name=gfs_{{.Year}}{{.Month}}{{.Day}}{{.Hour}}_{{.ForecastHour}}.grb2")
	cmd := exec.Command(getBinaryPath(t), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	require.NoError(t, cmd.Run(), "fetch failed: %s", stderr.String())

	target := filepath.Join(targetDir, "gfs_2026101200_024.grb2")
	assert.Equal(t, target, strings.TrimSpace(stdout.String()))
	content, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, []byte("fetch data"), content)
}

func TestFetchTargetNameEscapeCLI(t *testing.T) {
	rootDir := t.TempDir()
	source := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(source), 0755))
	require.NoError(t, os.WriteFile(source, []byte("fetch data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	targetDir := filepath.Join(rootDir, "work", "input")

	args := fetchArgs(configDir, targetDir, "--target-name=../escape.grb2")
	cmd := exec.Command(getBinaryPath(t), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "stderr: %s", stderr.String())
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "target name")
	_, statErr := os.Stat(filepath.Join(rootDir, "work", "escape.grb2"))
	assert.True(t, os.IsNotExist(statErr))
	_, statErr = os.Stat(targetDir)
	assert.True(t, os.IsNotExist(statErr))
}
