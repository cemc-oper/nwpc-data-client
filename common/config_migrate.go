package common

import (
	"fmt"
	"regexp"
	"strings"
)

// MigrateConfigContent migrates raw config content from its detected config
// format version to toVersion (ConfigVersion2 or ConfigVersion3). v1 -> v2
// renames template variables; v2 -> v3 replaces "{...}" delimiters with
// "{{...}}"; v1 -> v3 chains both steps. The migrated content gets a
// "_config_version" field and is validated with LintConfigContent before it
// is returned, so a failed migration never silently produces output.
//
// The source content is linted first: migrating a config that is already
// broken (e.g. a v1 file using functions in paths[].path, where client <= 0.8
// registered no FuncMap) is an error instead of a silent best-effort rewrite.
func MigrateConfigContent(content string, toVersion int) (string, error) {
	if toVersion != ConfigVersion2 && toVersion != ConfigVersion3 {
		return "", fmt.Errorf("unsupported target config format version %d: migration targets are v2 and v3", toVersion)
	}

	fromVersion, err := DetectConfigVersion(content)
	if err != nil {
		return "", fmt.Errorf("version detect: %v", err)
	}
	if fromVersion == toVersion {
		return "", fmt.Errorf("config is already at format v%d", toVersion)
	}
	if fromVersion > toVersion {
		return "", fmt.Errorf("cannot migrate config format v%d down to v%d", fromVersion, toVersion)
	}

	if fromVersion != ConfigVersionUnknown {
		if err := LintConfigContent(content); err != nil {
			return "", fmt.Errorf("source config is not a valid v%d config: %v", fromVersion, err)
		}
	}

	migrated := content
	if fromVersion == ConfigVersion1 {
		migrated = migrateV1ToV2(migrated)
	}
	if fromVersion != ConfigVersionUnknown && toVersion == ConfigVersion3 {
		migrated = migrateV2ToV3(migrated)
	}
	migrated = setConfigVersionField(migrated, toVersion)

	if err := LintConfigContent(migrated); err != nil {
		return "", fmt.Errorf("migrated config does not pass lint: %v", err)
	}
	return migrated, nil
}

// v1ToV2Replacements maps v1 template variables removed in client 0.9 to
// their v2 equivalents. Tokens carry both braces, so a plain string replace
// cannot mismatch ("{.Forecast}" is not a substring of "{.Forecast1HR}").
var v1ToV2Replacements = []struct {
	old string
	new string
}{
	// forecast hour +1, formatted "%03d"
	{"{.Forecast1HR}", `{generateForecastTime .ForecastTime "1h" | getForecastHour | printf "%03d"}`},
	// start time -1h
	{"{.Year1HR}", `{generateStartTime .StartTime -1 | getYear}`},
	{"{.Month1HR}", `{generateStartTime .StartTime -1 | getMonth}`},
	{"{.Day1HR}", `{generateStartTime .StartTime -1 | getDay}`},
	{"{.Hour1HR}", `{generateStartTime .StartTime -1 | getHour}`},
	// start time -3h
	{"{.Year4DV}", `{generateStartTime .StartTime -3 | getYear}`},
	{"{.Month4DV}", `{generateStartTime .StartTime -3 | getMonth}`},
	{"{.Day4DV}", `{generateStartTime .StartTime -3 | getDay}`},
	{"{.Hour4DV}", `{generateStartTime .StartTime -3 | getHour}`},
	// renamed variable
	{"{.Forecast}", `{.ForecastHour}`},
}

// migrateV1ToV2 renames the v1-only template variables to their v2
// equivalents. Variables shared by both versions ({.Year}, ...) and the
// custom functions (available in v1 file_name/file_names) are unchanged.
func migrateV1ToV2(content string) string {
	for _, replacement := range v1ToV2Replacements {
		content = strings.ReplaceAll(content, replacement.old, replacement.new)
	}
	return content
}

// v2ActionKeywords are the identifiers that may open a v2 template action
// after "{": the custom functions and the text/template builtin keywords.
// A "{" followed by anything else (e.g. shell brace expansion "{a,b}" in a
// path) is left untouched.
var v2ActionKeywords = map[string]bool{
	// custom functions (ConfigTemplateFuncMap)
	"generateStartTime": true, "getYear": true, "getMonth": true, "getDay": true,
	"getHour": true, "generateForecastTime": true, "getForecastHour": true, "getForecastMinute": true,
	// text/template builtins
	"if": true, "else": true, "end": true, "range": true, "with": true,
	"define": true, "template": true, "block": true, "break": true, "continue": true,
	"printf": true, "print": true, "println": true, "len": true,
	"and": true, "or": true, "not": true,
	"eq": true, "ne": true, "lt": true, "le": true, "gt": true, "ge": true,
	"index": true, "slice": true, "call": true,
}

// migrateV2ToV3 replaces "{...}" action delimiters with "{{...}}". Only
// actions starting with ".", "$", "-", a comment opener "/*", or a
// whitelisted function/keyword are rewritten; literal braces in paths stay
// as they are.
func migrateV2ToV3(content string) string {
	var builder strings.Builder
	builder.Grow(len(content))
	index := 0
	for index < len(content) {
		if content[index] == '{' && (index+1 >= len(content) || content[index+1] != '{') {
			if closeIndex, ok := scanV2Action(content, index); ok {
				builder.WriteString("{{")
				builder.WriteString(content[index+1 : closeIndex])
				builder.WriteString("}}")
				index = closeIndex + 1
				continue
			}
		}
		builder.WriteByte(content[index])
		index++
	}
	return builder.String()
}

// scanV2Action checks whether the "{" at openIndex starts a v2 template
// action. If so it returns the index of the matching closing "}".
func scanV2Action(content string, openIndex int) (int, bool) {
	index := openIndex + 1
	for index < len(content) && (content[index] == ' ' || content[index] == '\t') {
		index++
	}
	if index < len(content) && content[index] == '-' {
		index++
		for index < len(content) && (content[index] == ' ' || content[index] == '\t') {
			index++
		}
	}
	if index >= len(content) {
		return 0, false
	}

	first := content[index]
	isAction := false
	switch {
	case first == '.' || first == '$':
		// variable reference: "{.Year}", "{ $x := ... }", "{$x}"
		isAction = true
	case first == '/':
		// comment action: "{/* ... */}"
		isAction = index+1 < len(content) && content[index+1] == '*'
	case first == '_' || (first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z'):
		end := index
		for end < len(content) && (content[end] == '_' ||
			(content[end] >= 'a' && content[end] <= 'z') ||
			(content[end] >= 'A' && content[end] <= 'Z') ||
			(content[end] >= '0' && content[end] <= '9')) {
			end++
		}
		isAction = v2ActionKeywords[content[index:end]]
	}
	if !isAction {
		return 0, false
	}

	// find the closing "}", skipping over quoted strings
	inDoubleQuote := false
	inRawQuote := false
	escaped := false
	for index++; index < len(content); index++ {
		char := content[index]
		switch {
		case inRawQuote:
			if char == '`' {
				inRawQuote = false
			}
		case inDoubleQuote:
			if escaped {
				escaped = false
			} else if char == '\\' {
				escaped = true
			} else if char == '"' {
				inDoubleQuote = false
			}
		case char == '`':
			inRawQuote = true
		case char == '"':
			inDoubleQuote = true
		case char == '}':
			return index, true
		}
	}
	return 0, false
}

// configVersionLinePattern matches a whole top-level "_config_version:" line.
var configVersionLinePattern = regexp.MustCompile(`(?m)^` + ConfigVersionFieldName + `:[^\n]*(\n|$)`)

// setConfigVersionField replaces any existing "_config_version" line and
// inserts "_config_version: <version>" after the leading comment/blank block.
func setConfigVersionField(content string, version int) string {
	content = configVersionLinePattern.ReplaceAllString(content, "")

	versionLine := fmt.Sprintf("%s: %d\n\n", ConfigVersionFieldName, version)
	lines := strings.SplitAfter(content, "\n")
	insertIndex := 0
	for insertIndex < len(lines) {
		trimmed := strings.TrimSpace(lines[insertIndex])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			insertIndex++
			continue
		}
		break
	}

	var builder strings.Builder
	written := false
	for i, line := range lines {
		if i == insertIndex {
			builder.WriteString(versionLine)
			written = true
		}
		builder.WriteString(line)
	}
	if !written {
		// the whole file is comments/blank lines
		builder.WriteString(versionLine)
	}
	return builder.String()
}
