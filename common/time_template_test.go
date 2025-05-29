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
				StartTime:    time.Date(2025, 5, 9, 0, 0, 0, 0, time.UTC),
				ForecastTime: 24 * time.Hour,
				Year:         "2025",
				Month:        "05",
				Day:          "09",
				Hour:         "00",
				Forecast:     "024",
				Year4DV:      "2025",
				Month4DV:     "05",
				Day4DV:       "08",
				Hour4DV:      "21",
				Year1HR:      "2025",
				Month1HR:     "05",
				Day1HR:       "08",
				Hour1HR:      "23",
				Forecast1HR:  "025",
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
