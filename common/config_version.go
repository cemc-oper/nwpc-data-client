package common

import (
	"fmt"
	"regexp"
	"strconv"
)

// Config format versions. The config format version is decoupled from the
// client release version and increases monotonically (1, 2, 3, ...).
//
//	v1: client <= 0.8, field-level templates, "{...}" delimiters ({.Year}, {.Forecast})
//	v2: client 0.9, whole-file template, "{...}" delimiters ({.Year}, {.ForecastHour})
//	v3: client >= 0.10, whole-file template, "{{...}}" delimiters ({{.Year}}, {{.ForecastHour}})
const (
	ConfigVersion1 = 1
	ConfigVersion2 = 2
	ConfigVersion3 = 3

	// ConfigVersionUnknown means the content carries no template markers and
	// no explicit version field, so the format version cannot be told apart.
	ConfigVersionUnknown = 0

	// CurrentConfigVersion is the config format version this binary consumes.
	CurrentConfigVersion = ConfigVersion3
)

// ConfigVersionFieldName is the reserved top-level YAML key declaring the
// config format version. It must be static text (no template actions).
const ConfigVersionFieldName = "_config_version"

var (
	// explicitVersionPattern matches a top-level "_config_version: <n>" line.
	// A regex (instead of yaml.Unmarshal) is used on purpose: v2/v3 raw files
	// are Go templates and may contain actions that are not valid YAML.
	explicitVersionPattern = regexp.MustCompile(`(?m)^` + ConfigVersionFieldName + `:\s*([0-9]+)\s*(?:#.*)?$`)

	// v1-only variables: {.Forecast} (renamed to ForecastHour in v2) and the
	// pre-built offset variables (4DV / 1HR families) removed in v2.
	v1VariablePattern = regexp.MustCompile(`\{\.(Forecast|Year4DV|Month4DV|Day4DV|Hour4DV|Year1HR|Month1HR|Day1HR|Hour1HR|Forecast1HR)\}`)

	// v2 markers: "{...}" delimited actions referencing template variables or
	// the registered custom functions.
	v2VariablePattern = regexp.MustCompile(`\{\.[A-Za-z]`)
	v2FunctionPattern = regexp.MustCompile(`\{(generateStartTime|getYear|getMonth|getDay|getHour|generateForecastTime|getForecastHour|getForecastMinute)[\s|}]`)
)

// DetectConfigVersion detects the config format version (1/2/3) of raw config
// content. Detection rules, in priority order:
//
//  1. explicit "_config_version" field;
//  2. heuristic sniffing for old files without the field:
//     contains "{{" -> v3; v1-only variables -> v1; "{.Xxx}"/"{func ...}" -> v2;
//     no template markers at all -> ConfigVersionUnknown (undetectable).
//
// An explicit field conflicting with the sniffed version, or an invalid
// explicit value (e.g. 99), is an error.
func DetectConfigVersion(content string) (int, error) {
	explicit, hasExplicit, err := parseExplicitConfigVersion(content)
	if err != nil {
		return ConfigVersionUnknown, err
	}

	sniffed := sniffConfigVersion(content)

	if hasExplicit {
		if sniffed != ConfigVersionUnknown && sniffed != explicit {
			return ConfigVersionUnknown, fmt.Errorf(
				"config format version conflict: %s declares v%d, but content looks like v%d",
				ConfigVersionFieldName, explicit, sniffed)
		}
		return explicit, nil
	}

	return sniffed, nil
}

// parseExplicitConfigVersion extracts the explicit _config_version field.
// The third return value reports whether the field is present.
func parseExplicitConfigVersion(content string) (int, bool, error) {
	matches := explicitVersionPattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return ConfigVersionUnknown, false, nil
	}
	if len(matches) > 1 {
		return ConfigVersionUnknown, false, fmt.Errorf(
			"multiple %s fields found, expected exactly one", ConfigVersionFieldName)
	}
	version, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return ConfigVersionUnknown, false, fmt.Errorf(
			"parse %s value %q failed: %v", ConfigVersionFieldName, matches[0][1], err)
	}
	if version < ConfigVersion1 || version > ConfigVersion3 {
		return ConfigVersionUnknown, false, fmt.Errorf(
			"invalid %s value %d: supported config format versions are 1, 2, 3",
			ConfigVersionFieldName, version)
	}
	return version, true, nil
}

// sniffConfigVersion guesses the format version from template markers.
// Returns ConfigVersionUnknown when no marker is found.
func sniffConfigVersion(content string) int {
	if v3MarkerPattern.MatchString(content) {
		return ConfigVersion3
	}
	if v1VariablePattern.MatchString(content) {
		return ConfigVersion1
	}
	if v2VariablePattern.MatchString(content) || v2FunctionPattern.MatchString(content) {
		return ConfigVersion2
	}
	return ConfigVersionUnknown
}

// v3MarkerPattern matches "{{" of v3 (Go template default) delimiters.
var v3MarkerPattern = regexp.MustCompile(`\{\{`)

// CheckConfigVersion verifies that the config content matches the format
// version this binary supports (CurrentConfigVersion). Content with an
// undetectable version (ConfigVersionUnknown) is accepted, since a config
// without any template is valid for every format version.
func CheckConfigVersion(content string) error {
	version, err := DetectConfigVersion(content)
	if err != nil {
		return err
	}
	if version != ConfigVersionUnknown && version != CurrentConfigVersion {
		return fmt.Errorf(
			"config format v%d is not supported: this binary requires config format v%d "+
				"(migrate v1/v2 configs by renaming template variables and replacing \"{...}\" "+
				"delimiters with \"{{...}}\", then add \"_config_version: %d\" at the top of the file; "+
				"run \"nwpc_data_client config lint\" to validate the migrated file)",
			version, CurrentConfigVersion, CurrentConfigVersion)
	}
	return nil
}
