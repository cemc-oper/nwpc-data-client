package data_client

import "time"

type TemplateVariable struct {
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

func GenerateTemplateVariable(startTime time.Time, forecastTime string) TemplateVariable {
	startTime4DV := startTime.Add(time.Hour * -3)
	tpVariable := TemplateVariable{
		Year:     startTime.Format("2006"),
		Month:    startTime.Format("01"),
		Day:      startTime.Format("02"),
		Hour:     startTime.Format("15"),
		Forecast: forecastTime,
		Year4DV:  startTime4DV.Format("2006"),
		Month4DV: startTime4DV.Format("01"),
		Day4DV:   startTime4DV.Format("02"),
		Hour4DV:  startTime4DV.Format("15"),
	}
	return tpVariable
}
