package dataconfig

import (
	"fmt"
	"strings"
	"text/template"
)

// TargetNameTemplateVariable extends the data config template variables with
// the basename of the selected source file. It backs "fetch --target-name".
type TargetNameTemplateVariable struct {
	ConfigTemplateVariable
	// SourceName is the basename of the selected source file.
	SourceName string
}

// RenderTargetName renders a target file name template with the same variables
// and functions as a data config, plus {{.SourceName}}. An empty template keeps
// the source basename.
func RenderTargetName(nameTemplate string, vars ConfigTemplateVariable, sourceName string) (string, error) {
	if nameTemplate == "" {
		return sourceName, nil
	}

	nameTemplateParsed, err := template.New("targetName").Funcs(ConfigTemplateFuncMap()).Parse(nameTemplate)
	if err != nil {
		return "", fmt.Errorf("parse target name template %q failed: %w", nameTemplate, err)
	}

	variable := TargetNameTemplateVariable{
		ConfigTemplateVariable: vars,
		SourceName:             sourceName,
	}
	var builder strings.Builder
	if err := nameTemplateParsed.Execute(&builder, variable); err != nil {
		return "", fmt.Errorf("render target name template %q failed: %w", nameTemplate, err)
	}
	return builder.String(), nil
}
