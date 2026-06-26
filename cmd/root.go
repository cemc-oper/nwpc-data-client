package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const rootCommandDocString = `nwpc_data_client is a data file finder for CEMC operational systems.

It searches for data files across multiple directories and polls a series of
files until they are generated, running commands once they are found.

Common subcommands:
  local         Find a local file.
  check local   Wait for several files and run a command(s) when each file is found.
  version       Print version information.

Use "nwpc_data_client [command] --help" for detailed help on a subcommand.`

var rootCmd = &cobra.Command{
	Use:   "nwpc_data_client",
	Short: "Find data files for operational systems in CEMC.",
	Long:  rootCommandDocString,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
