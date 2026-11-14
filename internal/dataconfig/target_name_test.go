package dataconfig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func targetNameTestVariable() ConfigTemplateVariable {
	startTime := time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC)
	return GenerateConfigTemplateVariableWithStorageBase(startTime, 90*time.Minute, "014", "/CMADAAS")
}

func TestRenderTargetNameEmptyTemplateKeepsSourceName(t *testing.T) {
	name, err := RenderTargetName("", targetNameTestVariable(), "data.grb2")
	require.NoError(t, err)
	assert.Equal(t, "data.grb2", name)
}

func TestRenderTargetNamePlainString(t *testing.T) {
	name, err := RenderTargetName("input.grb2", targetNameTestVariable(), "data.grb2")
	require.NoError(t, err)
	assert.Equal(t, "input.grb2", name)
}

func TestRenderTargetNameVariables(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "start time and forecast components",
			template: "gfs_{{.Year}}{{.Month}}{{.Day}}{{.Hour}}_{{.ForecastHour}}.grb2",
			expected: "gfs_2026062400_001.grb2",
		},
		{
			name:     "forecast minute and member",
			template: "ens_{{.Member}}_{{.ForecastHour}}{{.ForecastMinute}}.grb2",
			expected: "ens_014_00130.grb2",
		},
		{
			name:     "source name",
			template: "pre_{{.SourceName}}",
			expected: "pre_data.grb2",
		},
		{
			name:     "template functions",
			template: "{{generateStartTime .StartTime -3 | getYear}}{{.ForecastTime | getForecastHour | printf \"%03d\"}}.grb2",
			expected: "2026001.grb2",
		},
		{
			name:     "storage base",
			template: "{{.StorageBase}}",
			expected: "/CMADAAS",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			name, err := RenderTargetName(testCase.template, targetNameTestVariable(), "data.grb2")
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, name)
		})
	}
}

func TestRenderTargetNameInvalidTemplate(t *testing.T) {
	_, err := RenderTargetName("{{.Year", targetNameTestVariable(), "data.grb2")
	assert.ErrorContains(t, err, "parse target name template")
}

func TestRenderTargetNameUnknownField(t *testing.T) {
	_, err := RenderTargetName("{{.NoSuchField}}", targetNameTestVariable(), "data.grb2")
	assert.ErrorContains(t, err, "render target name template")
}
