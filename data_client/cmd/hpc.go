package cmd

import (
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
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

	hpcCmd.Flags().StringVar(&StorageUser, "storage-user", user, "user name for storage.")
	hpcCmd.Flags().StringVar(&StorageHost, "storage-host", "10.40.140.44:22", "host for storage")
	hpcCmd.Flags().StringVar(&PrivateKeyFilePath, "private-key", fmt.Sprintf("%s/.ssh/id_rsa", home),
		"private key file path")
	hpcCmd.Flags().StringVar(&HostKeyFilePath, "host-key", fmt.Sprintf("%s/.ssh/known_hosts", home),
		"host key file path")

	hpcCmd.Flags().BoolVar(&ShowTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")
}

const hpcCommandDocString = `nwpc_data_client hpc
Find data path on hpc using config files in config dir.

Support both to find local files and to find files on storage nodes.

Args:
    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: FFF, such as 000`

var hpcCmd = &cobra.Command{
	Use:   "hpc",
	Short: "Find data path on hpc.",
	Long:  hpcCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if ShowTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")

		if len(args) != 2 {
			return errors.New("requires two arguments")
		}
		var err error
		StartTime, err = common.CheckStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check StartTime failed: %s", err)
		}

		ForecastTime, err = common.CheckForecastHour(args[1])
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
	configFilePath, err := common.FindConfig(ConfigDir, DataType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "model data type config is not found.\n")
		return
	}
	hpcDataConfig, err2 := common.LoadConfig(configFilePath)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %s\n", err2)
		return
	}
	filePath := findHpcFile(hpcDataConfig, StartTime, ForecastTime)
	fmt.Printf("%s\n", filePath.PathType)
	fmt.Printf("%s\n", filePath.Path)
}

func findHpcFile(config common.DataConfig, startTime time.Time, forecastTime time.Duration) common.PathItem {
	tpVar := common.GenerateTemplateVariable(startTime, forecastTime)

	fileNameTemplate := template.Must(template.New("fileName").
		Delims("{", "}").Parse(config.FileName))

	var fileNameBuilder strings.Builder
	err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file name template execute has error: %s\n", err)
		return common.PathItem{
			Path:     config.Default,
			PathType: config.Default,
		}
	}
	fileName := fileNameBuilder.String()

	for _, pathItem := range config.Paths {
		path := pathItem.Path
		pathType := pathItem.PathType
		dirPathTemplate := template.Must(template.New("dirPath").Delims("{", "}").Parse(path))

		var dirPathBuilder strings.Builder
		err = dirPathTemplate.Execute(&dirPathBuilder, tpVar)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dir path template execute has error: %s\n", err)
			continue
		}
		dirPath := dirPathBuilder.String()
		filePath := filepath.Join(dirPath, fileName)
		//fmt.Printf("%s\n", filePath)

		if pathType == "storage" {
			if common.CheckFileOverSSH(filePath, StorageUser, StorageHost, PrivateKeyFilePath, HostKeyFilePath) {
				return common.PathItem{
					Path:     filePath,
					PathType: pathType,
				}
			}
		} else if pathType == "local" {
			// check if file exists
			if common.CheckLocalFile(filePath) {
				return common.PathItem{
					Path:     filePath,
					PathType: pathType,
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "path type is not supported: %s", pathType)
		}
	}

	return common.PathItem{
		Path:     config.Default,
		PathType: config.Default,
	}
}
