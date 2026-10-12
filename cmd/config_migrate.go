package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	configMigrateCmd.Flags().SortFlags = false

	configMigrateCmd.Flags().StringVar(&migrateDataConfigFile, "data-config-file", "",
		"Path to a single data config file. If set, the config directory argument is ignored.")

	configMigrateCmd.Flags().IntVar(&migrateToVersion, "to-version", dataconfig.CurrentConfigVersion,
		"Target config format version (2 or 3).")

	configMigrateCmd.Flags().StringVar(&migrateOutputDir, "output-dir", "",
		"Write migrated files under this directory instead of overwriting the source files. "+
			"The directory structure of the config directory is mirrored.")
}

const configMigrateCommandName = "migrate"

const configMigrateCommandDocString = `Migrate data config files to a newer config format version.

Supported migrations:

  v1 -> v2: rename template variables ({.Forecast} -> {.ForecastHour},
            the 4DV/1HR offset variables -> generateStartTime /
            generateForecastTime pipelines)
  v2 -> v3: replace "{...}" action delimiters with "{{...}}" (only actions
            starting with a variable, a known function or a template keyword
            are rewritten; literal braces are kept)
  v1 -> v3: both steps in sequence

The migrated content gets a "_config_version" field and is validated with
the config lint machinery before being written. A file whose source does not
pass lint (e.g. a v1 config using functions in paths[].path, where no
functions were registered) or whose migrated result does not pass lint is
reported and left untouched.

Files are overwritten in place unless --output-dir is given. Files already
at the target version are skipped.`

const configMigrateCommandExample = `  # Migrate a single config file to v3, in place.
  nwpc_data_client config migrate --data-config-file=./my_config.yaml

  # Migrate all configs under a directory to v3, writing to another tree.
  nwpc_data_client config migrate ./old_conf/local --output-dir=./new_conf/local

  # Migrate to v2 only.
  nwpc_data_client config migrate --data-config-file=./my_config.yaml --to-version=2`

var configMigrateCmd = &cobra.Command{
	Use:     configMigrateCommandName + " [config-dir]",
	Short:   "Migrate data config files to a newer format version.",
	Long:    configMigrateCommandDocString,
	Example: configMigrateCommandExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if migrateDataConfigFile != "" {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("requires a config directory argument or --data-config-file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if migrateToVersion != dataconfig.ConfigVersion2 && migrateToVersion != dataconfig.ConfigVersion3 {
			log.Fatalf("invalid --to-version %d: migration targets are v2 and v3", migrateToVersion)
		}

		var configFilePaths []string
		var configDir string
		if migrateDataConfigFile != "" {
			configFilePaths = []string{migrateDataConfigFile}
		} else {
			configDir = args[0]
			configFilePaths = collectConfigFiles(configDir)
		}

		if len(configFilePaths) == 0 {
			log.Warn("no config file found")
			return
		}

		failedCount := 0
		migratedCount := 0
		skippedCount := 0
		for _, configFilePath := range configFilePaths {
			migrated, skip, err := migrateConfigFile(configFilePath, migrateToVersion)
			if err != nil {
				failedCount++
				log.Errorf("%v", err)
				continue
			}
			if skip {
				skippedCount++
				continue
			}

			targetPath := configFilePath
			if migrateOutputDir != "" {
				if configDir != "" {
					relPath, err := filepath.Rel(configDir, configFilePath)
					if err != nil {
						failedCount++
						log.Errorf("%s: compute relative path has error: %v", configFilePath, err)
						continue
					}
					targetPath = filepath.Join(migrateOutputDir, relPath)
				} else {
					targetPath = filepath.Join(migrateOutputDir, filepath.Base(configFilePath))
				}
				if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
					failedCount++
					log.Errorf("%s: create output directory has error: %v", targetPath, err)
					continue
				}
			}

			if err := os.WriteFile(targetPath, []byte(migrated), 0644); err != nil {
				failedCount++
				log.Errorf("%s: write has error: %v", targetPath, err)
				continue
			}
			migratedCount++
			log.Infof("migrated %s -> %s", configFilePath, targetPath)
		}

		if failedCount > 0 {
			log.Fatalf("config migrate failed: %d of %d file(s) have errors",
				failedCount, len(configFilePaths))
		}
		fmt.Printf("config migrate done: %d migrated, %d skipped (already at v%d)\n",
			migratedCount, skippedCount, migrateToVersion)
	},
}

// migrateConfigFile migrates one config file and returns the migrated
// content. skip is true when the file is already at the target version.
func migrateConfigFile(configFilePath string, toVersion int) (migrated string, skip bool, err error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", false, fmt.Errorf("%s: read: %v", configFilePath, err)
	}
	content := string(data)

	version, err := dataconfig.DetectConfigVersion(content)
	if err != nil {
		return "", false, fmt.Errorf("%s: version detect: %v", configFilePath, err)
	}
	if version == toVersion {
		log.Infof("%s: already at config format v%d, skipped", configFilePath, toVersion)
		return "", true, nil
	}

	migrated, err = dataconfig.MigrateConfigContent(content, toVersion)
	if err != nil {
		return "", false, fmt.Errorf("%s: %v", configFilePath, err)
	}
	return migrated, false, nil
}
