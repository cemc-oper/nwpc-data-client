package common

import (
	"fmt"
	"time"
)

type TimeTemplateVariable struct {
	StartTime    time.Time
	ForecastTime time.Duration
	Year         string
	Month        string
	Day          string
	Hour         string
	Forecast     string
	Year4DV      string
	Month4DV     string
	Day4DV       string
	Hour4DV      string
	Year1HR      string
	Month1HR     string
	Day1HR       string
	Hour1HR      string
	Forecast1HR  string
}

func GenerateTimeTemplateVariable(startTime time.Time, forecastTime time.Duration) TimeTemplateVariable {
	forecastHour := int(forecastTime.Hours())

	startTime4DV := startTime.Add(time.Hour * -3)
	startTime1HR := startTime.Add(time.Hour * -1)
	tpVariable := TimeTemplateVariable{
		StartTime:    startTime,
		ForecastTime: forecastTime,
		Year:         startTime.Format("2006"),
		Month:        startTime.Format("01"),
		Day:          startTime.Format("02"),
		Hour:         startTime.Format("15"),
		Forecast:     fmt.Sprintf("%03d", forecastHour),
		Year4DV:      startTime4DV.Format("2006"),
		Month4DV:     startTime4DV.Format("01"),
		Day4DV:       startTime4DV.Format("02"),
		Hour4DV:      startTime4DV.Format("15"),
		Year1HR:      startTime1HR.Format("2006"),
		Month1HR:     startTime1HR.Format("01"),
		Day1HR:       startTime1HR.Format("02"),
		Hour1HR:      startTime1HR.Format("15"),
		Forecast1HR:  fmt.Sprintf("%03d", forecastHour+1),
	}
	return tpVariable
}
