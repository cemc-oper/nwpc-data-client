package cmd

import (
	"github.com/cemc-oper/nwpc-data-client/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print version information.",
	Long:    "Print version, build time, and git commit baked into the binary.",
	Example: "  nwpc_data_client version",
	Run: func(cmd *cobra.Command, args []string) {
		common.PrintVersionInformation()
	},
}
