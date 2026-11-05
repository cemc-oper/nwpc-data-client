package checker

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return -100, fmt.Errorf("get file info has error:%v", err)
	}
	return fileInfo.Size(), nil
}

// waitForStableFile polls path until its size is non-zero and unchanged
// between two consecutive checks, and returns the final size. It gives up
// with an error after maxCount checks or when ctx is cancelled.
func waitForStableFile(ctx context.Context, path string, interval time.Duration, maxCount int) (int64, int, error) {
	var lastSize int64 = -1
	for i := 0; i < maxCount; i++ {
		if err := ctx.Err(); err != nil {
			return 0, i, err
		}
		currentSize, err := getFileSize(path)
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
