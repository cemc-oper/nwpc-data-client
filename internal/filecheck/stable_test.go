package filecheck

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenStableRegularFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.grb2")
	require.NoError(t, os.WriteFile(path, []byte("data"), 0644))

	file, info, err := OpenStableRegularFile(context.Background(), path, time.Millisecond)
	require.NoError(t, err)
	defer file.Close()
	assert.Equal(t, int64(4), info.Size())
	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.Equal(t, []byte("data"), content)
}

func TestOpenStableRegularFileRejectsEmptyAndDirectory(t *testing.T) {
	dir := t.TempDir()
	empty := filepath.Join(dir, "empty")
	require.NoError(t, os.WriteFile(empty, nil, 0644))

	_, _, err := OpenStableRegularFile(context.Background(), empty, time.Millisecond)
	assert.ErrorIs(t, err, ErrEmpty)
	_, _, err = OpenStableRegularFile(context.Background(), dir, time.Millisecond)
	assert.ErrorIs(t, err, ErrNotRegular)
}

func TestOpenStableRegularFileRejectsGrowth(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.grb2")
	require.NoError(t, os.WriteFile(path, []byte("data"), 0644))
	go func() {
		time.Sleep(10 * time.Millisecond)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			_, _ = file.WriteString("more")
			_ = file.Close()
		}
	}()

	_, _, err := OpenStableRegularFile(context.Background(), path, 50*time.Millisecond)
	assert.ErrorIs(t, err, ErrNotStable)
}

func TestOpenStableRegularFileReadsSameOpenedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.grb2")
	replaced := filepath.Join(dir, "old.grb2")
	require.NoError(t, os.WriteFile(path, []byte("old"), 0644))
	go func() {
		time.Sleep(10 * time.Millisecond)
		if os.Rename(path, replaced) == nil {
			_ = os.WriteFile(path, []byte("new"), 0644)
		}
	}()

	file, _, err := OpenStableRegularFile(context.Background(), path, 50*time.Millisecond)
	require.NoError(t, err)
	defer file.Close()
	content, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.Equal(t, []byte("old"), content)
}

func TestOpenStableRegularFileCancellation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.grb2")
	require.NoError(t, os.WriteFile(path, []byte("data"), 0644))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, _, err := OpenStableRegularFile(ctx, path, time.Hour)
		done <- err
	}()
	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.True(t, errors.Is(err, context.Canceled))
	case <-time.After(time.Second):
		t.Fatal("stable open did not return after cancellation")
	}
}
