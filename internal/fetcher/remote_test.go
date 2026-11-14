package fetcher

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/filecheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRemoteBackend struct {
	mu              sync.Mutex
	metadata        RemoteMetadata
	statErr         error
	afterMetadata   *RemoteMetadata
	fetchErr        error
	statCalls       int
	fetchCalls      int
	fetchStarted    chan struct{}
	delay           time.Duration
	fetchWriteSize  int64
	fetchWriteMtime time.Time
}

func (f *fakeRemoteBackend) Stat(ctx context.Context, remotePath string) (RemoteMetadata, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.statCalls++
	if f.statErr != nil {
		return RemoteMetadata{}, f.statErr
	}
	if f.statCalls >= 2 && f.afterMetadata != nil {
		return *f.afterMetadata, nil
	}
	return f.metadata, nil
}

func (f *fakeRemoteBackend) Fetch(ctx context.Context, remotePath string, localTempPath string) error {
	f.mu.Lock()
	f.fetchCalls++
	f.mu.Unlock()
	if f.fetchStarted != nil {
		select {
		case f.fetchStarted <- struct{}{}:
		default:
		}
	}
	if f.delay > 0 {
		timer := time.NewTimer(f.delay)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if f.fetchErr != nil {
		return f.fetchErr
	}
	size := f.metadata.Size
	if f.fetchWriteSize != 0 {
		size = f.fetchWriteSize
	}
	mtime := f.metadata.ModTime
	if !f.fetchWriteMtime.IsZero() {
		mtime = f.fetchWriteMtime
	}
	if err := os.WriteFile(localTempPath, bytes.Repeat([]byte("r"), int(size)), 0644); err != nil {
		return err
	}
	return os.Chtimes(localTempPath, mtime, mtime)
}

func remoteMetadata() RemoteMetadata {
	return RemoteMetadata{Size: 47, ModTime: time.Unix(1653633373, 0), Regular: true}
}

func remoteDataConfig() dataconfig.DataConfig {
	return dataconfig.DataConfig{
		FileName: "data.grb2",
		Paths: []dataconfig.PathItem{{
			PathType: "cmafs", LevelType: "storage", Path: "/sstorage/COMMONDATA",
		}},
	}
}

func TestFetchRemoteFallbackAndVerification(t *testing.T) {
	backend := &fakeRemoteBackend{metadata: remoteMetadata()}
	dataConfig := remoteDataConfig()
	dataConfig.Paths = append([]dataconfig.PathItem{{
		PathType: "local", LevelType: "storage", Path: filepath.Join(t.TempDir(), "missing"),
	}}, dataConfig.Paths...)
	config := testFetchConfig(t.TempDir(), MethodAuto)
	config.CmafsBackend = backend

	result, err := Fetch(context.Background(), config, dataConfig, nil)
	require.NoError(t, err)
	assert.Equal(t, "cmafs", result.SourceType)
	assert.Equal(t, MethodCopy, result.Method)
	assert.Equal(t, 2, backend.statCalls)
	assert.Equal(t, 1, backend.fetchCalls)
	content, err := os.ReadFile(result.FinalPath)
	require.NoError(t, err)
	assert.Len(t, content, 47)
	assert.NoError(t, assertNoPartialFiles(filepath.Dir(result.FinalPath)))
}

func TestFetchRemoteNotExistFallsBackToLocal(t *testing.T) {
	source, localConfig := writeFetchSource(t, "storage")
	backend := &fakeRemoteBackend{metadata: remoteMetadata(), statErr: ErrSourceNotExist}
	dataConfig := remoteDataConfig()
	dataConfig.Paths = append(dataConfig.Paths, localConfig.Paths...)
	config := testFetchConfig(t.TempDir(), MethodAuto)
	config.CmafsBackend = backend

	result, err := Fetch(context.Background(), config, dataConfig, nil)
	require.NoError(t, err)
	assert.Equal(t, "local", result.SourceType)
	assert.Equal(t, source, result.SourcePath)
	assert.Equal(t, 1, backend.statCalls)
	assert.Zero(t, backend.fetchCalls)
}

func TestFetchRemoteErrorDoesNotFallBack(t *testing.T) {
	_, localConfig := writeFetchSource(t, "storage")
	backend := &fakeRemoteBackend{metadata: remoteMetadata(), statErr: ErrBackendOperation}
	dataConfig := remoteDataConfig()
	dataConfig.Paths = append(dataConfig.Paths, localConfig.Paths...)
	config := testFetchConfig(t.TempDir(), MethodCopy)
	config.CmafsBackend = backend

	_, err := Fetch(context.Background(), config, dataConfig, nil)
	assert.ErrorIs(t, err, ErrBackendOperation)
	assert.Zero(t, backend.fetchCalls)
}

func TestFetchRemoteRequiresBackend(t *testing.T) {
	_, err := Fetch(context.Background(), testFetchConfig(t.TempDir(), MethodCopy), remoteDataConfig(), nil)
	assert.ErrorIs(t, err, ErrBackendNotConfigured)
}

func TestFetchRemoteRejectsExplicitSymlinkBeforeTransfer(t *testing.T) {
	backend := &fakeRemoteBackend{metadata: remoteMetadata()}
	config := testFetchConfig(t.TempDir(), MethodSymlink)
	config.CmafsBackend = backend

	_, err := Fetch(context.Background(), config, remoteDataConfig(), nil)
	assert.ErrorIs(t, err, ErrMethodUnavailable)
	assert.Equal(t, 1, backend.statCalls)
	assert.Zero(t, backend.fetchCalls)
}

func TestFetchRemoteRejectsEmptyAndNonRegular(t *testing.T) {
	for name, metadata := range map[string]RemoteMetadata{
		"empty":     {Size: 0, ModTime: time.Unix(1, 0), Regular: true},
		"directory": {Size: 47, ModTime: time.Unix(1, 0), Regular: false},
	} {
		t.Run(name, func(t *testing.T) {
			backend := &fakeRemoteBackend{metadata: metadata}
			config := testFetchConfig(t.TempDir(), MethodCopy)
			config.CmafsBackend = backend
			_, err := Fetch(context.Background(), config, remoteDataConfig(), nil)
			if name == "empty" {
				assert.ErrorIs(t, err, filecheck.ErrEmpty)
			} else {
				assert.ErrorIs(t, err, filecheck.ErrNotRegular)
			}
			assert.Zero(t, backend.fetchCalls)
		})
	}
}

func TestFetchRemoteSourceChangedDuringTransfer(t *testing.T) {
	changed := remoteMetadata()
	changed.Size++
	backend := &fakeRemoteBackend{metadata: remoteMetadata(), afterMetadata: &changed}
	config := testFetchConfig(t.TempDir(), MethodCopy)
	config.CmafsBackend = backend

	_, err := Fetch(context.Background(), config, remoteDataConfig(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "changed during transfer")
	_, statErr := os.Stat(filepath.Join(config.TargetDir, "data.grb2"))
	assert.True(t, os.IsNotExist(statErr))
	assert.NoError(t, assertNoPartialFiles(config.TargetDir))
}

func TestFetchRemoteDownloadVerification(t *testing.T) {
	t.Run("wrong size", func(t *testing.T) {
		backend := &fakeRemoteBackend{metadata: remoteMetadata(), fetchWriteSize: 2}
		config := testFetchConfig(t.TempDir(), MethodCopy)
		config.CmafsBackend = backend
		_, err := Fetch(context.Background(), config, remoteDataConfig(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "size/mtime/type")
		assert.NoError(t, assertNoPartialFiles(config.TargetDir))
	})
	t.Run("wrong mtime", func(t *testing.T) {
		backend := &fakeRemoteBackend{metadata: remoteMetadata(), fetchWriteMtime: time.Unix(1, 0)}
		config := testFetchConfig(t.TempDir(), MethodCopy)
		config.CmafsBackend = backend
		_, err := Fetch(context.Background(), config, remoteDataConfig(), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "size/mtime/type")
		assert.NoError(t, assertNoPartialFiles(config.TargetDir))
	})
}

func TestFetchRemoteReusesTargetWithoutDownload(t *testing.T) {
	metadata := remoteMetadata()
	backend := &fakeRemoteBackend{metadata: metadata}
	config := testFetchConfig(t.TempDir(), MethodAuto)
	config.CmafsBackend = backend
	target := filepath.Join(config.TargetDir, "data.grb2")
	require.NoError(t, os.WriteFile(target, bytes.Repeat([]byte("r"), int(metadata.Size)), 0644))
	require.NoError(t, os.Chtimes(target, metadata.ModTime, metadata.ModTime))

	result, err := Fetch(context.Background(), config, remoteDataConfig(), nil)
	require.NoError(t, err)
	assert.True(t, result.Reused)
	assert.Zero(t, backend.fetchCalls)
	assert.Equal(t, 1, backend.statCalls)
}

func TestFetchRemoteAutoDoesNotReuseSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks may require additional Windows privileges")
	}
	backend := &fakeRemoteBackend{metadata: remoteMetadata()}
	config := testFetchConfig(t.TempDir(), MethodAuto)
	config.CmafsBackend = backend
	dataConfig := remoteDataConfig()
	dataConfig.Paths[0].LevelType = "runtime"
	target := filepath.Join(config.TargetDir, "data.grb2")
	require.NoError(t, os.Symlink("/sstorage/COMMONDATA/data.grb2", target))

	_, err := Fetch(context.Background(), config, dataConfig, nil)
	assert.ErrorIs(t, err, ErrTargetConflict)
	assert.Zero(t, backend.fetchCalls)
}

func TestFetchRemoteConcurrentDownloadsOnce(t *testing.T) {
	backend := &fakeRemoteBackend{metadata: remoteMetadata(), delay: 30 * time.Millisecond}
	config := testFetchConfig(t.TempDir(), MethodCopy)
	config.CmafsBackend = backend
	const workers = 4
	var wg sync.WaitGroup
	errs := make([]error, workers)
	results := make([]Result, workers)
	for i := range workers {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index], errs[index] = Fetch(context.Background(), config, remoteDataConfig(), nil)
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
	assert.Equal(t, 1, backend.fetchCalls)
}

func TestFetchRemoteCancellationCleansTemporaryFile(t *testing.T) {
	backend := &fakeRemoteBackend{
		metadata:     remoteMetadata(),
		delay:        time.Hour,
		fetchStarted: make(chan struct{}, 1),
	}
	config := testFetchConfig(t.TempDir(), MethodCopy)
	config.CmafsBackend = backend
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-backend.fetchStarted
		cancel()
	}()

	_, err := Fetch(ctx, config, remoteDataConfig(), nil)
	assert.True(t, errors.Is(err, context.Canceled))
	assert.NoError(t, assertNoPartialFiles(config.TargetDir))
	_, statErr := os.Stat(filepath.Join(config.TargetDir, "data.grb2"))
	assert.True(t, os.IsNotExist(statErr))
}
