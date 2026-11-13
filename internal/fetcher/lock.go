package fetcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type targetLock struct {
	file *os.File
}

func acquireTargetLock(ctx context.Context, targetPath string) (*targetLock, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	lockPath := filepath.Join(filepath.Dir(targetPath), "."+filepath.Base(targetPath)+".nwpc-fetch.lock")
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	lock := &targetLock{file: file}
	for {
		locked, err := tryLockFile(file)
		if err != nil {
			_ = file.Close()
			return nil, err
		}
		if locked {
			return lock, nil
		}
		timer := time.NewTimer(100 * time.Millisecond)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			_ = file.Close()
			return nil, ctx.Err()
		}
	}
}

func (l *targetLock) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	unlockErr := unlockFile(l.file)
	closeErr := l.file.Close()
	l.file = nil
	if unlockErr != nil || closeErr != nil {
		return fmt.Errorf("release target lock: %w", errors.Join(unlockErr, closeErr))
	}
	return nil
}
