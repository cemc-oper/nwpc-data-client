package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.AddCommand(checkLocalCmd)
}

const checkCommandName = "check"

const checkCommandDocString = `Check data files for CEMC operational systems.

Use "nwpc_data_client check local --help" for detailed help.`

var checkCmd = &cobra.Command{
	Use:   checkCommandName,
	Short: "Wait for operational data and act when it is ready.",
	Long:  checkCommandDocString,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
