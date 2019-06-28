package cmd

import (
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/data_client"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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

	hpcCmd.Flags().BoolVar(&ShowTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")

	hpcCmd.Flags().StringVar(&StorageUser, "storage-user", user, "user name for storage.")
	hpcCmd.Flags().StringVar(&StorageHost, "storage-host", "10.40.140.44:22", "host for storage")
	hpcCmd.Flags().StringVar(&PrivateKeyFilePath, "private-key", fmt.Sprintf("%s/.ssh/id_rsa", home),
		"private key file path")
	hpcCmd.Flags().StringVar(&HostKeyFilePath, "host-key", fmt.Sprintf("%s/.ssh/known_hosts", home),
		"host key file path")
}

var hpcCmd = &cobra.Command{
	Use:   "hpc",
	Short: "Find data path on hpc.",
	Long: `Find data path on hpc using config files in config dir.

Support both to find local files and to find files on storage nodes.

Args:
    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: FFF, such as 000`,
	Args: func(cmd *cobra.Command, args []string) error {
		if ShowTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")

		if len(args) != 2 {
			return errors.New("requires two arguments")
		}
		var err error
		StartTime, err = data_client.CheckStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check StartTime failed: %s", err)
		}

		ForecastTime, err = data_client.CheckForecastTime(args[1])
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
	configFilePath, err := findConfig(ConfigDir, DataType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "model data type config is not found.\n")
		return
	}
	hpcDataConfig, err2 := loadHpcConfig(configFilePath)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %s\n", err2)
		return
	}
	filePath := findHpcFile(hpcDataConfig, StartTime, ForecastTime)
	fmt.Printf("%s\n", filePath.PathType)
	fmt.Printf("%s\n", filePath.Path)
}

type HpcPathItem struct {
	PathType string `yaml:"type"`
	Path     string `yaml:"path"`
}

type HpcDataConfig struct {
	Default  string        `yaml:"default"`
	FileName string        `yaml:"file_name"`
	Paths    []HpcPathItem `yaml:"paths"`
}

func loadHpcConfig(configFilePath string) (HpcDataConfig, error) {
	dataConfig := HpcDataConfig{}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return dataConfig, err
	}

	err = yaml.Unmarshal(data, &dataConfig)
	if err != nil {
		return dataConfig, err
	}

	return dataConfig, nil
}

func findHpcFile(config HpcDataConfig, startTime time.Time, forecastTime string) HpcPathItem {
	tpVar := data_client.GenerateTemplateVariable(startTime, forecastTime)

	fileNameTemplate := template.Must(template.New("fileName").
		Delims("{", "}").Parse(config.FileName))

	var fileNameBuilder strings.Builder
	err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file name template execute has error: %s\n", err)
		return HpcPathItem{
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
			if checkFileOverSSH(filePath) {
				return HpcPathItem{
					Path:     filePath,
					PathType: pathType,
				}
			}
		} else if pathType == "local" {
			// check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				continue
			}
			return HpcPathItem{
				Path:     filePath,
				PathType: pathType,
			}
		} else {
			fmt.Fprintf(os.Stderr, "path type is not supported: %s", pathType)
		}
	}

	return HpcPathItem{
		Path:     config.Default,
		PathType: config.Default,
	}
}

func checkFileOverSSH(filePath string) bool {
	key, err := ioutil.ReadFile(PrivateKeyFilePath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	hostKeyCallback, err := knownhosts.New(HostKeyFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config := &ssh.ClientConfig{
		User: StorageUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", StorageHost, config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}
	defer session.Close()

	b, err := session.Output(
		fmt.Sprintf("[[ -f %s ]] && echo \"File exists\" || echo \"File does not exist\"", filePath))
	if err != nil {
		log.Fatalf("failed to execute: %s", err)
	}
	outputString := strings.TrimSuffix(string(b), "\n")

	if outputString == "File exists" {
		return true
	}
	return false
}
