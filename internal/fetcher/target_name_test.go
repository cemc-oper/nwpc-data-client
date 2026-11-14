package fetcher

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func targetNameFetchConfig(targetDir string, method Method, targetName string) Config {
	config := testFetchConfig(targetDir, method)
	config.TargetName = targetName
	config.TemplateVariable = dataconfig.GenerateConfigTemplateVariableWithStorageBase(
		time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC), 24*time.Hour, "014", "")
	return config
}

func TestFetchTargetNameFixed(t *testing.T) {
	source, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()

	result, err := Fetch(context.Background(), targetNameFetchConfig(targetDir, MethodCopy, "input.grb2"), dataConfig, nil)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(targetDir, "input.grb2"), result.FinalPath)

	sourceContent, err := os.ReadFile(source)
	require.NoError(t, err)
	targetContent, err := os.ReadFile(result.FinalPath)
	require.NoError(t, err)
	assert.Equal(t, sourceContent, targetContent)
	assert.NoError(t, assertNoPartialFiles(targetDir))
}

func TestFetchTargetNameTemplate(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		targetName   string
		expectedName string
	}{
		{
			name:         "start time and forecast hour",
			targetName:   "gfs_{{.Year}}{{.Month}}{{.Day}}{{.Hour}}_{{.ForecastHour}}.grb2",
			expectedName: "gfs_2026062400_024.grb2",
		},
		{
			name:         "member",
			targetName:   "ens_{{.Member}}.grb2",
			expectedName: "ens_014.grb2",
		},
		{
			name:         "source name",
			targetName:   "pre_{{.SourceName}}",
			expectedName: "pre_data.grb2",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, dataConfig := writeFetchSource(t, "runtime")
			targetDir := t.TempDir()
			config := targetNameFetchConfig(targetDir, MethodCopy, testCase.targetName)

			result, err := Fetch(context.Background(), config, dataConfig, nil)
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(targetDir, testCase.expectedName), result.FinalPath)

			reused, err := Fetch(context.Background(), config, dataConfig, nil)
			require.NoError(t, err)
			assert.True(t, reused.Reused)
			assert.Equal(t, result.FinalPath, reused.FinalPath)
		})
	}
}

func TestFetchTargetNameSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks may require additional Windows privileges")
	}
	source, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()

	result, err := Fetch(context.Background(), targetNameFetchConfig(targetDir, MethodSymlink, "linked.grb2"), dataConfig, nil)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(targetDir, "linked.grb2"), result.FinalPath)
	linkTarget, err := os.Readlink(result.FinalPath)
	require.NoError(t, err)
	assert.Equal(t, source, linkTarget)
}

func TestFetchTargetNameInvalid(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		targetName string
	}{
		{name: "parent directory escape", targetName: "../escape.grb2"},
		{name: "nested path", targetName: "sub/data.grb2"},
		{name: "absolute path", targetName: "/tmp/escape.grb2"},
		{name: "backslash separator", targetName: `sub\data.grb2`},
		{name: "empty rendered name", targetName: "{{.Member}}"},
		{name: "dot", targetName: "."},
		{name: "dot dot", targetName: ".."},
		{name: "newline", targetName: "data\n.grb2"},
		{name: "invalid template syntax", targetName: "{{.Year"},
		{name: "unknown template field", targetName: "{{.NoSuchField}}"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, dataConfig := writeFetchSource(t, "runtime")
			targetDir := filepath.Join(t.TempDir(), "not-created")
			config := targetNameFetchConfig(targetDir, MethodCopy, testCase.targetName)
			if testCase.name == "empty rendered name" {
				config.TemplateVariable.Member = ""
			}

			_, err := Fetch(context.Background(), config, dataConfig, nil)
			require.Error(t, err)
			_, statErr := os.Stat(targetDir)
			assert.True(t, os.IsNotExist(statErr), "target directory must not be created for an invalid target name")
		})
	}
}

func TestFetchTargetNameConflict(t *testing.T) {
	_, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()
	conflictingPath := filepath.Join(targetDir, "input.grb2")
	require.NoError(t, os.WriteFile(conflictingPath, []byte("other content"), 0644))

	_, err := Fetch(context.Background(), targetNameFetchConfig(targetDir, MethodCopy, "input.grb2"), dataConfig, nil)
	assert.ErrorIs(t, err, ErrTargetConflict)
	content, readErr := os.ReadFile(conflictingPath)
	require.NoError(t, readErr)
	assert.Equal(t, "other content", string(content))
	assert.NoError(t, assertNoPartialFiles(targetDir))
}
