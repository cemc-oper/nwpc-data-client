package common

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// lintDummyStartTime and lintDummyForecastTime are fixed values used to
// execute config templates during lint. Lint only validates template syntax
// and rendered structure, so any valid time works.
var (
	lintDummyStartTime    = time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	lintDummyForecastTime = 24 * time.Hour
	lintDummyMember       = "000"
)

// ConfigTemplateFuncMap returns the custom template functions registered for
// v2/v3 whole-file config templates. Shared by ParseConfigContent and lint.
func ConfigTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"generateStartTime":    GenerateStartTime,
		"getYear":              GetYear,
		"getMonth":             GetMonth,
		"getDay":               GetDay,
		"getHour":              GetHour,
		"generateForecastTime": GenerateForecastTime,
		"getForecastHour":      GetForecastHour,
		"getForecastMinute":    GetForecastMinute,
	}
}

// TimeTemplateVariable is the v1 (client <= 0.8) field-level template
// variable struct, kept for linting and migrating v1 config files.
// Mirrors common/time_template.go at tag v0.8.1.
type TimeTemplateVariable struct {
	StartTime    time.Time
	ForecastTime time.Duration
	Year         string
	Month        string
	Day          string
	Hour         string
	Forecast     string

	Year4DV  string
	Month4DV string
	Day4DV   string
	Hour4DV  string

	Year1HR     string
	Month1HR    string
	Day1HR      string
	Hour1HR     string
	Forecast1HR string
}

// GenerateTimeTemplateVariable builds v1 template variables: start time
// components, a 3-digit forecast hour, the 4DV family (start time -3h) and
// the 1HR family (start time -1h, forecast hour +1).
func GenerateTimeTemplateVariable(startTime time.Time, forecastTime time.Duration) TimeTemplateVariable {
	time4DV := startTime.Add(-3 * time.Hour)
	time1HR := startTime.Add(-1 * time.Hour)
	return TimeTemplateVariable{
		StartTime:    startTime,
		ForecastTime: forecastTime,
		Year:         startTime.Format("2006"),
		Month:        startTime.Format("01"),
		Day:          startTime.Format("02"),
		Hour:         startTime.Format("15"),
		Forecast:     fmt.Sprintf("%03d", int(forecastTime.Hours())),

		Year4DV:  time4DV.Format("2006"),
		Month4DV: time4DV.Format("01"),
		Day4DV:   time4DV.Format("02"),
		Hour4DV:  time4DV.Format("15"),

		Year1HR:     time1HR.Format("2006"),
		Month1HR:    time1HR.Format("01"),
		Day1HR:      time1HR.Format("02"),
		Hour1HR:     time1HR.Format("15"),
		Forecast1HR: fmt.Sprintf("%03d", int(forecastTime.Hours())+1),
	}
}

// configTemplateVariableV2 is the v2 (client 0.9) template variable struct:
// identical to ConfigTemplateVariable but without Member.
type configTemplateVariableV2 struct {
	StartTime      time.Time
	ForecastTime   time.Duration
	Year           string
	Month          string
	Day            string
	Hour           string
	ForecastHour   string
	ForecastMinute string
}

// LintConfigFile lints a single data config file. The returned error carries
// the file path and the failing stage.
func LintConfigFile(configFilePath string) error {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("%s: read: %v", configFilePath, err)
	}
	if err := LintConfigContent(string(data)); err != nil {
		return fmt.Errorf("%s: %v", configFilePath, err)
	}
	return nil
}

// LintConfigContent validates raw config content of any format version:
//
//  1. template validation with the real text/template parser, using the
//     delimiters, FuncMap and variable struct matching the detected version
//     (unknown functions fail at parse, unknown variables at execute);
//  2. structural validation of the rendered result (see ValidateDataConfig).
func LintConfigContent(content string) error {
	version, err := DetectConfigVersion(content)
	if err != nil {
		return fmt.Errorf("version detect: %v", err)
	}
	if version == ConfigVersionUnknown {
		// no template markers: the content is valid for every format
		// version; lint it as the current version.
		version = CurrentConfigVersion
	}

	switch version {
	case ConfigVersion1:
		return lintConfigV1(content)
	case ConfigVersion2:
		tpVar := GenerateConfigTemplateVariable(lintDummyStartTime, lintDummyForecastTime, lintDummyMember)
		varV2 := configTemplateVariableV2{
			StartTime:      tpVar.StartTime,
			ForecastTime:   tpVar.ForecastTime,
			Year:           tpVar.Year,
			Month:          tpVar.Month,
			Day:            tpVar.Day,
			Hour:           tpVar.Hour,
			ForecastHour:   tpVar.ForecastHour,
			ForecastMinute: tpVar.ForecastMinute,
		}
		return lintWholeFileTemplate(content, "{", "}", varV2)
	case ConfigVersion3:
		tpVar := GenerateConfigTemplateVariable(lintDummyStartTime, lintDummyForecastTime, lintDummyMember)
		return lintWholeFileTemplate(content, "{{", "}}", tpVar)
	default:
		return fmt.Errorf("version detect: unsupported config format version %d", version)
	}
}

// lintWholeFileTemplate runs v2/v3 style validation: execute the whole file
// as one template, then validate the rendered YAML structure.
func lintWholeFileTemplate(content string, leftDelim string, rightDelim string, tpVar interface{}) error {
	contentTemplate, err := template.New("dataConfig").
		Funcs(ConfigTemplateFuncMap()).
		Delims(leftDelim, rightDelim).
		Option("missingkey=error").
		Parse(content)
	if err != nil {
		return fmt.Errorf("template parse: %v", err)
	}

	var configBuilder strings.Builder
	if err := contentTemplate.Execute(&configBuilder, tpVar); err != nil {
		return fmt.Errorf("template execute: %v", err)
	}

	if err := ValidateDataConfig(configBuilder.String()); err != nil {
		return fmt.Errorf("rendered structure: %v", err)
	}
	return nil
}

// lintConfigV1 runs v1 style validation: the raw file is plain YAML parsed
// first; file_name/file_names/paths[].path are then executed as field-level
// templates with "{...}" delimiters. As in client <= 0.8, custom functions
// are registered only for file_name/file_names, not for paths[].path.
func lintConfigV1(content string) error {
	dataConfig := DataConfig{}
	if err := yaml.Unmarshal([]byte(content), &dataConfig); err != nil {
		return fmt.Errorf("yaml parse: %v", err)
	}

	tpVar := GenerateTimeTemplateVariable(lintDummyStartTime, lintDummyForecastTime)

	renderField := func(name string, fieldTemplate string, withFuncs bool) (string, error) {
		tp := template.New(name).Delims("{", "}").Option("missingkey=error")
		if withFuncs {
			tp = tp.Funcs(ConfigTemplateFuncMap())
		}
		tp, err := tp.Parse(fieldTemplate)
		if err != nil {
			return "", fmt.Errorf("template parse: %v", err)
		}
		var builder strings.Builder
		if err := tp.Execute(&builder, tpVar); err != nil {
			return "", fmt.Errorf("template execute: %v", err)
		}
		return builder.String(), nil
	}

	rendered := dataConfig
	var err error
	if rendered.FileName, err = renderField("file_name", dataConfig.FileName, true); err != nil {
		return fmt.Errorf("file_name: %v", err)
	}
	rendered.FileNames = make([]string, len(dataConfig.FileNames))
	for i, fileName := range dataConfig.FileNames {
		if rendered.FileNames[i], err = renderField("file_names", fileName, true); err != nil {
			return fmt.Errorf("file_names[%d]: %v", i, err)
		}
	}
	rendered.Paths = make([]PathItem, len(dataConfig.Paths))
	for i, pathItem := range dataConfig.Paths {
		rendered.Paths[i] = pathItem
		if rendered.Paths[i].Path, err = renderField("paths[].path", pathItem.Path, false); err != nil {
			return fmt.Errorf("paths[%d].path: %v", i, err)
		}
	}

	if err := ValidateDataConfigStruct(rendered); err != nil {
		return fmt.Errorf("rendered structure: %v", err)
	}
	return nil
}
