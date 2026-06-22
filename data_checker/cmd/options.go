package cmd

import (
	"time"
)

var (
	// flagConfig holds values bound directly from CLI flags.
	flagConfig CheckerConfig
	// checkerConfigFile is the path to the optional YAML runtime config file.
	checkerConfigFile = ""

	startTime time.Time
)
