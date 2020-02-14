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
