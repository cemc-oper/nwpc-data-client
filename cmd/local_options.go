package cmd

import (
	"time"
)

var (
	localDataConfigDir      = ""
	localDataConfigFile     = ""
	localDataType           = ""
	localLocationLevels     = ""
	localStorageBase        = ""
	localStartTimeString    = ""
	localForecastTimeString = ""
	localMember             = ""

	localShowTypes = false
	localDebugMode = false

	localStartTime    time.Time
	localForecastTime time.Duration
)
