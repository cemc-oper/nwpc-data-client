package cmd

import (
	"time"
)

var (
	configDir      = ""
	configFile     = ""
	dataType       = ""
	locationLevels = ""

	startTimeSting     = ""
	forecastTimeString = ""
	member             = ""

	showTypes = false
	debugMode = false

	startTime    time.Time
	forecastTime time.Duration
)
