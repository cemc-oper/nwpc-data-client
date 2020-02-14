package common

import (
	"fmt"
	"time"
)

type TimeTemplateVariable struct {
	Year     string
	Month    string
	Day      string
	Hour     string
	Forecast string
	Year4DV  string
	Month4DV string
	Day4DV   string
	Hour4DV  string
}

func GenerateTimeTemplateVariable(startTime time.Time, forecastTime time.Duration) TimeTemplateVariable {
	forecastHour := int(forecastTime.Hours())

	startTime4DV := startTime.Add(time.Hour * -3)
	tpVariable := TimeTemplateVariable{
		Year:     startTime.Format("2006"),
		Month:    startTime.Format("01"),
		Day:      startTime.Format("02"),
		Hour:     startTime.Format("15"),
		Forecast: fmt.Sprintf("%03d", forecastHour),
		Year4DV:  startTime4DV.Format("2006"),
		Month4DV: startTime4DV.Format("01"),
		Day4DV:   startTime4DV.Format("02"),
		Hour4DV:  startTime4DV.Format("15"),
	}
	return tpVariable
}
