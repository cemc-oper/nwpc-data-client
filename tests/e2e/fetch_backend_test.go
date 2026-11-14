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

func TestFetchCmafsRequiresBackendCLI(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "remote.yaml")
	content := `_config_version: 3
file_name: data.grb2
paths:
  - type: cmafs
    level: storage
    path: /sstorage/COMMONDATA
`
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0644))
	cmd := exec.Command(getBinaryPath(t), "fetch",
		"--data-config-file="+configPath,
		"--start-time=2026101200",
		"--forecast-time=24h",
		"--target-dir="+filepath.Join(dir, "input"),
	)
	output, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "output: %s", output)
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Contains(t, string(output), "backend not configured")
}

func TestFetchExplicitBackendValidatedForLocalSourceCLI(t *testing.T) {
	rootDir := t.TempDir()
	source := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(source), 0755))
	require.NoError(t, os.WriteFile(source, []byte("data"), 0644))
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)
	backendPath := filepath.Join(rootDir, "backend.yaml")
	require.NoError(t, os.WriteFile(backendPath, []byte("version: 1\ncmafs: {}\n"), 0644))
	targetDir := filepath.Join(rootDir, "input")

	cmd := exec.Command(getBinaryPath(t), fetchArgs(configDir, targetDir, "--backend-config="+backendPath)...)
	output, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "output: %s", output)
	assert.Contains(t, string(output), "host is required")
	_, statErr := os.Stat(filepath.Join(targetDir, filepath.Base(source)))
	assert.True(t, os.IsNotExist(statErr))
}

func TestConfigLintBackendCLI(t *testing.T) {
	dir := t.TempDir()
	backendPath := filepath.Join(dir, "backend.yaml")
	knownHostsPath := filepath.Join(dir, "known_hosts")
	content := fmt.Sprintf(`version: 1
cmafs:
  host: 10.40.150.126
  known_hosts_file: %q
`, knownHostsPath)
	require.NoError(t, os.WriteFile(backendPath, []byte(content), 0644))
	output := runCLI(t, "config", "lint", "--backend-config="+backendPath)
	assert.Contains(t, output, "backend config lint passed")
}
