package cmd

import (
	"time"
)

var (
	configDir      = ""
	dataType       = ""
	locationLevels = ""
	maxCheckCount  = 2880
	checkInterval  = "5s"
	executeCommand = ""

	startTime time.Time
)
