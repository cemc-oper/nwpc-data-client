package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseStarTime(test *testing.T) {
	value, err := ParseStartTime("2025052912")
	assert.Nil(test, err)
	expected := time.Date(2025, 5, 29, 12, 0, 0, 0, time.UTC)
	assert.Equal(test, value, expected)
}

func TestParseStartTimeErrorFormat(test *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"YYYYMMDD", "20250529"},
		{"YYYYMMDDHHMM", "202505291200"},
		{"YYYY-MM-DD HH:MM", "2025-05-29 12:00"},
		{"YYYY-MM-DD HH:MM:SS", "2025-05-29 12:00:00"},
		{"not a date", "abcdefghij"},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			_, err := ParseStartTime(tt.value)
			assert.NotNil(test, err)
		})
	}
}

func TestParseForecastTime(test *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected time.Duration
	}{
		{"HHh", "24h", 24 * time.Hour},
		{"HHHh", "024h", 24 * time.Hour},
		{"HHHh 2", "000h", 0 * time.Hour},
		{"HHHh 3", "240h", 240 * time.Hour},
		{"HHHhMMm 1", "001h10m", 1*time.Hour + 10*time.Minute},
		{"HHHmm 2", "010h50m", 10*time.Hour + 50*time.Minute},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			result, err := ParseForecastTime(tt.value)
			assert.Nil(test, err)
			assert.Equal(test, tt.expected, result)
		})
	}
}

func TestParseForcastTimeErrorFormat(test *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"two digits", "24"},
		{"three digits 1", "000"},
		{"three digits 2", "024"},
		{"five digits", "02410"},
		{"lack m", "024h10"},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			_, err := ParseForecastTime(tt.value)
			assert.NotNil(t, err)
		})
	}
}

func TestGenerateStartTime(t *testing.T) {
	startTime := time.Date(2025, 5, 29, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		hour     int
		expected time.Time
	}{
		{"add 3 hours", 3, time.Date(2025, 5, 29, 15, 0, 0, 0, time.UTC)},
		{"subtract 2 hours", -2, time.Date(2025, 5, 29, 10, 0, 0, 0, time.UTC)},
		{"add 0 hours", 0, time.Date(2025, 5, 29, 12, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateStartTime(startTime, tt.hour)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetYear(t *testing.T) {
	startTime := time.Date(2025, 5, 29, 12, 0, 0, 0, time.UTC)
	expected := "2025"
	got := GetYear(startTime)
	assert.Equal(t, expected, got)
}

func TestGetMonth(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		expected  string
	}{
		{"January", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), "01"},
		{"November", time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), "11"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMonth(tt.startTime)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetDay(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		expected  string
	}{
		{"Single digit day", time.Date(2025, 5, 9, 0, 0, 0, 0, time.UTC), "09"},
		{"Double digit day", time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC), "29"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDay(tt.startTime)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetHour(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		expected  string
	}{
		{"Single digit hour", time.Date(2025, 5, 29, 9, 0, 0, 0, time.UTC), "09"},
		{"Double digit hour", time.Date(2025, 5, 29, 21, 0, 0, 0, time.UTC), "21"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHour(tt.startTime)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGenerateForecastTime(test *testing.T) {
	tests := []struct {
		name         string
		forecastTime time.Duration
		timeInterval string
		expected     time.Duration
	}{
		{"1 hour", 24 * time.Hour, "1h", 25 * time.Hour},
		{"2 hours", 24 * time.Hour, "2h", 26 * time.Hour},
		{"-3 hours", 24 * time.Hour, "-3h", 21 * time.Hour},
		{"10 minute", 24 * time.Hour, "10m", 24*time.Hour + 10*time.Minute},
		{"-10 minute", 24 * time.Hour, "-10m", 24*time.Hour - 10*time.Minute},
		{"1 hour 10 min", 24 * time.Hour, "1h10m", 25*time.Hour + 10*time.Minute},
		{"0 hour", 24 * time.Hour, "0h", 24 * time.Hour},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			forecastTime := GenerateForecastTime(tt.forecastTime, tt.timeInterval)
			assert.Equal(test, tt.expected, forecastTime)
		})
	}
}

// WARNING
func TestGenerateForecastTimeErrorFormat(test *testing.T) {
	tests := []struct {
		name         string
		forecastTime time.Duration
		timeInterval string
	}{
		{"lack h", 24 * time.Hour, "10"},
		{"lack m", 24 * time.Hour, "abc"},
	}
	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			forecastTime := GenerateForecastTime(tt.forecastTime, tt.timeInterval)
			assert.Equal(test, tt.forecastTime, forecastTime)
		})
	}
}

func TestGetForecastHour(test *testing.T) {
	tests := []struct {
		name         string
		forecastTime time.Duration
		expected     int
	}{
		{"0h", 0 * time.Hour, 0},
		{"2h", 2 * time.Hour, 2},
		{"20h", 20 * time.Hour, 20},
		{"40h", 40 * time.Hour, 40},
		{"1h10m", 1*time.Hour + 10*time.Minute, 1},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			result := GetForecastHour(tt.forecastTime)
			assert.Equal(test, tt.expected, result)
		})
	}
}

func TestGetForecastMinute(test *testing.T) {
	tests := []struct {
		name         string
		forecastTime time.Duration
		expected     int
	}{
		{"0h", 0 * time.Hour, 0},
		{"1h", 1 * time.Hour, 0},
		{"20h", 20 * time.Hour, 0},
		{"1h01m", 1*time.Hour + 1*time.Minute, 1},
		{"1h10m", 1*time.Hour + 10*time.Minute, 10},
		{"48h55m", 48*time.Hour + 55*time.Minute, 55},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			result := GetForecastMinute(tt.forecastTime)
			assert.Equal(test, tt.expected, result)
		})
	}
}

func TestFormatForecastTimeShort(test *testing.T) {
	tests := []struct {
		name         string
		forecastTime time.Duration
		expected     string
	}{
		{"0h", 0 * time.Hour, "000h"},
		{"1h", 1 * time.Hour, "001h"},
		{"20h", 20 * time.Hour, "020h"},
		{"1h10m", 1*time.Hour + 10*time.Minute, "001h10m"},
		{"24h5m", 24*time.Hour + 5*time.Minute, "024h05m"},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			result := FormatForecastTimeShort(tt.forecastTime)
			assert.Equal(test, tt.expected, result)
		})
	}
}
