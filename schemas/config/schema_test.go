package config

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

// The schema patterns use only RE2-compatible constructs (no lookahead or
// backreferences), so they can be exercised directly with Go's regexp.

type schemaProperty struct {
	Pattern string `json:"pattern"`
}

type schemaDoc struct {
	Properties map[string]schemaProperty `json:"properties"`
}

func loadSchemaPattern(t *testing.T, schemaFile string, keyPath ...string) *regexp.Regexp {
	data, err := os.ReadFile(schemaFile)
	require.NoError(t, err)

	var node map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &node))
	var raw json.RawMessage
	for _, key := range keyPath {
		var ok bool
		raw, ok = node[key]
		require.True(t, ok, "schema %s missing key %q", schemaFile, key)
		require.NoError(t, json.Unmarshal(raw, &node))
	}
	var prop schemaProperty
	require.NoError(t, json.Unmarshal(raw, &prop))
	require.NotEmpty(t, prop.Pattern, "schema %s key %v has no pattern", schemaFile, keyPath)
	return regexp.MustCompile(prop.Pattern)
}

func TestV1SchemaFileNamePattern(t *testing.T) {
	pattern := loadSchemaPattern(t, "data-config-v1.schema.json", "properties", "file_name")

	cases := []struct {
		value string
		want  bool
	}{
		{"gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2", true},
		{"a.{.Year4DV}{.Month1HR}{.Forecast1HR}.grb2", true},
		{"a.{generateStartTime .StartTime -3 | getYear}.grb2", true}, // functions allowed in file_name
		{"a.{.StartTime | getYear}.grb2", true},                      // raw time objects usable with functions
		{"a.{.ForecastHour}.grb2", false},                            // v2 variable name
		{"a.{.Member}.grb2", false},                                  // v3-only variable
		{"a.{noSuchFunc .Year}.grb2", false},                         // unregistered function
	}
	for _, c := range cases {
		require.Equal(t, c.want, pattern.MatchString(c.value), "file_name pattern on %q", c.value)
	}
}

func TestV1SchemaPathPattern(t *testing.T) {
	pattern := loadSchemaPattern(t, "data-config-v1.schema.json", "properties", "paths", "items", "properties", "path")

	cases := []struct {
		value string
		want  bool
	}{
		{"/g2/{.Year}{.Month}{.Day}{.Hour}/run", true},
		{"/g2/{.Hour}/fcst", true},
		{"/g2/plain/path", true},
		{"/g2/{generateStartTime .StartTime -3 | getYear}/run", false}, // no FuncMap on path in v1
		{"/g2/{.ForecastHour}/run", false},
	}
	for _, c := range cases {
		require.Equal(t, c.want, pattern.MatchString(c.value), "path pattern on %q", c.value)
	}
}

func TestV1SchemaLevelPattern(t *testing.T) {
	pattern := loadSchemaPattern(t, "data-config-v1.schema.json", "properties", "paths", "items", "properties", "level")

	for _, ok := range []string{"runtime", "archive", "storage", "runtime/archive"} {
		require.True(t, pattern.MatchString(ok), "level pattern should accept %q", ok)
	}
	for _, bad := range []string{"nowhere", "", "runtime/nowhere"} {
		require.False(t, pattern.MatchString(bad), "level pattern should reject %q", bad)
	}
}
