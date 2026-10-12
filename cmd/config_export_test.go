package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cemc-oper/nwpc-data-client/common"
	"github.com/cemc-oper/nwpc-data-client/internal/embedded"
	"github.com/stretchr/testify/require"
)

// The exported tree must be byte-identical to the conf/ directory the
// embedded configs were generated from, and every exported file must pass
// config lint.
func TestExportEmbeddedConfigs_MatchesConfDir(t *testing.T) {
	tempDir := t.TempDir()

	writtenCount, err := exportEmbeddedConfigs(tempDir, false)
	require.NoError(t, err)
	require.Equal(t, len(embedded.EmbeddedConfigs), writtenCount)

	confDir := filepath.Join("..", "conf")
	confCount := 0
	err = filepath.Walk(confDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != common.ConfigFileBasename {
			return nil
		}
		confCount++

		relPath, err := filepath.Rel(confDir, path)
		require.NoError(t, err)
		exportedPath := filepath.Join(tempDir, relPath)

		expected, err := os.ReadFile(path)
		require.NoError(t, err)
		exported, err := os.ReadFile(exportedPath)
		require.NoError(t, err, "exported file missing: %s", relPath)
		require.Equal(t, string(expected), string(exported), "exported content differs: %s", relPath)

		require.NoError(t, common.LintConfigFile(exportedPath), "exported file does not pass lint: %s", relPath)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, confCount, writtenCount, "exported file count differs from conf/")
}

func TestExportEmbeddedConfigs_NoOverwriteWithoutForce(t *testing.T) {
	tempDir := t.TempDir()

	_, err := exportEmbeddedConfigs(tempDir, false)
	require.NoError(t, err)

	_, err = exportEmbeddedConfigs(tempDir, false)
	require.Error(t, err, "re-export without --force should fail")
	require.Contains(t, err.Error(), "--force")

	writtenCount, err := exportEmbeddedConfigs(tempDir, true)
	require.NoError(t, err)
	require.Equal(t, len(embedded.EmbeddedConfigs), writtenCount)
}
