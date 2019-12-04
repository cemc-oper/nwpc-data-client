package cmd

import (
	"time"
)

var (
	configDir      = ""
	dataType       = ""
	locationLevels = ""
	maxCheckCount  = 240
	checkInterval  = "30s"

	startTime time.Time
)
