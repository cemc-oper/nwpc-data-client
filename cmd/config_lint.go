package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	configLintCmd.Flags().SortFlags = false

	configLintCmd.Flags().StringVar(&lintDataConfigFile, "data-config-file", "",
		"Path to a single data config file. If set, the config directory argument is ignored.")
}

const configLintCommandName = "lint"

const configLintCommandDocString = `Lint data config files of any config format version (v1/v2/v3).

Each file is validated in two stages:

  1. template validation with the real text/template parser, using the
     delimiters, functions and variables matching the detected config format
     version (the "_config_version" field wins; heuristic sniffing is the
     fallback). Unknown functions fail at parse, unknown variables at execute.
  2. structural validation of the rendered result against the data config
     structure (paths items require type/level/path, and so on).

Pass a config directory to lint every *.yaml file under it recursively, or
use --data-config-file to lint a single file.`

const configLintCommandExample = `  # Lint all configs under conf/local recursively.
  nwpc_data_client config lint conf/local

  # Lint a single config file.
  nwpc_data_client config lint --data-config-file=./my_config.yaml`

var configLintCmd = &cobra.Command{
	Use:     configLintCommandName + " [config-dir]",
	Short:   "Validate data config files.",
	Long:    configLintCommandDocString,
	Example: configLintCommandExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if lintDataConfigFile != "" {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("requires a config directory argument or --data-config-file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var configFilePaths []string
		if lintDataConfigFile != "" {
			configFilePaths = []string{lintDataConfigFile}
		} else {
			configFilePaths = collectConfigFiles(args[0])
		}

		if len(configFilePaths) == 0 {
			log.Warn("no config file found")
			return
		}

		failedCount := 0
		for _, configFilePath := range configFilePaths {
			if err := dataconfig.LintConfigFile(configFilePath); err != nil {
				failedCount++
				log.Errorf("%v", err)
			}
		}

		if failedCount > 0 {
			log.Fatalf("config lint failed: %d of %d file(s) have errors", failedCount, len(configFilePaths))
		}
		fmt.Printf("config lint passed: %d file(s)\n", len(configFilePaths))
	},
}

// collectConfigFiles walks a config directory and returns all *.yaml files.
func collectConfigFiles(configDir string) []string {
	var configFilePaths []string
	err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			configFilePaths = append(configFilePaths, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("walk config directory %s has error: %s", configDir, err)
	}
	return configFilePaths
}
