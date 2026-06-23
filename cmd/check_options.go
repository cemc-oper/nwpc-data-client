package cmd

import (
	"time"
)

var (
	// checkFlagConfig holds values bound directly from CLI flags.
	checkFlagConfig CheckerConfig
	// checkConfigFile is the path to the optional YAML runtime config file.
	checkConfigFile = ""

	checkStartTime time.Time
)
