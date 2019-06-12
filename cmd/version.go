package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version = "Unknown Version"
	Date    = "Unknown Date"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %s\n", Version)
		fmt.Printf("Build at %s\n", Date)
	},
}
