package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runLocal runs "nwpc_data_client local" with the given extra arguments and
// returns its trimmed stdout.
func runLocal(t *testing.T, args ...string) string {
	t.Helper()
	binPath := getBinaryPath(t)
	cmd := exec.Command(binPath, append([]string{"local"}, args...)...)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "CLI failed: %s", string(output))
	return strings.TrimSpace(string(output))
}

// TestLocalFileExists is scenario E1: the data file exists under the
// temporary data root, so the CLI prints its absolute path.
func TestLocalFileExists(t *testing.T) {
	dataType := "e2e_demo/current/grib2/orig"
	configRelPath := dataType + ".yaml"

	rootDir := t.TempDir()
	expected := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(expected), 0755))
	require.NoError(t, os.WriteFile(expected, []byte("data"), 0644))

	configDir := writeConfigWithRoot(t, configRelPath, rootDir)

	got := runLocal(t,
		"--data-type="+dataType,
		"--data-config-dir="+configDir,
		"--start-time=2026101200",
		"--forecast-time=24h",
	)

	assert.Equal(t, expected, got)
}
