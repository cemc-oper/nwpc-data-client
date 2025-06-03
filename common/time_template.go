package common

import (
	"fmt"
	"time"
)

type TimeTemplateVariable struct {
	StartTime      time.Time
	ForecastTime   time.Duration
	Year           string
	Month          string
	Day            string
	Hour           string
	ForecastHour   string
	ForecastMinute string
}

func GenerateTimeTemplateVariable(startTime time.Time, forecastTime time.Duration) TimeTemplateVariable {
	forecastHour := int(forecastTime.Hours())
	forecastMinute := int(forecastTime.Minutes()) % 60

	tpVariable := TimeTemplateVariable{
		StartTime:      startTime,
		ForecastTime:   forecastTime,
		Year:           startTime.Format("2006"),
		Month:          startTime.Format("01"),
		Day:            startTime.Format("02"),
		Hour:           startTime.Format("15"),
		ForecastHour:   fmt.Sprintf("%03d", forecastHour),
		ForecastMinute: fmt.Sprintf("%02d", forecastMinute),
	}
	return tpVariable
}
