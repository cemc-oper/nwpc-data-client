package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	e2eDemoDataType   = "e2e_demo/current/grib2/orig"
	e2eDemoConfigPath = e2eDemoDataType + ".yaml"
)

// localArgs builds the common argument list for the "local" command against
// the e2e_demo fixture config rooted at configDir.
func localArgs(configDir string, extra ...string) []string {
	args := []string{
		"local",
		"--data-type=" + e2eDemoDataType,
		"--data-config-dir=" + configDir,
		"--start-time=2026101200",
		"--forecast-time=24h",
	}
	return append(args, extra...)
}

// TestLocalFileExists is scenario E1: the data file exists under the
// temporary data root, so the CLI prints its absolute path.
func TestLocalFileExists(t *testing.T) {
	rootDir := t.TempDir()
	expected := filepath.Join(rootDir, "runtime", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(expected), 0755))
	require.NoError(t, os.WriteFile(expected, []byte("data"), 0644))

	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)

	got := runCLI(t, localArgs(configDir)...)

	assert.Equal(t, expected, got)
}

// TestLocalFileNotFound is scenario E2: no data file exists, so the CLI
// prints the config's default value.
func TestLocalFileNotFound(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)

	got := runCLI(t, localArgs(configDir)...)

	assert.Equal(t, "NOTFOUND", got)
}

// TestLocalLocationLevelFilter is scenario E3: the data file exists only
// under the archive path, and searching with --location-level=runtime skips
// it, so the CLI prints the default value.
func TestLocalLocationLevelFilter(t *testing.T) {
	rootDir := t.TempDir()
	archiveFile := filepath.Join(rootDir, "archive", "2026101200", "data.2026101200024.grb2")
	require.NoError(t, os.MkdirAll(filepath.Dir(archiveFile), 0755))
	require.NoError(t, os.WriteFile(archiveFile, []byte("data"), 0644))

	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, rootDir)

	got := runCLI(t, localArgs(configDir, "--location-level=runtime")...)
	assert.Equal(t, "NOTFOUND", got)

	// Without the filter the archive path is found.
	got = runCLI(t, localArgs(configDir)...)
	assert.Equal(t, archiveFile, got)
}

// TestLocalShowTypesEmbedded is scenario E4: --show-types without a config
// dir lists the embedded data types.
func TestLocalShowTypesEmbedded(t *testing.T) {
	got := runCLI(t, "local", "--show-types")

	assert.Contains(t, got, "cma_gfs_gmf/current/grib2/orig")
}

// TestLocalShowTypesConfigDir is scenario E5: --show-types with a config dir
// lists the fixture data type from that directory.
func TestLocalShowTypesConfigDir(t *testing.T) {
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, t.TempDir())

	got := runCLI(t, "local", "--show-types", "--data-config-dir="+configDir)

	// The CLI prints relative paths with OS separators; normalize for Windows.
	assert.Contains(t, filepath.ToSlash(got), e2eDemoDataType)
}
