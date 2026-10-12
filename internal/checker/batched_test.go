package checker

import (
	"os"
	"path/filepath"
	"strings"
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
	// Submission order is nondeterministic: the two stability observations
	// run concurrently, so index by forecast time instead of position.
	byForecastTime := make(map[time.Duration]CheckResult)
	for _, result := range submitted {
		assert.NoError(t, result.Error)
		byForecastTime[result.ForecastTime] = result
	}
	assert.Equal(t, filepath.Join(runtimeDir, "data.000.grb2"), byForecastTime[0].FilePath)
	assert.Equal(t, filepath.Join(archiveDir, "data.003.grb2"), byForecastTime[3*time.Hour].FilePath)
}

func TestResolveArrivalOrder(t *testing.T) {
	// CLI/YAML overrides the data config in both directions.
	assert.Equal(t, "sequential", resolveArrivalOrder("sequential", "unordered"))
	assert.Equal(t, "unordered", resolveArrivalOrder("unordered", "sequential"))
	// auto (and the not-merged empty value) follows the data config.
	assert.Equal(t, "sequential", resolveArrivalOrder(ArrivalOrderAuto, "sequential"))
	assert.Equal(t, "unordered", resolveArrivalOrder(ArrivalOrderAuto, "unordered"))
	assert.Equal(t, "sequential", resolveArrivalOrder("", "sequential"))
	// Default: unordered.
	assert.Equal(t, "unordered", resolveArrivalOrder(ArrivalOrderAuto, ""))
	assert.Equal(t, "unordered", resolveArrivalOrder("", ""))
}

func TestGateParticipants(t *testing.T) {
	targets := []*batchedTarget{
		{forecastTime: 6 * time.Hour},
		{forecastTime: 0},
		{forecastTime: 3 * time.Hour},
	}

	// While the gate is held, nothing participates.
	assert.Empty(t, gateParticipants(targets, true))

	// Otherwise only the earliest forecast time participates.
	got := gateParticipants(targets, false)
	require.Len(t, got, 1)
	assert.Equal(t, time.Duration(0), got[0].forecastTime)

	assert.Empty(t, gateParticipants(nil, false))
}

// gatedLoopConfig builds a two-time config whose files live in one directory,
// with the given arrival order (empty omits the field).
func gatedLoopConfig(dir string, arrivalOrder string) string {
	header := ""
	if arrivalOrder != "" {
		header = "_arrival_order: " + arrivalOrder + "\n"
	}
	return header + `
default: NOTFOUND
file_name: data.{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: ` + filepath.ToSlash(dir) + `
`
}

// runGatedLoop runs checkBatched over two forecast times where the 003h file
// exists from the start but the 000h file only appears after a delay. The
// execute command appends the forecast hour to orderFile, recording the
// command execution order. configOrder sets the data config's
// _arrival_order, checkerOrder sets the checker's arrival_order override.
func runGatedLoop(t *testing.T, configOrder string, checkerOrder string) []string {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "data.003.grb2"), []byte("data"), 0644))

	orderFile := filepath.Join(dir, "order.txt")
	templates, err := buildCommandTemplates("echo {{.ForecastHour}} >> "+filepath.ToSlash(orderFile), nil)
	require.NoError(t, err)

	exec := newCommandExecutor(2, templates, time.Now(), "", nil)

	config := DefaultCheckerConfig()
	config.CheckMode = CheckModeBatched
	config.CheckIntervalJitter = 0
	config.ArrivalOrder = checkerOrder

	// The 000h file arrives late.
	go func() {
		time.Sleep(400 * time.Millisecond)
		_ = os.WriteFile(filepath.Join(dir, "data.000.grb2"), []byte("data"), 0644)
	}()

	checkBatched(
		config,
		time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC),
		gatedLoopConfig(dir, configOrder),
		nil,
		20*time.Millisecond,
		[]time.Duration{0, 3 * time.Hour},
		exec,
	)

	data, err := os.ReadFile(orderFile)
	require.NoError(t, err)
	return strings.Fields(string(data))
}

func TestCheckBatchedSequentialGatesLaterTimes(t *testing.T) {
	// auto follows the data config, which declares sequential: the later
	// 003h file (present from the start) must wait for the delayed 000h.
	assert.Equal(t, []string{"000", "003"}, runGatedLoop(t, "sequential", ArrivalOrderAuto))
	// Same with an explicit CLI/YAML override upgrading unordered→sequential.
	assert.Equal(t, []string{"000", "003"}, runGatedLoop(t, "unordered", "sequential"))
}

func TestCheckBatchedUnorderedDoesNotGate(t *testing.T) {
	// unordered: the 003h file exists immediately and completes first, even
	// though 000h is the earlier forecast time.
	assert.Equal(t, []string{"003", "000"}, runGatedLoop(t, "unordered", ArrivalOrderAuto))
	// The emergency switch: CLI/YAML override downgrading sequential→unordered.
	assert.Equal(t, []string{"003", "000"}, runGatedLoop(t, "sequential", "unordered"))
}
