package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.AddCommand(checkLocalCmd)
}

const checkCommandName = "check"

const checkCommandDocString = `nwpc_data_client check
Check data for operation systems in NWPC.
`

var checkCmd = &cobra.Command{
	Use:   checkCommandName,
	Short: "Check data for operation systems in NWPC.",
	Long:  checkCommandDocString,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
