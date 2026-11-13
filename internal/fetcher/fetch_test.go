package fetcher

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/filecheck"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFetchSource(t *testing.T, level string) (string, dataconfig.DataConfig) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "data.grb2")
	require.NoError(t, os.WriteFile(path, bytes.Repeat([]byte("data"), 256*1024), 0644))
	mtime := time.Unix(1700000000, 123456789)
	require.NoError(t, os.Chtimes(path, mtime, mtime))
	config := dataconfig.DataConfig{
		Default:  "NOTFOUND",
		FileName: "data.grb2",
		Paths: []dataconfig.PathItem{{
			PathType: "local", LevelType: level, Path: dir,
		}},
	}
	return path, config
}

func testFetchConfig(targetDir string, method Method) Config {
	return Config{TargetDir: targetDir, Method: method, StableInterval: time.Millisecond}
}

func TestFetchCopyAndReuse(t *testing.T) {
	source, dataConfig := writeFetchSource(t, "runtime")
	targetDir := filepath.Join(t.TempDir(), "nested", "input")

	result, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodCopy), dataConfig, nil)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(targetDir, "data.grb2"), result.FinalPath)
	assert.Equal(t, MethodCopy, result.Method)
	assert.False(t, result.Reused)

	sourceContent, err := os.ReadFile(source)
	require.NoError(t, err)
	targetContent, err := os.ReadFile(result.FinalPath)
	require.NoError(t, err)
	assert.Equal(t, sourceContent, targetContent)
	sourceInfo, err := os.Stat(source)
	require.NoError(t, err)
	targetInfo, err := os.Stat(result.FinalPath)
	require.NoError(t, err)
	assert.Equal(t, sourceInfo.Mode().Perm(), targetInfo.Mode().Perm())
	assert.True(t, sourceInfo.ModTime().Equal(targetInfo.ModTime()))

	reused, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodCopy), dataConfig, nil)
	require.NoError(t, err)
	assert.True(t, reused.Reused)
	assert.NoError(t, assertNoPartialFiles(targetDir))
}

func TestFetchSymlinkAndReuse(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks may require additional Windows privileges")
	}
	source, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()

	result, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodSymlink), dataConfig, nil)
	require.NoError(t, err)
	assert.Equal(t, MethodSymlink, result.Method)
	linkTarget, err := os.Readlink(result.FinalPath)
	require.NoError(t, err)
	assert.Equal(t, source, linkTarget)

	reused, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodSymlink), dataConfig, nil)
	require.NoError(t, err)
	assert.True(t, reused.Reused)
}

func TestFetchAutoMethod(t *testing.T) {
	for _, level := range []string{"runtime", "archive", "storage", "runtime/storage"} {
		t.Run(level, func(t *testing.T) {
			_, dataConfig := writeFetchSource(t, level)
			targetDir := t.TempDir()
			expectSymlink := !hasStorageLevel(level)
			if runtime.GOOS == "windows" && expectSymlink {
				t.Skip("creating symlinks may require additional Windows privileges")
			}
			result, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodAuto), dataConfig, nil)
			require.NoError(t, err)

			if expectSymlink {
				assert.Equal(t, MethodSymlink, result.Method)
				info, err := os.Lstat(result.FinalPath)
				require.NoError(t, err)
				assert.NotZero(t, info.Mode()&os.ModeSymlink)
			} else {
				assert.Equal(t, MethodCopy, result.Method)
				info, err := os.Lstat(result.FinalPath)
				require.NoError(t, err)
				assert.True(t, info.Mode().IsRegular())
			}
		})
	}
}

func TestFetchNotFoundDoesNotCreateTargetDirectory(t *testing.T) {
	_, dataConfig := writeFetchSource(t, "runtime")
	dataConfig.FileName = "missing.grb2"
	targetDir := filepath.Join(t.TempDir(), "not-created")

	_, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodCopy), dataConfig, nil)
	assert.ErrorIs(t, err, ErrDataNotFound)
	_, statErr := os.Stat(targetDir)
	assert.True(t, os.IsNotExist(statErr))
}

func TestFetchRejectsEmptyAndDirectorySources(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, dataConfig := writeFetchSource(t, "runtime")
		path := filepath.Join(dataConfig.Paths[0].Path, "data.grb2")
		require.NoError(t, os.WriteFile(path, nil, 0644))
		_, err := Fetch(context.Background(), testFetchConfig(t.TempDir(), MethodCopy), dataConfig, nil)
		assert.ErrorIs(t, err, filecheck.ErrEmpty)
	})

	t.Run("directory", func(t *testing.T) {
		_, dataConfig := writeFetchSource(t, "runtime")
		path := filepath.Join(dataConfig.Paths[0].Path, "data.grb2")
		require.NoError(t, os.Remove(path))
		require.NoError(t, os.Mkdir(path, 0755))
		_, err := Fetch(context.Background(), testFetchConfig(t.TempDir(), MethodCopy), dataConfig, nil)
		assert.ErrorIs(t, err, filecheck.ErrNotRegular)
	})
}

func TestFetchRejectsGrowingSource(t *testing.T) {
	_, dataConfig := writeFetchSource(t, "runtime")
	path := filepath.Join(dataConfig.Paths[0].Path, "data.grb2")
	config := testFetchConfig(t.TempDir(), MethodCopy)
	config.StableInterval = 200 * time.Millisecond
	go func() {
		time.Sleep(20 * time.Millisecond)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			_, _ = file.WriteString("growth")
			_ = file.Close()
		}
	}()

	_, err := Fetch(context.Background(), config, dataConfig, nil)
	assert.ErrorIs(t, err, filecheck.ErrNotStable)
	assert.NoError(t, assertNoPartialFiles(config.TargetDir))
}

func TestFetchTargetConflicts(t *testing.T) {
	source, dataConfig := writeFetchSource(t, "runtime")
	sourceInfo, err := os.Stat(source)
	require.NoError(t, err)

	t.Run("different regular file", func(t *testing.T) {
		targetDir := t.TempDir()
		target := filepath.Join(targetDir, "data.grb2")
		require.NoError(t, os.WriteFile(target, []byte("different"), 0644))
		_, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodCopy), dataConfig, nil)
		assert.ErrorIs(t, err, ErrTargetConflict)
		content, readErr := os.ReadFile(target)
		require.NoError(t, readErr)
		assert.Equal(t, []byte("different"), content)
	})

	t.Run("regular file does not satisfy explicit symlink", func(t *testing.T) {
		targetDir := t.TempDir()
		target := filepath.Join(targetDir, "data.grb2")
		content, err := os.ReadFile(source)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(target, content, 0644))
		require.NoError(t, os.Chtimes(target, sourceInfo.ModTime(), sourceInfo.ModTime()))
		_, err = Fetch(context.Background(), testFetchConfig(targetDir, MethodSymlink), dataConfig, nil)
		assert.ErrorIs(t, err, ErrTargetConflict)
	})

	if runtime.GOOS != "windows" {
		t.Run("symlink does not satisfy copy", func(t *testing.T) {
			targetDir := t.TempDir()
			target := filepath.Join(targetDir, "data.grb2")
			require.NoError(t, os.Symlink(source, target))
			_, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodCopy), dataConfig, nil)
			assert.ErrorIs(t, err, ErrTargetConflict)
		})

		t.Run("auto storage rejects symlink", func(t *testing.T) {
			_, storageConfig := writeFetchSource(t, "storage")
			storageSource := filepath.Join(storageConfig.Paths[0].Path, "data.grb2")
			targetDir := t.TempDir()
			target := filepath.Join(targetDir, "data.grb2")
			require.NoError(t, os.Symlink(storageSource, target))
			_, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodAuto), storageConfig, nil)
			assert.ErrorIs(t, err, ErrTargetConflict)
		})
	}
}

func TestFetchReuseComparesModTimeAtUnixSeconds(t *testing.T) {
	source, dataConfig := writeFetchSource(t, "storage")
	targetDir := t.TempDir()
	target := filepath.Join(targetDir, "data.grb2")
	content, err := os.ReadFile(source)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(target, content, 0644))
	sourceInfo, err := os.Stat(source)
	require.NoError(t, err)
	targetTime := time.Unix(sourceInfo.ModTime().Unix(), 999999999)
	require.NoError(t, os.Chtimes(target, targetTime, targetTime))

	result, err := Fetch(context.Background(), testFetchConfig(targetDir, MethodAuto), dataConfig, nil)
	require.NoError(t, err)
	assert.True(t, result.Reused)
}

func TestFetchCancellationLeavesNoPublishedTarget(t *testing.T) {
	_, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	config := testFetchConfig(targetDir, MethodCopy)
	config.StableInterval = time.Hour
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	_, err := Fetch(ctx, config, dataConfig, nil)
	assert.True(t, errors.Is(err, context.Canceled))
	_, statErr := os.Stat(filepath.Join(targetDir, "data.grb2"))
	assert.True(t, os.IsNotExist(statErr))
	assert.NoError(t, assertNoPartialFiles(targetDir))
}

func TestFetchCancellationDuringCopyCleansTemporaryFile(t *testing.T) {
	_, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.StandardLogger()
	oldHooks := logger.Hooks
	logger.Hooks = log.LevelHooks{}
	logger.Hooks.Add(&cancelOnCopyHook{cancel: cancel})
	defer func() { logger.Hooks = oldHooks }()

	_, err := Fetch(ctx, testFetchConfig(targetDir, MethodCopy), dataConfig, nil)
	assert.ErrorIs(t, err, context.Canceled)
	_, statErr := os.Stat(filepath.Join(targetDir, "data.grb2"))
	assert.True(t, os.IsNotExist(statErr))
	assert.NoError(t, assertNoPartialFiles(targetDir))
}

type cancelOnCopyHook struct {
	cancel context.CancelFunc
	once   sync.Once
}

func (h *cancelOnCopyHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *cancelOnCopyHook) Fire(entry *log.Entry) error {
	if strings.Contains(entry.Message, "fetch copying to temporary file") {
		h.once.Do(h.cancel)
	}
	return nil
}

func TestTargetLockCancellation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "data.grb2")
	lock, err := acquireTargetLock(context.Background(), target)
	require.NoError(t, err)
	defer lock.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = acquireTargetLock(ctx, target)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestFetchConcurrentGoroutinesCopyOnce(t *testing.T) {
	_, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()
	const workers = 8
	var wg sync.WaitGroup
	results := make([]Result, workers)
	errs := make([]error, workers)
	for i := range workers {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index], errs[index] = Fetch(context.Background(), testFetchConfig(targetDir, MethodCopy), dataConfig, nil)
		}(i)
	}
	wg.Wait()
	created := 0
	for i := range workers {
		require.NoError(t, errs[i])
		if !results[i].Reused {
			created++
		}
	}
	assert.Equal(t, 1, created)
	assert.NoError(t, assertNoPartialFiles(targetDir))
}

func TestFetchHelperProcess(t *testing.T) {
	if os.Getenv("NWPC_FETCH_HELPER") != "1" {
		return
	}
	dataConfig := dataconfig.DataConfig{
		FileName: "data.grb2",
		Paths: []dataconfig.PathItem{{
			PathType: "local", LevelType: "runtime", Path: os.Getenv("NWPC_FETCH_SOURCE_DIR"),
		}},
	}
	result, err := Fetch(context.Background(), testFetchConfig(os.Getenv("NWPC_FETCH_TARGET_DIR"), MethodCopy), dataConfig, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("reused=%t\n", result.Reused)
}

func TestFetchMultipleProcessesCopyOnce(t *testing.T) {
	_, dataConfig := writeFetchSource(t, "runtime")
	targetDir := t.TempDir()
	const workers = 4
	commands := make([]*exec.Cmd, workers)
	outputs := make([]*bytes.Buffer, workers)
	for i := range workers {
		output := &bytes.Buffer{}
		cmd := exec.Command(os.Args[0], "-test.run=^TestFetchHelperProcess$")
		cmd.Env = append(os.Environ(),
			"NWPC_FETCH_HELPER=1",
			"NWPC_FETCH_SOURCE_DIR="+dataConfig.Paths[0].Path,
			"NWPC_FETCH_TARGET_DIR="+targetDir,
		)
		cmd.Stdout = output
		cmd.Stderr = output
		require.NoError(t, cmd.Start())
		commands[i] = cmd
		outputs[i] = output
	}

	copyProcesses := 0
	for i, cmd := range commands {
		require.NoError(t, cmd.Wait(), "helper failed: %s", outputs[i].String())
		if strings.Contains(outputs[i].String(), "fetch copying to temporary file") {
			copyProcesses++
		}
	}
	assert.Equal(t, 1, copyProcesses)
	assert.NoError(t, assertNoPartialFiles(targetDir))
}

func assertNoPartialFiles(dir string) error {
	matches, err := filepath.Glob(filepath.Join(dir, "*.nwpc-fetch-partial"))
	if err != nil {
		return err
	}
	if len(matches) != 0 {
		return fmt.Errorf("partial files remain: %v", matches)
	}
	return nil
}
