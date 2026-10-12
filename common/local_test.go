package common

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	localTestStartTime    = time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC)
	localTestForecastTime = 24 * time.Hour
)

func writeLocalTestFile(t *testing.T, dir string, name string) string {
	t.Helper()
	err := os.MkdirAll(dir, 0755)
	assert.NoError(t, err)
	filePath := filepath.Join(dir, name)
	err = os.WriteFile(filePath, []byte("data"), 0644)
	assert.NoError(t, err)
	return filePath
}

func TestFindLocalFileFirstPathHit(t *testing.T) {
	firstDir := t.TempDir()
	secondDir := t.TempDir()
	expected := writeLocalTestFile(t, firstDir, "data.grb2")

	config := DataConfig{
		Default:  "NOTFOUND",
		FileName: "data.grb2",
		Paths: []PathItem{
			{PathType: "local", LevelType: "runtime", Path: firstDir},
			{PathType: "local", LevelType: "archive", Path: secondDir},
		},
	}

	result := FindLocalFile(config, nil, localTestStartTime, localTestForecastTime)

	assert.Equal(t, expected, result.Path)
	assert.Equal(t, "local", result.PathType)
}

func TestFindLocalFileFallbackToSecondPath(t *testing.T) {
	firstDir := t.TempDir()
	secondDir := t.TempDir()
	expected := writeLocalTestFile(t, secondDir, "data.grb2")

	config := DataConfig{
		Default:  "NOTFOUND",
		FileName: "data.grb2",
		Paths: []PathItem{
			{PathType: "local", LevelType: "runtime", Path: firstDir},
			{PathType: "local", LevelType: "archive", Path: secondDir},
		},
	}

	result := FindLocalFile(config, nil, localTestStartTime, localTestForecastTime)

	assert.Equal(t, expected, result.Path)
	assert.Equal(t, "local", result.PathType)
}

func TestFindLocalFileFileNamesInOrder(t *testing.T) {
	dir := t.TempDir()
	// The first name in file_names does not exist, so the second one is used.
	expected := writeLocalTestFile(t, dir, "second.grb2")

	config := DataConfig{
		Default:   "NOTFOUND",
		FileNames: []string{"first.grb2", "second.grb2"},
		Paths: []PathItem{
			{PathType: "local", LevelType: "runtime", Path: dir},
		},
	}

	result := FindLocalFile(config, nil, localTestStartTime, localTestForecastTime)

	assert.Equal(t, expected, result.Path)
}

func TestFindLocalFileFileNamesTakePrecedenceOverFileName(t *testing.T) {
	dir := t.TempDir()
	writeLocalTestFile(t, dir, "single.grb2")
	expected := writeLocalTestFile(t, dir, "listed.grb2")

	config := DataConfig{
		Default:   "NOTFOUND",
		FileName:  "single.grb2",
		FileNames: []string{"listed.grb2"},
		Paths: []PathItem{
			{PathType: "local", LevelType: "runtime", Path: dir},
		},
	}

	result := FindLocalFile(config, nil, localTestStartTime, localTestForecastTime)

	assert.Equal(t, expected, result.Path)
}

func TestFindLocalFileLocationLevelFilter(t *testing.T) {
	runtimeDir := t.TempDir()
	archiveDir := t.TempDir()
	expected := writeLocalTestFile(t, archiveDir, "data.grb2")

	config := DataConfig{
		Default:  "NOTFOUND",
		FileName: "data.grb2",
		Paths: []PathItem{
			{PathType: "local", LevelType: "runtime", Path: runtimeDir},
			{PathType: "local", LevelType: "archive", Path: archiveDir},
		},
	}

	// Searching only "runtime" skips the archive entry even though the file exists there.
	result := FindLocalFile(config, []string{"runtime"}, localTestStartTime, localTestForecastTime)
	assert.Equal(t, "NOTFOUND", result.Path)
	assert.Equal(t, "NOTFOUND", result.PathType)

	// Searching "archive" finds it.
	result = FindLocalFile(config, []string{"archive"}, localTestStartTime, localTestForecastTime)
	assert.Equal(t, expected, result.Path)

	// Empty levels search everything.
	result = FindLocalFile(config, nil, localTestStartTime, localTestForecastTime)
	assert.Equal(t, expected, result.Path)
}

func TestFindLocalFileNotFoundReturnsDefault(t *testing.T) {
	config := DataConfig{
		Default:  "NOTFOUND",
		FileName: "data.grb2",
		Paths: []PathItem{
			{PathType: "local", LevelType: "runtime", Path: t.TempDir()},
			{PathType: "local", LevelType: "archive", Path: t.TempDir()},
		},
	}

	result := FindLocalFile(config, nil, localTestStartTime, localTestForecastTime)

	assert.Equal(t, config.Default, result.Path)
	assert.Equal(t, config.Default, result.PathType)
}

func TestCheckLocalFile(t *testing.T) {
	dir := t.TempDir()
	existingFile := writeLocalTestFile(t, dir, "data.grb2")

	assert.True(t, CheckLocalFile(existingFile))
	assert.False(t, CheckLocalFile(filepath.Join(dir, "not-exist.grb2")))
	// A directory also passes os.Stat; this documents the current boundary behavior.
	assert.True(t, CheckLocalFile(dir))
}

func TestPrepareLocalDir(t *testing.T) {
	base := t.TempDir()
	filePath := filepath.Join(base, "level1", "level2", "level3", "data.grb2")

	PrepareLocalDir(filePath)

	info, err := os.Stat(filepath.Dir(filePath))
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
