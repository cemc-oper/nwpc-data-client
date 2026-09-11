package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cemc-oper/nwpc-data-client/common/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	configExportCmd.Flags().SortFlags = false

	configExportCmd.Flags().BoolVar(&exportForce, "force", false,
		"Overwrite existing files in the output directory.")
}

const configExportCommandName = "export"

const configExportCommandDocString = `Export the embedded data configs into a directory tree.

The exported tree mirrors the conf/ directory of the nwpc-data-client
repository: each embedded config "local/<system>/<version>/<format>/<name>"
is written to "<output-dir>/local/<system>/<version>/<format>/<name>.yaml",
byte-identical to the source file in conf/.

This is the transitional answer to the "--data-config-dir replaces the
embedded configs entirely" limitation: a deployment pipeline can merge an
exported baseline with an incremental config repository before distribution,
instead of copying the embedded configs by hand.

Existing files are not overwritten unless --force is given.`

const configExportCommandExample = `  # Export the embedded configs into ./exported_conf.
  nwpc_data_client config export ./exported_conf`

var configExportCmd = &cobra.Command{
	Use:     configExportCommandName + " <output-dir>",
	Short:   "Export the embedded data configs into a directory tree.",
	Long:    configExportCommandDocString,
	Example: configExportCommandExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputDir := args[0]

		writtenCount, err := exportEmbeddedConfigs(outputDir, exportForce)
		if err != nil {
			log.Fatalf("config export failed: %v", err)
		}
		fmt.Printf("config export done: %d file(s) written to %s\n", writtenCount, outputDir)
	},
}

// exportEmbeddedConfigs writes every embedded config to
// <output-dir>/<data-type>.yaml, creating parent directories as needed.
// Existing files are left untouched unless force is true.
func exportEmbeddedConfigs(outputDir string, force bool) (int, error) {
	writtenCount := 0
	for _, item := range config.EmbeddedConfigs {
		dataType := item[0]
		content := item[1]

		configFilePath := filepath.Join(outputDir, filepath.FromSlash(dataType)+".yaml")
		if !force {
			if _, err := os.Stat(configFilePath); err == nil {
				return writtenCount, fmt.Errorf("%s: file exists, use --force to overwrite", configFilePath)
			}
		}

		if err := os.MkdirAll(filepath.Dir(configFilePath), 0755); err != nil {
			return writtenCount, fmt.Errorf("%s: create directory has error: %v", configFilePath, err)
		}
		if err := os.WriteFile(configFilePath, []byte(content), 0644); err != nil {
			return writtenCount, fmt.Errorf("%s: write has error: %v", configFilePath, err)
		}
		writtenCount++
		log.Debugf("exported %s", configFilePath)
	}
	return writtenCount, nil
}
