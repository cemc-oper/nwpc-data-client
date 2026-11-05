package checker

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitForStableFileAlreadyStable(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "data.grb2")
	content := []byte("some data content")
	err := os.WriteFile(filePath, content, 0644)
	assert.NoError(t, err)

	size, _, err := waitForStableFile(context.Background(), filePath, 10*time.Millisecond, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(content)), size)
}

func TestWaitForStableFileStableAfterGrowth(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "data.grb2")
	err := os.WriteFile(filePath, []byte("chunk0"), 0644)
	assert.NoError(t, err)

	// Grow the file during the first checks, then leave it alone.
	// The appends must finish well within one check interval: with the
	// sleep matching the checker's 10 ms interval, the checker can sample
	// the same size twice between two appends (observed on loaded Windows
	// CI runners) and declare the file stable before all writes landed.
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 4; i++ {
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			_, _ = f.WriteString("chunk")
			_ = f.Close()
			time.Sleep(1 * time.Millisecond)
		}
		<-stop
	}()

	size, _, err := waitForStableFile(context.Background(), filePath, 10*time.Millisecond, 50)
	close(stop)
	<-done

	assert.NoError(t, err)
	info, statErr := os.Stat(filePath)
	assert.NoError(t, statErr)
	assert.Equal(t, info.Size(), size)
	assert.Equal(t, int64(len("chunk0")+4*len("chunk")), size)
}

func TestWaitForStableFileKeepsGrowing(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "data.grb2")
	err := os.WriteFile(filePath, []byte("chunk0"), 0644)
	assert.NoError(t, err)

	// Keep growing the file so it never stabilizes within the check budget.
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
			}
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			_, _ = f.WriteString("chunk")
			_ = f.Close()
			time.Sleep(5 * time.Millisecond)
		}
	}()

	_, _, err = waitForStableFile(context.Background(), filePath, 10*time.Millisecond, 5)
	close(stop)
	<-done

	assert.Error(t, err)
}

func TestWaitForStableFileMissing(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "not-exist.grb2")

	_, _, err := waitForStableFile(context.Background(), filePath, 10*time.Millisecond, 5)

	assert.Error(t, err)
}

func TestGetFileSize(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "data.grb2")
	content := []byte("some data content")
	assert.NoError(t, os.WriteFile(filePath, content, 0644))

	size, err := getFileSize(filePath)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(content)), size)

	_, err = getFileSize(filepath.Join(dir, "not-exist.grb2"))
	assert.Error(t, err)
}

func TestWaitForStableFileCancellation(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "data.grb2")
	require.NoError(t, os.WriteFile(filePath, []byte("data"), 0644))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, _, err := waitForStableFile(ctx, filePath, time.Hour, 10)
		done <- err
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("stable-file wait did not return after cancellation")
	}
}
