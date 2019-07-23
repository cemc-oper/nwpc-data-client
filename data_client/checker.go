package data_client

import (
	"fmt"
	"time"
)

func CheckStartTime(value string) (time.Time, error) {
	if len(value) != 10 {
		return time.Time{}, fmt.Errorf("length of start_time must be 10")
	}
	s, err := time.Parse("2006010215", value)
	if err != nil {
		return s, err
	}
	return s, nil
}

func CheckForecastHour(value string) (time.Duration, error) {
	return CheckForecastTime(fmt.Sprintf("%sh", value))
}

func CheckForecastTime(value string) (time.Duration, error) {
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse duration error: %v", err)
	}

	return d, nil
}
