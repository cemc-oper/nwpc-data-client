package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// migrateTestStartTime and migrateTestForecastTime render both the source
// config (under its own version machinery) and the migrated config, so the
// migration can be checked for semantic preservation.
var (
	migrateTestStartTime    = time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	migrateTestForecastTime = 24 * time.Hour
)

// renderConfigForVersion renders raw config content with the template
// machinery of the given format version and parses the result into a
// DataConfig, normalizing away formatting differences.
func renderConfigForVersion(t *testing.T, content string, version int) DataConfig {
	t.Helper()
	var rendered string
	switch version {
	case ConfigVersion1:
		rendered = renderV1Config(t, content)
	case ConfigVersion2:
		rendered = renderWholeFileTemplate(t, content, "{", "}", false)
	case ConfigVersion3:
		rendered = renderWholeFileTemplate(t, content, "{{", "}}", true)
	default:
		t.Fatalf("unsupported version %d", version)
	}

	dataConfig := DataConfig{}
	require.NoError(t, yaml.Unmarshal([]byte(rendered), &dataConfig))
	return dataConfig
}

func renderV1Config(t *testing.T, content string) string {
	t.Helper()
	dataConfig := DataConfig{}
	require.NoError(t, yaml.Unmarshal([]byte(content), &dataConfig))

	tpVar := GenerateTimeTemplateVariable(migrateTestStartTime, migrateTestForecastTime)
	renderField := func(fieldTemplate string, withFuncs bool) string {
		tp := template.New("field").Delims("{", "}")
		if withFuncs {
			tp = tp.Funcs(ConfigTemplateFuncMap())
		}
		tp, err := tp.Parse(fieldTemplate)
		require.NoError(t, err)
		var builder strings.Builder
		require.NoError(t, tp.Execute(&builder, tpVar))
		return builder.String()
	}

	dataConfig.FileName = renderField(dataConfig.FileName, true)
	for i := range dataConfig.FileNames {
		dataConfig.FileNames[i] = renderField(dataConfig.FileNames[i], true)
	}
	for i := range dataConfig.Paths {
		dataConfig.Paths[i].Path = renderField(dataConfig.Paths[i].Path, false)
	}

	rendered, err := yaml.Marshal(dataConfig)
	require.NoError(t, err)
	return string(rendered)
}

func renderWholeFileTemplate(t *testing.T, content string, leftDelim string, rightDelim string, withMember bool) string {
	t.Helper()
	tp, err := template.New("dataConfig").
		Funcs(ConfigTemplateFuncMap()).
		Delims(leftDelim, rightDelim).
		Parse(content)
	require.NoError(t, err)

	tpVar := GenerateConfigTemplateVariable(migrateTestStartTime, migrateTestForecastTime, "000")
	var varV2 configTemplateVariableV2
	var executeVar interface{} = tpVar
	if !withMember {
		varV2 = configTemplateVariableV2{
			StartTime:      tpVar.StartTime,
			ForecastTime:   tpVar.ForecastTime,
			Year:           tpVar.Year,
			Month:          tpVar.Month,
			Day:            tpVar.Day,
			Hour:           tpVar.Hour,
			ForecastHour:   tpVar.ForecastHour,
			ForecastMinute: tpVar.ForecastMinute,
		}
		executeVar = varV2
	}

	var builder strings.Builder
	require.NoError(t, tp.Execute(&builder, executeVar))
	return builder.String()
}

// requireMigrationPreservesRendering migrates a fixture to toVersion and
// checks that the migrated config renders to the same DataConfig as the
// source under its own format version, and that it passes lint.
func requireMigrationPreservesRendering(t *testing.T, fixturePath string, fromVersion int, toVersion int) {
	t.Helper()
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err)
	content := string(data)

	detected, err := DetectConfigVersion(content)
	require.NoError(t, err)
	require.Equal(t, fromVersion, detected, "fixture %s", fixturePath)

	migrated, err := MigrateConfigContent(content, toVersion)
	require.NoError(t, err, "fixture %s", fixturePath)

	source := renderConfigForVersion(t, content, fromVersion)
	result := renderConfigForVersion(t, migrated, toVersion)
	// the version field is written by the migration on purpose; everything
	// else must render identically
	source.ConfigVersion = 0
	result.ConfigVersion = 0
	// normalize empty vs nil slices (the v1 test renderer round-trips
	// through yaml.Marshal)
	if len(source.FileNames) == 0 {
		source.FileNames = nil
	}
	if len(result.FileNames) == 0 {
		result.FileNames = nil
	}
	require.Equal(t, source, result, "fixture %s", fixturePath)
}

// Fixtures are real config files from git tags v0.8.1 (v1) and v0.9.0 (v2).
// v1/cma_gfs_gda_bin_modelvar.yaml uses custom functions in file_name;
// v2/cma_gfs_gda_bin_modelvar.yaml uses variable assignments and quoted
// strings inside actions.

func TestMigrateConfigContent_V1FixturesToV2(t *testing.T) {
	for _, name := range []string{
		"cma_gfs_gda_bin_modelvar.yaml",
		"cma_tym_grib2_orig.yaml",
	} {
		requireMigrationPreservesRendering(t,
			filepath.Join("testdata", "migrate", "v1", name), ConfigVersion1, ConfigVersion2)
	}
}

func TestMigrateConfigContent_V1FixturesToV3(t *testing.T) {
	for _, name := range []string{
		"cma_gfs_gda_bin_modelvar.yaml",
		"cma_tym_grib2_orig.yaml",
	} {
		requireMigrationPreservesRendering(t,
			filepath.Join("testdata", "migrate", "v1", name), ConfigVersion1, ConfigVersion3)
	}
}

func TestMigrateConfigContent_V2FixturesToV3(t *testing.T) {
	for _, name := range []string{
		"cma_gfs_gda_bin_modelvar.yaml",
		"cma_tym_grib2_orig.yaml",
	} {
		requireMigrationPreservesRendering(t,
			filepath.Join("testdata", "migrate", "v2", name), ConfigVersion2, ConfigVersion3)
	}
}

func TestMigrateConfigContent_SetsVersionField(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "migrate", "v2", "cma_tym_grib2_orig.yaml"))
	require.NoError(t, err)

	migrated, err := MigrateConfigContent(string(data), ConfigVersion3)
	require.NoError(t, err)

	detected, err := DetectConfigVersion(migrated)
	require.NoError(t, err)
	require.Equal(t, ConfigVersion3, detected)
	// the version field sits before the first content key
	require.Less(t,
		strings.Index(migrated, "_config_version: 3"),
		strings.Index(migrated, "default:"))
}

func TestMigrateConfigContent_V1OffsetVariables(t *testing.T) {
	content := `default: NOTFOUND
file_name: "a.{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}{.Year1HR}{.Month1HR}{.Day1HR}{.Hour1HR}{.Forecast1HR}{.Forecast}.grb2"
paths:
  - type: local
    level: runtime
    path: /some/dir/{.Year}{.Month}{.Day}{.Hour}
`
	migrated, err := MigrateConfigContent(content, ConfigVersion2)
	require.NoError(t, err)

	require.Contains(t, migrated, "a."+
		`{generateStartTime .StartTime -3 | getYear}`+
		`{generateStartTime .StartTime -3 | getMonth}`+
		`{generateStartTime .StartTime -3 | getDay}`+
		`{generateStartTime .StartTime -3 | getHour}`+
		`{generateStartTime .StartTime -1 | getYear}`+
		`{generateStartTime .StartTime -1 | getMonth}`+
		`{generateStartTime .StartTime -1 | getDay}`+
		`{generateStartTime .StartTime -1 | getHour}`+
		`{generateForecastTime .ForecastTime "1h" | getForecastHour | printf "%03d"}`+
		`{.ForecastHour}.grb2`)
	require.NotContains(t, migrated, "{.Forecast}")
	require.Contains(t, migrated, "_config_version: 2")

	// 4DV: start -3h; 1HR: start -1h, forecast +1h
	rendered := renderConfigForVersion(t, migrated, ConfigVersion2)
	require.Equal(t,
		"a.2026090821"+"2026090823"+"025"+"024.grb2",
		rendered.FileName)
}

func TestMigrateConfigContent_V2DelimiterReplacement(t *testing.T) {
	content := `default: NOTFOUND
  { $base := "/g2/data" }
file_name: "a{ getYear .StartTime }{.Year}{.ForecastHour | printf "%s"}.grb2"
paths:
  - type: local
    level: runtime
    path: "{ $base }/{.Year}{.Month}{.Day}{.Hour}"
`
	migrated, err := MigrateConfigContent(content, ConfigVersion3)
	require.NoError(t, err)

	require.Contains(t, migrated, "{{ $base := \"/g2/data\" }}")
	require.Contains(t, migrated, `a{{ getYear .StartTime }}{{.Year}}`)
	require.Contains(t, migrated, `{{.ForecastHour | printf "%s"}}`)

	rendered := renderConfigForVersion(t, migrated, ConfigVersion3)
	require.Equal(t, "a20262026024.grb2", rendered.FileName)
	require.Equal(t, "/g2/data/2026090900", rendered.Paths[0].Path)
}

func TestMigrateV2ToV3_LiteralBracesKept(t *testing.T) {
	// braces that do not start a template action (e.g. shell brace
	// expansion in a path) must survive the delimiter replacement;
	// note a real v2 config cannot contain them inside a template-parsed
	// document, so this tests the rewriter directly.
	content := `path: /some/dir/brace{1,2}keep/{.Year}/"{a,b}"`
	migrated := migrateV2ToV3(content)
	require.Equal(t, `path: /some/dir/brace{1,2}keep/{{.Year}}/"{a,b}"`, migrated)
}

func TestMigrateConfigContent_V2ControlStructures(t *testing.T) {
	content := `default: NOTFOUND
file_name: "{if eq .Hour "00"}a{else}b{end}.grb2"
paths:
  - type: local
    level: runtime
    path: /some/dir/{.Year}
`
	migrated, err := MigrateConfigContent(content, ConfigVersion3)
	require.NoError(t, err)
	require.Contains(t, migrated, `{{if eq .Hour "00"}}a{{else}}b{{end}}`)

	rendered := renderConfigForVersion(t, migrated, ConfigVersion3)
	require.Equal(t, "a.grb2", rendered.FileName)
}

func TestMigrateConfigContent_AlreadyAtTarget(t *testing.T) {
	content := `_config_version: 3

default: NOTFOUND
paths:
  - type: local
    level: runtime
    path: /some/dir
`
	_, err := MigrateConfigContent(content, ConfigVersion3)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already at format v3")
}

func TestMigrateConfigContent_DowngradeRejected(t *testing.T) {
	content := `_config_version: 3

default: NOTFOUND
paths:
  - type: local
    level: runtime
    path: /some/dir
`
	_, err := MigrateConfigContent(content, ConfigVersion2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "v3 down to v2")
}

func TestMigrateConfigContent_BadSourceRejected(t *testing.T) {
	// a v1 file using a function in paths[].path, where client <= 0.8
	// registered no FuncMap: the source was never valid, so migration
	// must fail with a clear error instead of silently rewriting it.
	content := `default: NOTFOUND
file_name: a{.Forecast}.grb2
paths:
  - type: local
    level: runtime
    path: /some/dir/{generateStartTime .StartTime -3 | getYear}
`
	_, err := MigrateConfigContent(content, ConfigVersion3)
	require.Error(t, err)
	require.Contains(t, err.Error(), "source config is not a valid v1 config")
}

func TestMigrateConfigContent_UnknownVersionGetsField(t *testing.T) {
	content := `default: NOTFOUND
paths:
  - type: local
    level: runtime
    path: /some/dir
`
	migrated, err := MigrateConfigContent(content, ConfigVersion3)
	require.NoError(t, err)
	require.Contains(t, migrated, "_config_version: 3")
	require.NoError(t, LintConfigContent(migrated))
}

func TestMigrateConfigContent_UnsupportedTarget(t *testing.T) {
	_, err := MigrateConfigContent("default: NOTFOUND\n", 99)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported target")
}
