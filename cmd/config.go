package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configLintCmd)
}

const configCommandName = "config"

const configCommandDocString = `Manage data config files.

Use "nwpc_data_client config lint --help" for detailed help.`

var configCmd = &cobra.Command{
	Use:   configCommandName,
	Short: "Validate and maintain data config files.",
	Long:  configCommandDocString,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
