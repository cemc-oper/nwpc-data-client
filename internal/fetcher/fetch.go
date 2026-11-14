package fetcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/filecheck"
	"github.com/cemc-oper/nwpc-data-client/internal/finder"
	log "github.com/sirupsen/logrus"
)

// Fetch resolves the first existing candidate and materializes it. Only an
// explicit not-exist result advances to the next candidate.
func Fetch(ctx context.Context, config Config, dataConfig dataconfig.DataConfig, locationLevels []string) (Result, error) {
	if err := config.Validate(); err != nil {
		return Result{}, err
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	for _, candidate := range finder.CandidatePaths(dataConfig, locationLevels) {
		switch candidate.PathType {
		case "local":
			info, err := os.Stat(candidate.Path)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return Result{}, fmt.Errorf("stat local source %s: %w", candidate.Path, err)
			}
			if !info.Mode().IsRegular() {
				return Result{}, fmt.Errorf("stat local source %s: %w", candidate.Path, filecheck.ErrNotRegular)
			}
			if info.Size() == 0 {
				return Result{}, fmt.Errorf("stat local source %s: %w", candidate.Path, filecheck.ErrEmpty)
			}

			sourcePath, err := filepath.Abs(candidate.Path)
			if err != nil {
				return Result{}, fmt.Errorf("resolve source path %s: %w", candidate.Path, err)
			}
			candidate.Path = filepath.Clean(sourcePath)
			return fetchLocalSelected(ctx, config, candidate)

		case "cmafs":
			if config.CmafsBackend == nil {
				return Result{}, fmt.Errorf("%w: cmafs backend is required for %s", ErrBackendNotConfigured, candidate.Path)
			}
			candidate.Path = path.Clean(candidate.Path)
			metadata, err := config.CmafsBackend.Stat(ctx, candidate.Path)
			if errors.Is(err, ErrSourceNotExist) {
				continue
			}
			if err != nil {
				return Result{}, fmt.Errorf("stat cmafs source %s: %w", candidate.Path, err)
			}
			if !metadata.Regular {
				return Result{}, fmt.Errorf("stat cmafs source %s: %w", candidate.Path, filecheck.ErrNotRegular)
			}
			if metadata.Size == 0 {
				return Result{}, fmt.Errorf("stat cmafs source %s: %w", candidate.Path, filecheck.ErrEmpty)
			}
			if config.Method == MethodSymlink {
				return Result{}, fmt.Errorf("%w: cmafs source %s cannot satisfy explicit symlink mode", ErrMethodUnavailable, candidate.Path)
			}
			return fetchRemoteSelected(ctx, config, candidate, metadata)

		default:
			return Result{}, fmt.Errorf("%w %q: %s", ErrUnsupportedSource, candidate.PathType, candidate.Path)
		}
	}

	return Result{}, fmt.Errorf("%w for the supplied data identity and location levels", ErrDataNotFound)
}

func fetchLocalSelected(ctx context.Context, config Config, source dataconfig.PathItem) (Result, error) {
	targetPath, err := targetPathForSource(config, source)
	if err != nil {
		return Result{}, err
	}
	log.Infof("fetch selected local source: %s", source.Path)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return Result{}, fmt.Errorf("create target directory %s: %w", filepath.Dir(targetPath), err)
	}

	lock, err := acquireTargetLock(ctx, targetPath)
	if err != nil {
		return Result{}, fmt.Errorf("acquire target lock for %s: %w", targetPath, err)
	}
	defer func() {
		_ = lock.Close()
	}()

	log.Infof("checking local source stability: %s", source.Path)
	sourceFile, sourceInfo, err := filecheck.OpenStableRegularFile(ctx, source.Path, config.StableInterval)
	if err != nil {
		return Result{}, fmt.Errorf("check local source %s: %w", source.Path, err)
	}
	defer func() {
		_ = sourceFile.Close()
	}()

	reused, reusedMethod, err := inspectTarget(targetPath, config.Method, source, metadataFromFileInfo(sourceInfo))
	if err != nil {
		return Result{}, err
	}
	if reused {
		log.Infof("fetch reusing existing target: %s", targetPath)
		return newResult(targetPath, source, sourceInfo.Size(), reusedMethod, true), nil
	}

	method := config.Method
	if method == MethodAuto {
		method = MethodSymlink
		if hasStorageLevel(source.LevelType) {
			method = MethodCopy
		}
	}

	switch method {
	case MethodSymlink:
		reused, err = publishSymlink(source, sourceInfo, targetPath, config.Method)
	case MethodCopy:
		reused, err = copyAndPublish(ctx, sourceFile, sourceInfo, source, targetPath, config.Method)
	default:
		err = fmt.Errorf("unsupported effective method %q", method)
	}
	if err != nil {
		return Result{}, err
	}
	return newResult(targetPath, source, sourceInfo.Size(), method, reused), nil
}

func fetchRemoteSelected(ctx context.Context, config Config, source dataconfig.PathItem, metadata RemoteMetadata) (Result, error) {
	targetPath, err := targetPathForSource(config, source)
	if err != nil {
		return Result{}, err
	}
	log.Infof("fetch selected cmafs source: %s", source.Path)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return Result{}, fmt.Errorf("create target directory %s: %w", filepath.Dir(targetPath), err)
	}

	lock, err := acquireTargetLock(ctx, targetPath)
	if err != nil {
		return Result{}, fmt.Errorf("acquire target lock for %s: %w", targetPath, err)
	}
	defer func() {
		_ = lock.Close()
	}()

	sourceMetadata := fileMetadata{Size: metadata.Size, ModTime: metadata.ModTime}
	reused, _, err := inspectTarget(targetPath, config.Method, source, sourceMetadata)
	if err != nil {
		return Result{}, err
	}
	if reused {
		log.Infof("fetch reusing existing target: %s", targetPath)
		return newResult(targetPath, source, metadata.Size, MethodCopy, true), nil
	}

	temp, err := os.CreateTemp(filepath.Dir(targetPath), "."+filepath.Base(targetPath)+".*.nwpc-fetch-partial")
	if err != nil {
		return Result{}, fmt.Errorf("create cmafs temporary target: %w", err)
	}
	tempPath := temp.Name()
	if err := temp.Close(); err != nil {
		_ = os.Remove(tempPath)
		return Result{}, fmt.Errorf("close cmafs temporary target %s: %w", tempPath, err)
	}
	defer func() {
		if tempPath != "" {
			_ = os.Remove(tempPath)
		}
	}()

	log.Infof("fetch downloading cmafs source to temporary file: %s", tempPath)
	if err := config.CmafsBackend.Fetch(ctx, source.Path, tempPath); err != nil {
		return Result{}, fmt.Errorf("fetch cmafs source %s: %w", source.Path, err)
	}
	localInfo, err := os.Stat(tempPath)
	if err != nil {
		return Result{}, fmt.Errorf("inspect downloaded temporary target %s: %w", tempPath, err)
	}
	if !localInfo.Mode().IsRegular() || localInfo.Size() == 0 || localInfo.Size() != metadata.Size ||
		localInfo.ModTime().Unix() != metadata.ModTime.Unix() {
		return Result{}, fmt.Errorf("downloaded temporary target %s failed size/mtime/type verification", tempPath)
	}

	after, err := config.CmafsBackend.Stat(ctx, source.Path)
	if err != nil {
		return Result{}, fmt.Errorf("recheck cmafs source %s after transfer: %w", source.Path, err)
	}
	if !after.Regular || after.Size != metadata.Size || after.ModTime.Unix() != metadata.ModTime.Unix() {
		return Result{}, fmt.Errorf("cmafs source %s changed during transfer", source.Path)
	}

	reused, err = publishRegularFile(tempPath, targetPath, source, sourceMetadata, config.Method)
	if err != nil {
		return Result{}, err
	}
	if err := os.Remove(tempPath); err != nil {
		return Result{}, fmt.Errorf("remove cmafs temporary target %s: %w", tempPath, err)
	}
	tempPath = ""
	return newResult(targetPath, source, metadata.Size, MethodCopy, reused), nil
}

// targetPathForSource resolves the published path. The name defaults to the
// source basename and is replaced by the rendered Config.TargetName template
// when one is supplied.
func targetPathForSource(config Config, source dataconfig.PathItem) (string, error) {
	sourceName := filepath.Base(source.Path)
	if source.PathType == "cmafs" {
		sourceName = path.Base(source.Path)
	}
	if sourceName == "" || sourceName == "." || sourceName == ".." || sourceName == "/" || sourceName == "\\" {
		return "", fmt.Errorf("invalid source basename %q", sourceName)
	}

	name, err := dataconfig.RenderTargetName(config.TargetName, config.TemplateVariable, sourceName)
	if err != nil {
		return "", err
	}
	if err := validateTargetName(name); err != nil {
		return "", err
	}

	absoluteTargetDir, err := filepath.Abs(config.TargetDir)
	if err != nil {
		return "", fmt.Errorf("resolve target directory %s: %w", config.TargetDir, err)
	}
	return filepath.Join(absoluteTargetDir, name), nil
}

// validateTargetName keeps the published file directly inside the target
// directory: the name must be a single file name, never a path.
func validateTargetName(name string) error {
	if name == "" {
		return fmt.Errorf("target name must not be empty")
	}
	if strings.ContainsAny(name, "/\\\x00\r\n") {
		return fmt.Errorf("target name %q must be a file name without path separators, NUL or newlines", name)
	}
	if name == "." || name == ".." {
		return fmt.Errorf("invalid target name %q", name)
	}
	if filepath.IsAbs(name) || filepath.Base(name) != name {
		return fmt.Errorf("target name %q must not be a path", name)
	}
	return nil
}

func newResult(targetPath string, source dataconfig.PathItem, size int64, method Method, reused bool) Result {
	return Result{
		FinalPath:     targetPath,
		SourcePath:    source.Path,
		SourceType:    source.PathType,
		LocationLevel: source.LevelType,
		Method:        method,
		Size:          size,
		Reused:        reused,
	}
}
