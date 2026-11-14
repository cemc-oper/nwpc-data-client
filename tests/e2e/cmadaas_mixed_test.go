package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	e2eMixedDataType   = "e2e_mixed/current/grib2/orig"
	e2eMixedConfigPath = e2eMixedDataType + ".yaml"
)

func mixedLocalArgs(configDir string, storageBase string, extra ...string) []string {
	args := []string{
		"local",
		"--data-type=" + e2eMixedDataType,
		"--data-config-dir=" + configDir,
		"--start-time=2026101200",
		"--forecast-time=24h",
		"--storage-base=" + storageBase,
	}
	return append(args, extra...)
}

func writeMixedFile(t *testing.T, path string, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

func mixedCMADAASOnly(t *testing.T) (string, string, string) {
	t.Helper()
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eMixedConfigPath, rootDir)
	storageBase := filepath.Join(rootDir, "cmadaas")
	source := filepath.Join(storageBase, "2026", "20261012", "cmadaas.2026101200024.grib2")
	writeMixedFile(t, source, "cmadaas data")
	return configDir, storageBase, source
}

func runMixedFetch(t *testing.T, configDir string, storageBase string, targetDir string, extra ...string) (string, string, error) {
	t.Helper()
	args := []string{
		"fetch",
		"--data-type=" + e2eMixedDataType,
		"--data-config-dir=" + configDir,
		"--start-time=2026101200",
		"--forecast-time=24h",
		"--target-dir=" + targetDir,
		"--stable-interval=1ms",
		"--storage-base=" + storageBase,
	}
	args = append(args, extra...)
	cmd := exec.Command(getBinaryPath(t), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return strings.TrimSpace(stdout.String()), stderr.String(), err
}

func TestMixedSourceFilenamePriority(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eMixedConfigPath, rootDir)
	storageBase := filepath.Join(rootDir, "cmadaas")
	directory := filepath.Join(storageBase, "2026", "20261012")
	hpcFile := filepath.Join(directory, "hpc.2026101200024.grb2")
	cmadaasFile := filepath.Join(directory, "cmadaas.2026101200024.grib2")
	writeMixedFile(t, hpcFile, "hpc name")
	writeMixedFile(t, cmadaasFile, "cmadaas name")

	got := runCLI(t, mixedLocalArgs(configDir, storageBase)...)

	assert.Equal(t, hpcFile, got, "file_names must take precedence over file_name within each path")
}

func TestMixedSourceOriginalLocalPathPriority(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eMixedConfigPath, rootDir)
	storageBase := filepath.Join(rootDir, "cmadaas")
	hpcFile := filepath.Join(rootDir, "archive", "2026101200", "hpc.2026101200024.grb2")
	cmadaasFile := filepath.Join(storageBase, "2026", "20261012", "cmadaas.2026101200024.grib2")
	writeMixedFile(t, hpcFile, "original local source")
	writeMixedFile(t, cmadaasFile, "CMADAAS fallback")

	got := runCLI(t, mixedLocalArgs(configDir, storageBase)...)
	assert.Equal(t, hpcFile, got)

	targetDir := filepath.Join(rootDir, "target")
	stdout, stderr, err := runMixedFetch(t, configDir, storageBase, targetDir, "--method=copy")
	require.NoError(t, err, "fetch failed: %s", stderr)
	target := filepath.Join(targetDir, filepath.Base(hpcFile))
	assert.Equal(t, target, stdout)
	content, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, "original local source", string(content))
	_, err = os.Stat(filepath.Join(targetDir, filepath.Base(cmadaasFile)))
	assert.True(t, os.IsNotExist(err), "CMADAAS fallback must not be fetched after the original local source hits")
}

func TestMixedSourceCMADAASFallbackAcrossCommands(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		configDir, storageBase, source := mixedCMADAASOnly(t)
		got := runCLI(t, mixedLocalArgs(configDir, storageBase)...)
		assert.Equal(t, source, got)
	})

	t.Run("fetch without cmafs backend", func(t *testing.T) {
		configDir, storageBase, source := mixedCMADAASOnly(t)
		targetDir := t.TempDir()
		stdout, stderr, err := runMixedFetch(t, configDir, storageBase, targetDir, "--method=copy")
		require.NoError(t, err, "CMADAAS must be selected before the unconfigured cmafs source: %s", stderr)
		target := filepath.Join(targetDir, filepath.Base(source))
		assert.Equal(t, target, stdout)
		content, err := os.ReadFile(target)
		require.NoError(t, err)
		assert.Equal(t, "cmadaas data", string(content))
		assert.NotContains(t, stderr, "backend")
	})

	for _, mode := range []string{"per_time", "batched"} {
		t.Run("check local "+mode, func(t *testing.T) {
			configDir, storageBase, source := mixedCMADAASOnly(t)
			reportFile := filepath.Join(t.TempDir(), "report.json")
			checkerConfigPath := writeCheckerConfig(t, checkerConfig{
				DataConfigDir:  filepath.ToSlash(configDir),
				DataType:       e2eMixedDataType,
				StorageBase:    filepath.ToSlash(storageBase),
				CheckInterval:  "10ms",
				MaxCheckCount:  20,
				DelayTime:      "1ms",
				CheckMode:      mode,
				ForecastTimes:  []string{"24h"},
				ExecuteCommand: `echo {{.StorageBase}} > {{.FilePath}}.storage-base`,
				ReportFile:     reportFile,
			})

			output, err := runCheckLocal(t, checkerConfigPath)
			require.NoError(t, err, "check local failed: %s", output)
			report := readRunReport(t, reportFile)
			assert.Equal(t, "success", report.Status)
			require.Len(t, report.Results, 1)
			assert.Equal(t, source, report.Results[0].FilePath)
			assert.Equal(t, "success", report.Results[0].Status)
			content, err := os.ReadFile(source + ".storage-base")
			require.NoError(t, err)
			assert.Equal(t, filepath.ToSlash(storageBase), strings.TrimSpace(string(content)))
		})
	}
}

func TestMixedSourceArchiveLevelSkipsStorageAndCmafs(t *testing.T) {
	rootDir := t.TempDir()
	configDir := writeConfigWithRoot(t, e2eMixedConfigPath, rootDir)
	storageBase := filepath.Join(rootDir, "cmadaas")
	storageFile := filepath.Join(rootDir, "storage", "2026101200", "hpc.2026101200024.grb2")
	writeMixedFile(t, storageFile, "storage data")

	got := runCLI(t, mixedLocalArgs(configDir, storageBase, "--location-level=archive")...)
	assert.Equal(t, "NOTFOUND", got)

	targetDir := filepath.Join(rootDir, "target-must-not-exist")
	stdout, stderr, err := runMixedFetch(t, configDir, storageBase, targetDir, "--location-level=archive")
	require.Error(t, err)
	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "data not found")
	assert.NotContains(t, stderr, "backend not configured", "archive filtering must skip the cmafs source")
	_, statErr := os.Stat(targetDir)
	assert.True(t, os.IsNotExist(statErr), "a not-found fetch must not create the target directory")
}

func TestMixedSourceWithoutStorageBaseUsesOtherSources(t *testing.T) {
	newLocalCommand := func(configDir string) *exec.Cmd {
		return exec.Command(getBinaryPath(t),
			"local",
			"--data-type="+e2eMixedDataType,
			"--data-config-dir="+configDir,
			"--start-time=2026101200",
			"--forecast-time=24h",
		)
	}

	t.Run("original local source remains available", func(t *testing.T) {
		rootDir := t.TempDir()
		configDir := writeConfigWithRoot(t, e2eMixedConfigPath, rootDir)
		hpcFile := filepath.Join(rootDir, "archive", "2026101200", "hpc.2026101200024.grb2")
		writeMixedFile(t, hpcFile, "HPC data")
		output, err := newLocalCommand(configDir).CombinedOutput()
		require.NoError(t, err, "local failed without storage base: %s", output)
		assert.Equal(t, hpcFile, strings.TrimSpace(string(output)))
	})
}
