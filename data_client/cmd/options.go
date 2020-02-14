package cmd

import (
	"fmt"
	"os"
	"time"
)

var (
	configDir      = ""
	dataType       = ""
	locationLevels = ""
	showTypes      = false

	startTimeSting     = ""
	forecastTimeString = ""

	startTime    time.Time
	forecastTime time.Duration

	serviceAddress  = ""
	serviceAction   = ""
	outputDirectory = "."

	home               = os.Getenv("HOME")
	user               = os.Getenv("USER")
	storageHost        = ""
	storageUser        = ""
	hostKeyFilePath    = fmt.Sprintf("%s/.ssh/known_hosts", home)
	privateKeyFilePath = fmt.Sprintf("%s/.ssh/id_rsa", home)
)
