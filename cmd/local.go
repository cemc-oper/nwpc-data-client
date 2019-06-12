package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var (
	ConfigDir = ""
	DateType  = ""
	//ShowTypes = false
	startTime    time.Time
	forecastTime = ""
)

func checkStartTime(value string) (time.Time, error) {
	if len(value) != 10 {
		return time.Time{}, fmt.Errorf("length of start_time must be 10")
	}
	s, err := time.Parse("2006010215", value)
	if err != nil {
		return s, err
	}
	return s, nil
}

func checkForecastTime(value string) (string, error) {
	if len(value) > 3 {
		return "", fmt.Errorf("length of forecast time must less or equal to 3")
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d", intValue), nil
}

func init() {
	rootCmd.AddCommand(localCmd)
	localCmd.Flags().StringVar(&ConfigDir, "config-dir", "",
		"Config dir")
	localCmd.Flags().StringVar(&DateType, "date-type", "",
		"Data type used to locate config file path in config dir.")
	//localCmd.Flags().BoolVar(&ShowTypes, "show-types", false,
	//	"Show supported data types defined in config dir and exit.")
	localCmd.MarkFlagRequired("date-type")
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Find local data path.",
	Long: `Find local data path using config files in config dir.

    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: FFF, such as 000`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires two arguments")
		}
		var err error
		startTime, err = checkStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check startTime failed: %e", err)
		}

		forecastTime, err = checkForecastTime(args[1])
		if err != nil {
			return fmt.Errorf("check forecastTime failed: %e", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s, %s\n", startTime, forecastTime)
	},
}
