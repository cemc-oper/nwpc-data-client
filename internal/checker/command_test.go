package checker

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCommandTemplatesSingle(t *testing.T) {
	templates, err := buildCommandTemplates("echo {{.FilePath}}", nil)
	require.NoError(t, err)
	assert.Len(t, templates, 1)
}

func TestBuildCommandTemplatesList(t *testing.T) {
	templates, err := buildCommandTemplates("", []string{"echo {{.FilePath}}", "cat {{.FilePath}}"})
	require.NoError(t, err)
	assert.Len(t, templates, 2)
}

func TestBuildCommandTemplatesMutualExclusion(t *testing.T) {
	_, err := buildCommandTemplates("echo {{.FilePath}}", []string{"cat {{.FilePath}}"})
	assert.Error(t, err)
}

func TestBuildCommandTemplatesEmpty(t *testing.T) {
	templates, err := buildCommandTemplates("", nil)
	require.NoError(t, err)
	assert.Empty(t, templates)
}

func TestBuildCommandTemplatesFilePathTemplate(t *testing.T) {
	templates, err := buildCommandTemplates("echo {{.FilePath}}", nil)

	assert.NoError(t, err)
	assert.Len(t, templates, 1)

	var builder strings.Builder
	err = templates[0].Execute(&builder, CheckerTemplateVariable{FilePath: "/some/data.grb2"})
	assert.NoError(t, err)
	assert.Equal(t, "echo /some/data.grb2", builder.String())
}

func TestBuildCommandTemplatesInvalidTemplate(t *testing.T) {
	_, err := buildCommandTemplates("echo {{.FilePath", nil)

	assert.Error(t, err)
}
