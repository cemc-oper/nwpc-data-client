package cmd

import (
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %s (%s)\n", common.Version, common.GitCommit)
		fmt.Printf("Build at %s\n", common.BuildTime)
		fmt.Printf("Please visit https://github.com/nwpc-oper/nwpc-data-client for more information.\n")
	},
}
