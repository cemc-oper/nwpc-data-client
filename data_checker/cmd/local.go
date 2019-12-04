package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.Flags().SortFlags = false

	localCmd.Flags().StringVar(&configDir, "data-config-dir", "",
		"Config dir.")

	localCmd.Flags().StringVar(&dataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	localCmd.Flags().StringVar(&locationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")
}

const localCommandName = "local"

const localCommandDocString = `nwpc_data_checker local
Check local data path using config files in config dir.

Args:
    start_time: YYYYMMDDHH, such as 2018080100`

var localCmd = &cobra.Command{
	Use:   localCommandName,
	Short: "Check local data.",
	Long:  localCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires one arguments")
		}
		var err error
		startTime, err = common.CheckStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check startTime failed: %s", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanWords)
		startTimeString := startTime.Format("2006010215")
		var forecastTimeList []time.Duration
		for scanner.Scan() {
			forecastTimeString := scanner.Text()
			forecastTime, err := time.ParseDuration(forecastTimeString)
			if err != nil {
				log.Fatalf("parse input has error:%v", err)
			}
			forecastTimeList = append(forecastTimeList, forecastTime)
			log.Infof("checking for data of %s + %03d...", startTimeString, int(forecastTime.Hours()))
		}
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		levels := strings.Split(locationLevels, ",")
		if len(configDir) == 0 {
			dataType = localCommandName + "/" + dataType
		}
		config, err2 := common.LoadConfig(configDir, dataType)
		for _, oneTime := range forecastTimeList {
			if err2 != nil {
				log.Fatalf("load config failed: %v\n", err2)
				return
			}
			findLocalFile(config, levels, oneTime)
		}
	},
}

func findLocalFile(config common.DataConfig, levels []string, forecastTime time.Duration) {
	pathItem := common.FindLocalFile(config, levels, startTime, forecastTime)
	fmt.Printf("%s\n", pathItem.Path)
}
