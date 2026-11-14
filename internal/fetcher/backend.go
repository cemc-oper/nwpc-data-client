package fetcher

import (
	"context"
	"os"
	"time"
)

// RemoteMetadata is the source identity used for verification and idempotent
// target reuse. ModTime is compared at Unix-second precision.
type RemoteMetadata struct {
	Size    int64
	ModTime time.Time
	Regular bool
}

// RemoteBackend reads one remote file without deciding its final local name.
type RemoteBackend interface {
	Stat(ctx context.Context, remotePath string) (RemoteMetadata, error)
	Fetch(ctx context.Context, remotePath string, localTempPath string) error
}

type fileMetadata struct {
	Size    int64
	ModTime time.Time
}

func metadataFromFileInfo(info os.FileInfo) fileMetadata {
	return fileMetadata{Size: info.Size(), ModTime: info.ModTime()}
}
