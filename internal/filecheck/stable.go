// Package filecheck provides local file readiness checks shared by checker
// polling and one-shot fetch materialization.
package filecheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNotRegular = errors.New("not a regular file")
	ErrEmpty      = errors.New("file is empty")
	ErrNotStable  = errors.New("file is not stable")
)

// GetFileSize returns a file's current size.
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return -100, fmt.Errorf("get file info has error:%v", err)
	}
	return fileInfo.Size(), nil
}

// WaitForStableFile polls path until its size is non-zero and unchanged
// between two consecutive checks. It preserves the checker's historical
// logging and check-budget semantics.
func WaitForStableFile(ctx context.Context, path string, interval time.Duration, maxCount int) (int64, int, error) {
	var lastSize int64 = -1
	for i := range maxCount {
		if err := ctx.Err(); err != nil {
			return 0, i, err
		}
		currentSize, err := GetFileSize(path)
		if err != nil {
			log.Warnf("checking size...get file size failed: %v", err)
			lastSize = -1
		} else if currentSize > 0 && currentSize == lastSize {
			log.Infof("checking size...success %d/%d", i, maxCount)
			return currentSize, i + 1, nil
		} else if currentSize == 0 {
			log.Warnf("get size 0, continue to wait for a valid data size %d/%d", i, maxCount)
			lastSize = currentSize
		} else {
			log.Infof("checking size...changed %d/%d", i, maxCount)
			lastSize = currentSize
		}
		if err := sleepContext(ctx, interval); err != nil {
			return 0, i + 1, err
		}
	}
	return 0, maxCount, fmt.Errorf("file is not stable after %d checks: %s", maxCount, path)
}

// OpenStableRegularFile opens path and performs both stability observations on
// the same open file description. The returned file is positioned at offset
// zero and must be closed by the caller. Size and modification time must both
// remain unchanged across one interval.
func OpenStableRegularFile(ctx context.Context, path string, interval time.Duration) (*os.File, os.FileInfo, error) {
	if interval <= 0 {
		return nil, nil, fmt.Errorf("stability interval must be greater than zero")
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	pathInfo, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}
	if !pathInfo.Mode().IsRegular() {
		return nil, nil, fmt.Errorf("inspect source %s: %w", path, ErrNotRegular)
	}
	if pathInfo.Size() == 0 {
		return nil, nil, fmt.Errorf("inspect source %s: %w", path, ErrEmpty)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	keepOpen := false
	defer func() {
		if !keepOpen {
			_ = file.Close()
		}
	}()

	first, err := stableFileInfo(file)
	if err != nil {
		return nil, nil, fmt.Errorf("inspect source %s: %w", path, err)
	}
	if err := sleepContext(ctx, interval); err != nil {
		return nil, nil, err
	}
	second, err := stableFileInfo(file)
	if err != nil {
		return nil, nil, fmt.Errorf("recheck source %s: %w", path, err)
	}
	if first.Size() != second.Size() || !first.ModTime().Equal(second.ModTime()) {
		return nil, nil, fmt.Errorf("%w: %s changed during stability observation", ErrNotStable, path)
	}

	keepOpen = true
	return file, second, nil
}

func stableFileInfo(file *os.File) (os.FileInfo, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, ErrNotRegular
	}
	if info.Size() == 0 {
		return nil, ErrEmpty
	}
	return info, nil
}

func sleepContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
