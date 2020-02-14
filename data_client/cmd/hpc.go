package cmd

import (
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(hpcCmd)

	hpcCmd.Flags().SortFlags = false

	hpcCmd.Flags().StringVar(&configDir, "config-dir", "",
		"Config dir")

	hpcCmd.Flags().StringVar(&dataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	hpcCmd.Flags().StringVar(&locationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	hpcCmd.Flags().StringVar(&startTimeSting, "start-time", "",
		"start time, YYYYMMDDHH, such as 2020021400")
	hpcCmd.Flags().StringVar(&forecastTimeString, "forecast-time", "",
		"forecast time, FFFh, such as 0h, 120h")

	hpcCmd.Flags().StringVar(&storageUser, "storage-user", user, "user name for storage.")
	hpcCmd.Flags().StringVar(&storageHost, "storage-host", "10.40.140.44:22", "host for storage")
	hpcCmd.Flags().StringVar(&privateKeyFilePath, "private-key", fmt.Sprintf("%s/.ssh/id_rsa", home),
		"private key file path")
	hpcCmd.Flags().StringVar(&hostKeyFilePath, "host-key", fmt.Sprintf("%s/.ssh/known_hosts", home),
		"host key file path")

	hpcCmd.Flags().BoolVar(&showTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")
}

const hpcCommandName = "hpc"

const hpcCommandDocString = `nwpc_data_client hpc
Find data path on hpc using config files in config dir.

Support both to find local files and to find files on storage nodes.
`

var hpcCmd = &cobra.Command{
	Use:   hpcCommandName,
	Short: "Find data path on hpc.",
	Long:  hpcCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if showTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")
		cmd.MarkFlagRequired("start-time")
		cmd.MarkFlagRequired("forecast-time")

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if showTypes {
			showDataTypes(cmd, args)
		} else {
			runHpcCommand(cmd, args)
		}
	},
}

func runHpcCommand(cmd *cobra.Command, args []string) {
	startTime, err := common.ParseStartTime(startTimeSting)
	if err != nil {
		log.Errorf("check startTime failed: %s", err)
		return
	}

	forecastTime, err = common.ParseForecastTime(forecastTimeString)
	if err != nil {
		log.Errorf("check forecastTime failed: %s", err)
		return
	}

	if len(configDir) == 0 {
		dataType = hpcCommandName + "/" + dataType
	}

	dataConfig, err2 := common.LoadConfig(configDir, dataType)

	if err2 != nil {
		log.Errorf("load config failed: %s", err2)
		return
	}

	levels := strings.Split(locationLevels, ",")
	filePath := common.FindHpcFile(dataConfig, levels, startTime, forecastTime,
		storageUser, storageHost, privateKeyFilePath, hostKeyFilePath)

	fmt.Printf("%s\n", filePath.PathType)
	fmt.Printf("%s\n", filePath.Path)
}
