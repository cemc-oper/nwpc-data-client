package checker

import (
	"context"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/filecheck"
)

func getFileSize(filePath string) (int64, error) {
	return filecheck.GetFileSize(filePath)
}

func waitForStableFile(ctx context.Context, path string, interval time.Duration, maxCount int) (int64, int, error) {
	return filecheck.WaitForStableFile(ctx, path, interval, maxCount)
}
