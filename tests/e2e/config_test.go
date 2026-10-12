package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigLint is part of scenario E6: linting a config dir holding the
// fixture config exits 0.
func TestConfigLint(t *testing.T) {
	configDir := writeConfigWithRoot(t, e2eDemoConfigPath, t.TempDir())

	runCLI(t, "config", "lint", configDir)
}

// TestConfigMigrate is part of scenario E6: a v2 config (single-brace
// delimiters) is migrated to the current format into --output-dir.
func TestConfigMigrate(t *testing.T) {
	v2Config := `default: NOTFOUND
file_name: data.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2
paths:
  - type: local
    level: runtime
    path: /some/dir/{.Year}{.Month}{.Day}{.Hour}
`
	configDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "demo.yaml"), []byte(v2Config), 0644))

	outputDir := filepath.Join(t.TempDir(), "out")
	runCLI(t, "config", "migrate", configDir, "--output-dir="+outputDir)

	migrated, err := os.ReadFile(filepath.Join(outputDir, "demo.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(migrated), "{{.Year}}")
}

// TestConfigExport is part of scenario E6: the embedded configs are exported
// byte-identical to the conf/ tree.
func TestConfigExport(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "exported")
	runCLI(t, "config", "export", outputDir)

	relPath := filepath.Join("local", "cma_gfs_gmf", "current", "grib2", "orig.yaml")
	exported, err := os.ReadFile(filepath.Join(outputDir, relPath))
	require.NoError(t, err)
	expected, err := os.ReadFile(filepath.Join(getConfigDir(t), relPath))
	require.NoError(t, err)
	assert.Equal(t, expected, exported)
}
