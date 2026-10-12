package dataconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// v3 sample: adapted from conf/local/cma_gfs_gmf/current/grib2/orig.yaml
const v3ConfigSample = `# cma-gfs
#   grib2 orig
_config_version: 3

default: NOTFOUND

file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/data/output/grib2_orig/
`

// v2 sample: client 0.9 format, "{...}" delimiters with renamed variables.
const v2ConfigSample = `# cma-gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig/
`

// v1 sample: client <= 0.8 format, "{...}" delimiters with old variable names.
const v1ConfigSample = `# cma-gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig/
`

func TestDetectConfigVersion_ExplicitField(t *testing.T) {
	version, err := DetectConfigVersion(v3ConfigSample)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion3, version)
}

func TestDetectConfigVersion_SniffV2(t *testing.T) {
	version, err := DetectConfigVersion(v2ConfigSample)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion2, version)
}

func TestDetectConfigVersion_SniffV1(t *testing.T) {
	version, err := DetectConfigVersion(v1ConfigSample)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion1, version)
}

func TestDetectConfigVersion_SniffV3WithoutExplicitField(t *testing.T) {
	content := "default: NOTFOUND\nfile_name: a.{{.Year}}{{.ForecastHour}}.grb2\n"
	version, err := DetectConfigVersion(content)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion3, version)
}

func TestDetectConfigVersion_V1OffsetVariables(t *testing.T) {
	// 4DV/1HR variables only exist in v1.
	content := "file_name: a.{.Year4DV}{.Month1HR}{.Forecast1HR}.grb2\n"
	version, err := DetectConfigVersion(content)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion1, version)
}

func TestDetectConfigVersion_V2Function(t *testing.T) {
	content := "file_name: a.{generateStartTime .StartTime -3 | getYear}.grb2\n"
	version, err := DetectConfigVersion(content)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion2, version)
}

func TestDetectConfigVersion_PlainText(t *testing.T) {
	content := "default: NOTFOUND\nfile_name: gmf.gra.2026090900024.grb2\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/\n"
	version, err := DetectConfigVersion(content)
	require.NoError(t, err)
	require.Equal(t, ConfigVersionUnknown, version)
}

func TestDetectConfigVersion_ExplicitFieldWithoutTemplate(t *testing.T) {
	// explicit field wins even when there is nothing to sniff.
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.grb2\n"
	version, err := DetectConfigVersion(content)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion3, version)
}

func TestDetectConfigVersion_Conflict(t *testing.T) {
	content := "_config_version: 2\ndefault: NOTFOUND\nfile_name: a.{{.Year}}.grb2\n"
	_, err := DetectConfigVersion(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflict")
}

func TestDetectConfigVersion_InvalidVersionValue(t *testing.T) {
	content := "_config_version: 99\ndefault: NOTFOUND\nfile_name: a.{{.Year}}.grb2\n"
	_, err := DetectConfigVersion(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid")
}

func TestDetectConfigVersion_DuplicateField(t *testing.T) {
	content := "_config_version: 3\n_config_version: 3\ndefault: NOTFOUND\n"
	_, err := DetectConfigVersion(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "multiple")
}

func TestCheckConfigVersion(t *testing.T) {
	require.NoError(t, CheckConfigVersion(v3ConfigSample))

	// no template markers: accepted, valid for every format version.
	require.NoError(t, CheckConfigVersion("default: NOTFOUND\nfile_name: a.grb2\n"))

	err := CheckConfigVersion(v1ConfigSample)
	require.Error(t, err)
	require.Contains(t, err.Error(), "v1")
	require.Contains(t, err.Error(), "v3")

	err = CheckConfigVersion(v2ConfigSample)
	require.Error(t, err)
	require.Contains(t, err.Error(), "v2")
	require.Contains(t, err.Error(), "v3")
}

func TestLoadConfigContentFromFile_VersionMismatch(t *testing.T) {
	// a v1 format config fed to the current binary must fail with a clear
	// error naming the detected version, the required version and migration
	// guidance, instead of a low-level template execution error.
	configDir := t.TempDir()
	configFile := filepath.Join(configDir, "v1_config.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte(v1ConfigSample), 0o644))

	_, err := LoadConfigContentFromFile(configFile)
	require.Error(t, err)
	require.Contains(t, err.Error(), "config format v1")
	require.Contains(t, err.Error(), "requires config format v3")
	require.Contains(t, err.Error(), "migrate")
}

func TestLoadConfigContent_VersionMismatch(t *testing.T) {
	configDir := t.TempDir()
	dataType := "some_system/current/grib2/orig"
	configFile := filepath.Join(configDir, dataType+".yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(configFile), 0o755))
	require.NoError(t, os.WriteFile(configFile, []byte(v2ConfigSample), 0o644))

	_, err := LoadConfigContent(configDir, dataType)
	require.Error(t, err)
	require.Contains(t, err.Error(), "config format v2")
	require.Contains(t, err.Error(), "requires config format v3")
}
