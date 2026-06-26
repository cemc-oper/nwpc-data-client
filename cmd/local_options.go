package cmd

import (
	"time"
)

var (
	localDataConfigDir      = ""
	localDataConfigFile     = ""
	localDataType           = ""
	localLocationLevels     = ""
	localStartTimeString    = ""
	localForecastTimeString = ""
	localMember             = ""

	localShowTypes = false
	localDebugMode = false

	localStartTime    time.Time
	localForecastTime time.Duration
)
