package cmd

import (
	"time"
)

var (
	localConfigDir          = ""
	localConfigFile         = ""
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
