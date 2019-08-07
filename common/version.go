package common

import "fmt"

var (
	Version   = "Unknown version"
	BuildTime = "Unknown build time"
	GitCommit = "Unknown GitCommit"
)

func PrintVersionInformation() {
	fmt.Printf("Version %s (%s)\n", Version, GitCommit)
	fmt.Printf("Build at %s\n", BuildTime)
	fmt.Printf("Please visit https://github.com/nwpc-oper/nwpc-data-client for more information.\n")
}
