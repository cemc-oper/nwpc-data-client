package checker

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestExecutor(capacity int) *commandExecutor {
	return newCommandExecutor(capacity, nil, time.Now(), "", nil)
}

func pathItems(dirs ...string) []dataconfig.PathItem {
	var items []dataconfig.PathItem
	for _, dir := range dirs {
		items = append(items, dataconfig.PathItem{PathType: "local", LevelType: "runtime", Path: dir})
	}
	return items
}

func dirSetsWith(files map[string][]string) map[string]map[string]struct{} {
	sets := make(map[string]map[string]struct{})
	for dir, names := range files {
		set := make(map[string]struct{}, len(names))
		for _, name := range names {
			set[name] = struct{}{}
		}
		sets[dir] = set
	}
	return sets
}

func TestExpandBatchedTargetsLevelFilter(t *testing.T) {
	content := `
default: NOTFOUND
file_name: data.{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: /data/runtime
  - type: local
    level: archive
    path: /data/archive
`
	exec := newTestExecutor(2)
	defer func() { exec.finish(); exec.wait() }()

	startTime := time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC)

	targets := expandBatchedTargets(content, startTime, []time.Duration{0}, "", nil, exec)
	require.Len(t, targets, 1)
	assert.Len(t, targets[0].paths, 2)
	assert.Equal(t, "NOTFOUND", targets[0].defaultPath)
	assert.Equal(t, []string{"data.000.grb2"}, targets[0].fileNames)

	targets = expandBatchedTargets(content, startTime, []time.Duration{0}, "", []string{"runtime"}, exec)
	require.Len(t, targets, 1)
	require.Len(t, targets[0].paths, 1)
	assert.Equal(t, "runtime", targets[0].paths[0].LevelType)
}

func TestExpandBatchedTargetsFileNames(t *testing.T) {
	content := `
default: NOTFOUND
file_name: b.grb2
file_names:
  - a.grb2
paths:
  - type: local
    level: runtime
    path: /data
`
	exec := newTestExecutor(1)
	defer func() { exec.finish(); exec.wait() }()

	targets := expandBatchedTargets(content, time.Now(), []time.Duration{0}, "", nil, exec)
	require.Len(t, targets, 1)
	// file_names comes first, file_name is appended after it (FindLocalFile order).
	assert.Equal(t, []string{"a.grb2", "b.grb2"}, targets[0].fileNames)
}

func TestCollectDirectoriesDedup(t *testing.T) {
	targets := []*batchedTarget{
		{paths: pathItems("/a/runtime", "/a/archive")},
		{paths: pathItems("/a/runtime", "/b/runtime")},
	}
	assert.Equal(t,
		[]string{"/a/runtime", "/a/archive", "/b/runtime"},
		collectDirectories(targets))
}

func TestMatchTargetPriorityFirstPathWins(t *testing.T) {
	target := &batchedTarget{
		fileNames: []string{"data.grb2"},
		paths:     pathItems("/data/runtime", "/data/archive"),
	}
	// The file exists in both directories; the first path (runtime) wins.
	sets := dirSetsWith(map[string][]string{
		"/data/runtime": {"data.grb2"},
		"/data/archive": {"data.grb2"},
	})

	filePath, found := matchTarget(target, sets)

	assert.True(t, found)
	assert.Equal(t, filepath.Join("/data/runtime", "data.grb2"), filePath)
}

func TestMatchTargetHitInLaterPath(t *testing.T) {
	target := &batchedTarget{
		fileNames: []string{"data.grb2"},
		paths:     pathItems("/data/runtime", "/data/archive"),
	}
	// The file only exists in the second (archive) directory.
	sets := dirSetsWith(map[string][]string{
		"/data/runtime": {},
		"/data/archive": {"data.grb2"},
	})

	filePath, found := matchTarget(target, sets)

	assert.True(t, found)
	assert.Equal(t, filepath.Join("/data/archive", "data.grb2"), filePath)
}

func TestMatchTargetMissingDirectory(t *testing.T) {
	target := &batchedTarget{
		fileNames: []string{"data.grb2"},
		paths:     pathItems("/data/runtime", "/data/archive"),
	}
	// /data/runtime was never read this round (e.g. skipped after an error);
	// it must not match, while the read archive directory still can.
	sets := dirSetsWith(map[string][]string{
		"/data/archive": {"data.grb2"},
	})

	filePath, found := matchTarget(target, sets)

	assert.True(t, found)
	assert.Equal(t, filepath.Join("/data/archive", "data.grb2"), filePath)

	// No readable directory contains the file.
	_, found = matchTarget(target, dirSetsWith(map[string][]string{
		"/data/runtime": {"other.grb2"},
	}))
	assert.False(t, found)
}

func TestReadDirSetMissingDir(t *testing.T) {
	set, err := readDirSet(filepath.Join(t.TempDir(), "not-exist"))
	assert.NoError(t, err)
	assert.Empty(t, set)
}

func TestReadDirSetListsEntries(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.grb2"), []byte("x"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.grb2"), []byte("x"), 0644))

	set, err := readDirSet(dir)

	assert.NoError(t, err)
	assert.Equal(t, map[string]struct{}{"a.grb2": {}, "b.grb2": {}}, set)
}

func TestReadDirSetOnFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "a-file")
	require.NoError(t, os.WriteFile(filePath, []byte("x"), 0644))

	_, err := readDirSet(filePath)
	assert.Error(t, err)
}

func TestCheckBatchedStableFiles(t *testing.T) {
	runtimeDir := t.TempDir()
	archiveDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(runtimeDir, "data.000.grb2"), []byte("data0"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(archiveDir, "data.003.grb2"), []byte("data3"), 0644))

	content := `
default: NOTFOUND
file_name: data.{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: ` + filepath.ToSlash(runtimeDir) + `
  - type: local
    level: archive
    path: ` + filepath.ToSlash(archiveDir) + `
`
	var submitted []CheckResult
	exec := newCommandExecutor(2, nil, time.Now(), "", func(result CheckResult) {
		submitted = append(submitted, result)
	})

	config := DefaultCheckerConfig()
	config.CheckMode = CheckModeBatched
	config.CheckIntervalJitter = 0

	checkBatched(
		config,
		time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC),
		content,
		nil,
		10*time.Millisecond,
		[]time.Duration{0, 3 * time.Hour},
		exec,
	)

	require.Len(t, submitted, 2)
	for _, result := range submitted {
		assert.NoError(t, result.Error)
	}
	assert.Equal(t, time.Duration(0), submitted[0].ForecastTime)
	assert.Equal(t, filepath.Join(runtimeDir, "data.000.grb2"), submitted[0].FilePath)
	assert.Equal(t, 3*time.Hour, submitted[1].ForecastTime)
	assert.Equal(t, filepath.Join(archiveDir, "data.003.grb2"), submitted[1].FilePath)
}
