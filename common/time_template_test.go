package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateTimeTemplateVariable(t *testing.T) {
	tests := []struct {
		name         string
		startTime    time.Time
		forecastTime time.Duration
		expected     TimeTemplateVariable
	}{
		{
			"Test 1",
			time.Date(2025, 5, 9, 0, 0, 0, 0, time.UTC),
			24 * time.Hour,
			TimeTemplateVariable{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateVariable := GenerateTimeTemplateVariable(tt.startTime, tt.forecastTime)
			assert.Equal(t, tt.expected, templateVariable)
		})
	}
}
