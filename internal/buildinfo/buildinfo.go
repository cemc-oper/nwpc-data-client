// Package buildinfo holds the version, build time, and git commit injected
// at link time via -ldflags, and prints them for the version command.
package buildinfo

import "fmt"

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func PrintVersionInformation() {
	fmt.Printf("Version %s (%s)\n", Version, GitCommit)
	fmt.Printf("Build at %s\n", BuildTime)
	fmt.Printf("Please visit https://github.com/cemc-oper/nwpc-data-client for more information.\n")
}
