package common

import (
	"fmt"
	"time"
)

func ParseStartTime(value string) (time.Time, error) {
	if len(value) != 10 {
		return time.Time{}, fmt.Errorf("length of start_time must be 10")
	}
	s, err := time.Parse("2006010215", value)
	if err != nil {
		return s, err
	}
	return s, nil
}

func ParseForecastTime(value string) (time.Duration, error) {
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse duration error: %v", err)
	}

	return d, nil
}

// GenerateStartTime generates a new time add or minus some hours.
// used in template files.
// should be used with getYear, getMonth, getDay, getHour functions.
// usage:
//
//	{generateStartTime .StartTime -3 | getYear}
func GenerateStartTime(startTime time.Time, hour int) time.Time {
	newStartTime := startTime.Add(time.Hour * time.Duration(hour))
	return newStartTime
}

func GetYear(startTime time.Time) string {
	return startTime.Format("2006")
}

func GetMonth(startTime time.Time) string {
	return startTime.Format("01")
}

func GetDay(startTime time.Time) string {
	return startTime.Format("02")
}

func GetHour(startTime time.Time) string {
	return startTime.Format("15")
}

// GenerateForecastTime generate a new time duration calculated from a forecast time and a time interval.
// should be used with getForecastHour, getForecastMinute functions.
func GenerateForecastTime(forecastTime time.Duration, timeInterval string) time.Duration {
	t, _ := time.ParseDuration(timeInterval)
	newForecastTime := forecastTime + t
	return newForecastTime
}

// GetForecastHour get hour from forecast time. used in template files.
// Usage:
//
//	{.ForecastTime | getForecastHour | printf "%03d"}
func GetForecastHour(forecastTime time.Duration) int {
	return int(forecastTime.Hours())
}

// GetForecastMinute get minute from forecast time. used in template files.
// Usage:
//
//	{.ForecastTime | getForecastMinute | printf "%02d"}
func GetForecastMinute(forecastTime time.Duration) int {
	return int(forecastTime.Minutes()) % 60
}

func FormatForecastTimeShort(forecastTime time.Duration) string {
	forecastHour := GetForecastHour(forecastTime)
	forecastMinute := GetForecastMinute(forecastTime)
	if forecastMinute == 0 {
		return fmt.Sprintf("%03dh", forecastHour)
	} else {
		return fmt.Sprintf("%03dh%02dm", forecastHour, forecastMinute)
	}
}
