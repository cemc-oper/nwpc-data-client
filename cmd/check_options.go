package cmd

import (
	"github.com/cemc-oper/nwpc-data-client/internal/checker"
)

var (
	// checkFlagConfig holds values bound directly from CLI flags.
	checkFlagConfig checker.CheckerConfig
	// checkConfigFile is the path to the optional YAML runtime config file.
	checkConfigFile = ""
)
