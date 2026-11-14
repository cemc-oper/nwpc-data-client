package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
)

func inspectTarget(targetPath string, requested Method, source dataconfig.PathItem, sourceMetadata fileMetadata) (bool, Method, error) {
	targetInfo, err := os.Lstat(targetPath)
	if os.IsNotExist(err) {
		return false, "", nil
	}
	if err != nil {
		return false, "", fmt.Errorf("inspect target %s: %w", targetPath, err)
	}

	if targetInfo.Mode()&os.ModeSymlink != 0 {
		allowed := requested == MethodSymlink ||
			(requested == MethodAuto && source.PathType == "local" && !hasStorageLevel(source.LevelType))
		if !allowed {
			return false, "", targetConflict(targetPath, "a symlink does not satisfy the requested method")
		}
		linkTarget, err := os.Readlink(targetPath)
		if err != nil {
			return false, "", fmt.Errorf("read target symlink %s: %w", targetPath, err)
		}
		if !filepath.IsAbs(linkTarget) {
			linkTarget = filepath.Join(filepath.Dir(targetPath), linkTarget)
		}
		if filepath.Clean(linkTarget) != source.Path {
			return false, "", targetConflict(targetPath, "symlink points to a different source")
		}
		return true, MethodSymlink, nil
	}

	if !targetInfo.Mode().IsRegular() {
		return false, "", targetConflict(targetPath, "existing target is not a regular file")
	}
	if requested != MethodCopy && requested != MethodAuto {
		return false, "", targetConflict(targetPath, "a regular file does not satisfy explicit symlink mode")
	}
	if targetInfo.Size() == 0 || targetInfo.Size() != sourceMetadata.Size ||
		targetInfo.ModTime().Unix() != sourceMetadata.ModTime.Unix() {
		return false, "", targetConflict(targetPath, "existing file size or modification time differs from the source")
	}
	return true, MethodCopy, nil
}

func publishSymlink(source dataconfig.PathItem, sourceInfo os.FileInfo, targetPath string, requested Method) (bool, error) {
	currentInfo, err := os.Stat(source.Path)
	if err != nil {
		return false, fmt.Errorf("recheck symlink source %s: %w", source.Path, err)
	}
	if !os.SameFile(sourceInfo, currentInfo) || currentInfo.Size() != sourceInfo.Size() ||
		!currentInfo.ModTime().Equal(sourceInfo.ModTime()) {
		return false, fmt.Errorf("symlink source %s changed after stability check", source.Path)
	}

	if err := os.Symlink(source.Path, targetPath); err == nil {
		log.Infof("fetch published symlink: %s", targetPath)
		return false, nil
	} else if !os.IsExist(err) {
		return false, fmt.Errorf("create target symlink %s: %w", targetPath, err)
	}

	reused, _, err := inspectTarget(targetPath, requested, source, metadataFromFileInfo(sourceInfo))
	if err != nil {
		return false, err
	}
	if !reused {
		return false, targetConflict(targetPath, "target disappeared during symlink publication; retry the request")
	}
	return true, nil
}

func copyAndPublish(ctx context.Context, sourceFile *os.File, sourceInfo os.FileInfo, source dataconfig.PathItem, targetPath string, requested Method) (bool, error) {
	targetDir := filepath.Dir(targetPath)
	temp, err := os.CreateTemp(targetDir, "."+filepath.Base(targetPath)+".*.nwpc-fetch-partial")
	if err != nil {
		return false, fmt.Errorf("create temporary target in %s: %w", targetDir, err)
	}
	tempPath := temp.Name()
	closed := false
	defer func() {
		if !closed {
			_ = temp.Close()
		}
		if tempPath != "" {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := sourceFile.Seek(0, io.SeekStart); err != nil {
		return false, fmt.Errorf("rewind source %s: %w", source.Path, err)
	}
	log.Infof("fetch copying to temporary file: %s", tempPath)
	if err := copyWithContext(ctx, temp, sourceFile); err != nil {
		return false, fmt.Errorf("copy source %s: %w", source.Path, err)
	}
	if err := temp.Chmod(sourceInfo.Mode().Perm()); err != nil {
		return false, fmt.Errorf("preserve source permissions on %s: %w", tempPath, err)
	}
	if err := os.Chtimes(tempPath, time.Now(), sourceInfo.ModTime()); err != nil {
		return false, fmt.Errorf("preserve source modification time on %s: %w", tempPath, err)
	}
	if err := temp.Sync(); err != nil {
		return false, fmt.Errorf("sync temporary target %s: %w", tempPath, err)
	}
	if err := temp.Close(); err != nil {
		return false, fmt.Errorf("close temporary target %s: %w", tempPath, err)
	}
	closed = true

	currentSourceInfo, err := sourceFile.Stat()
	if err != nil {
		return false, fmt.Errorf("recheck copied source %s: %w", source.Path, err)
	}
	if currentSourceInfo.Size() != sourceInfo.Size() || !currentSourceInfo.ModTime().Equal(sourceInfo.ModTime()) {
		return false, fmt.Errorf("source %s changed while it was being copied", source.Path)
	}
	tempInfo, err := os.Stat(tempPath)
	if err != nil {
		return false, fmt.Errorf("inspect temporary target %s: %w", tempPath, err)
	}
	if !tempInfo.Mode().IsRegular() || tempInfo.Size() == 0 || tempInfo.Size() != sourceInfo.Size() {
		return false, fmt.Errorf("temporary target %s failed size/type verification", tempPath)
	}

	reused, err := publishRegularFile(tempPath, targetPath, source, metadataFromFileInfo(sourceInfo), requested)
	if err != nil {
		return false, err
	}
	if err := os.Remove(tempPath); err != nil {
		return false, fmt.Errorf("remove temporary target %s: %w", tempPath, err)
	}
	tempPath = ""
	return reused, nil
}

func publishRegularFile(tempPath string, targetPath string, source dataconfig.PathItem, sourceMetadata fileMetadata, requested Method) (bool, error) {
	if err := os.Link(tempPath, targetPath); err == nil {
		log.Infof("fetch published regular file: %s", targetPath)
		return false, nil
	} else if !os.IsExist(err) {
		return false, fmt.Errorf("publish temporary target %s to %s without replacement: %w", tempPath, targetPath, err)
	}

	reused, _, err := inspectTarget(targetPath, requested, source, sourceMetadata)
	if err != nil {
		return false, err
	}
	if !reused {
		return false, targetConflict(targetPath, "target disappeared during publication; retry the request")
	}
	return true, nil
}

func copyWithContext(ctx context.Context, destination io.Writer, source io.Reader) error {
	buffer := make([]byte, 1024*1024)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		n, readErr := source.Read(buffer)
		if n > 0 {
			written := 0
			for written < n {
				if err := ctx.Err(); err != nil {
					return err
				}
				count, err := destination.Write(buffer[written:n])
				if err != nil {
					return err
				}
				if count == 0 {
					return io.ErrShortWrite
				}
				written += count
			}
		}
		if errors.Is(readErr, io.EOF) {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}

func targetConflict(targetPath, reason string) error {
	return fmt.Errorf("%w: %s: %s", ErrTargetConflict, targetPath, reason)
}
