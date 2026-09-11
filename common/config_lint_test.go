package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLintConfigContent_V3(t *testing.T) {
	require.NoError(t, LintConfigContent(v3ConfigSample))
}

func TestLintConfigContent_V3WithTemplateDirectives(t *testing.T) {
	// template actions may appear anywhere in a v3 file (variable assignment,
	// control structures), not only inside string values.
	content := `{{- $base := "/g2/op_post/OPER/WORKDIR" -}}
{{- $start := generateStartTime .StartTime -3 -}}
_config_version: 3
default: NOTFOUND
file_name: gmf.gra.{{getYear $start}}{{getMonth $start}}{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: {{ $base }}/{{ $.Year }}{{ $.Month }}{{ $.Day }}{{ $.Hour }}/data/
  {{- if eq $.Hour "00" }}
  - type: local
    level: archive
    path: /g3/COMMONDATA/{{ getYear $start }}/
  {{- end }}
`
	require.NoError(t, LintConfigContent(content))
}

func TestLintConfigContent_V2(t *testing.T) {
	require.NoError(t, LintConfigContent(v2ConfigSample))
}

func TestLintConfigContent_V2MemberRejected(t *testing.T) {
	// Member does not exist in the v2 variable struct.
	content := "default: NOTFOUND\nfile_name: a.{.Member}.grb2\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "template execute")
}

func TestLintConfigContent_V1(t *testing.T) {
	require.NoError(t, LintConfigContent(v1ConfigSample))
}

func TestLintConfigContent_V1FunctionInFileName(t *testing.T) {
	// functions are registered for file_name in v1. {.Forecast} makes the
	// heuristic detect this file as v1.
	content := "default: NOTFOUND\nfile_name: a.{generateStartTime .StartTime -3 | getYear}{.Forecast}.grb2\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/{.Year}/\n"
	require.NoError(t, LintConfigContent(content))
}

func TestLintConfigContent_V1FunctionInPathRejected(t *testing.T) {
	// v1 path templates have no FuncMap: functions must fail. {.Forecast}
	// makes the heuristic detect this file as v1.
	content := "default: NOTFOUND\nfile_name: a.{.Forecast}.grb2\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/{generateStartTime .StartTime -3 | getYear}/\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "paths[0].path")
	require.Contains(t, err.Error(), "template parse")
}

func TestLintConfigContent_UnknownFunction(t *testing.T) {
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.{{noSuchFunc .Year}}.grb2\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "template parse")
	require.Contains(t, err.Error(), "noSuchFunc")
}

func TestLintConfigContent_UnknownVariable(t *testing.T) {
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.{{.NoSuchVar}}.grb2\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "template execute")
	require.Contains(t, err.Error(), "NoSuchVar")
}

func TestLintConfigContent_MissingPaths(t *testing.T) {
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.{{.Year}}.grb2\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "rendered structure")
	require.Contains(t, err.Error(), "paths")
}

func TestLintConfigContent_InvalidPathItem(t *testing.T) {
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.grb2\npaths:\n  - type: remote\n    level: runtime\n    path: /some/dir/\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported type")
}

func TestLintConfigContent_InvalidLevel(t *testing.T) {
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.grb2\npaths:\n  - type: local\n    level: nowhere\n    path: /some/dir/\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid level")
}

func TestLintConfigContent_CombinedLevel(t *testing.T) {
	// "runtime/archive" is used by existing configs and must be accepted.
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.grb2\npaths:\n  - type: local\n    level: runtime/archive\n    path: /some/dir/\n"
	require.NoError(t, LintConfigContent(content))
}

func TestLintConfigContent_UnknownFieldRejected(t *testing.T) {
	content := "_config_version: 3\ndefault: NOTFOUND\nfile_name: a.grb2\ntypo_field: oops\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/\n"
	err := LintConfigContent(content)
	require.Error(t, err)
	require.Contains(t, err.Error(), "typo_field")
}

func TestLintConfigContent_PlainText(t *testing.T) {
	// no template markers: linted as the current version.
	content := "default: NOTFOUND\nfile_name: a.grb2\npaths:\n  - type: local\n    level: runtime\n    path: /some/dir/\n"
	require.NoError(t, LintConfigContent(content))
}

func TestLintConfigFile(t *testing.T) {
	configDir := t.TempDir()

	goodFile := filepath.Join(configDir, "good.yaml")
	require.NoError(t, os.WriteFile(goodFile, []byte(v3ConfigSample), 0o644))
	require.NoError(t, LintConfigFile(goodFile))

	badFile := filepath.Join(configDir, "bad.yaml")
	require.NoError(t, os.WriteFile(badFile, []byte("default: NOTFOUND\n"), 0o644))
	err := LintConfigFile(badFile)
	require.Error(t, err)
	require.Contains(t, err.Error(), badFile)
}
