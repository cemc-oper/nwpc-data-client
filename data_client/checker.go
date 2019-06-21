package data_client

import (
	"fmt"
	"strconv"
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

func CheckForecastTime(value string) (string, error) {
	if len(value) > 3 {
		return "", fmt.Errorf("length of forecast time must less or equal to 3")
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d", intValue), nil
}
