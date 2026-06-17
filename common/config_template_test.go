package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateConfigTemplateVariable(t *testing.T) {
	tests := []struct {
		name         string
		startTime    time.Time
		forecastTime time.Duration
		member       string
		expected     ConfigTemplateVariable
	}{
		{
			"Test 1",
			time.Date(2025, 5, 9, 0, 0, 0, 0, time.UTC),
			24 * time.Hour,
			"",
			ConfigTemplateVariable{
				StartTime:      time.Date(2025, 5, 9, 0, 0, 0, 0, time.UTC),
				ForecastTime:   24 * time.Hour,
				Year:           "2025",
				Month:          "05",
				Day:            "09",
				Hour:           "00",
				ForecastHour:   "024",
				ForecastMinute: "00",
			},
		},
		{
			"Test 2",
			time.Date(2025, 5, 9, 12, 0, 0, 0, time.UTC),
			120 * time.Hour,
			"003",
			ConfigTemplateVariable{
				StartTime:      time.Date(2025, 5, 9, 12, 0, 0, 0, time.UTC),
				ForecastTime:   120 * time.Hour,
				Year:           "2025",
				Month:          "05",
				Day:            "09",
				Hour:           "12",
				ForecastHour:   "120",
				ForecastMinute: "00",
				Member:         "003",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateVariable := GenerateConfigTemplateVariable(tt.startTime, tt.forecastTime, tt.member)
			assert.Equal(t, tt.expected, templateVariable)
		})
	}
}
