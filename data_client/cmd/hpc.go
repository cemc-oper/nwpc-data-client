package cmd

import (
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	home               = os.Getenv("HOME")
	user               = os.Getenv("USER")
	StorageHost        = ""
	StorageUser        = ""
	HostKeyFilePath    = fmt.Sprintf("%s/.ssh/known_hosts", home)
	PrivateKeyFilePath = fmt.Sprintf("%s/.ssh/id_rsa", home)
)

func init() {
	rootCmd.AddCommand(hpcCmd)

	hpcCmd.Flags().SortFlags = false

	hpcCmd.Flags().StringVar(&ConfigDir, "config-dir", "",
		"Config dir")

	hpcCmd.Flags().StringVar(&DataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	hpcCmd.Flags().StringVar(&LocationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	hpcCmd.Flags().StringVar(&StorageUser, "storage-user", user, "user name for storage.")
	hpcCmd.Flags().StringVar(&StorageHost, "storage-host", "10.40.140.44:22", "host for storage")
	hpcCmd.Flags().StringVar(&PrivateKeyFilePath, "private-key", fmt.Sprintf("%s/.ssh/id_rsa", home),
		"private key file path")
	hpcCmd.Flags().StringVar(&HostKeyFilePath, "host-key", fmt.Sprintf("%s/.ssh/known_hosts", home),
		"host key file path")

	hpcCmd.Flags().BoolVar(&ShowTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")
}

const hpcCommandName = "hpc"

const hpcCommandDocString = `nwpc_data_client hpc
Find data path on hpc using config files in config dir.

Support both to find local files and to find files on storage nodes.

Args:
    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: FFFh, such as 120h`

var hpcCmd = &cobra.Command{
	Use:   hpcCommandName,
	Short: "Find data path on hpc.",
	Long:  hpcCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if ShowTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")

		if len(args) != 2 {
			return errors.New("requires two arguments: startTime and forecastTime")
		}
		var err error
		StartTime, err = common.CheckStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check StartTime failed: %s", err)
		}

		ForecastTime, err = common.CheckForecastTime(args[1])
		if err != nil {
			return fmt.Errorf("check ForecastTime failed: %s", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ShowTypes {
			showDataTypes(cmd, args)
		} else {
			runHpcCommand(cmd, args)
		}
	},
}

func runHpcCommand(cmd *cobra.Command, args []string) {
	if len(ConfigDir) == 0 {
		DataType = hpcCommandName + "/" + DataType
	}

	dataConfig, err2 := common.LoadConfig(ConfigDir, DataType)

	if err2 != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %s\n", err2)
		return
	}

	levels := strings.Split(LocationLevels, ",")
	filePath := common.FindHpcFile(dataConfig, levels, StartTime, ForecastTime,
		StorageUser, StorageHost, PrivateKeyFilePath, HostKeyFilePath)

	fmt.Printf("%s\n", filePath.PathType)
	fmt.Printf("%s\n", filePath.Path)
}
