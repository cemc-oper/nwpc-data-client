package fetcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/filecheck"
	"github.com/cemc-oper/nwpc-data-client/internal/finder"
	log "github.com/sirupsen/logrus"
)

// Fetch resolves the first existing local candidate and materializes it. Only
// an explicit not-exist result advances to the next candidate.
func Fetch(ctx context.Context, config Config, dataConfig dataconfig.DataConfig, locationLevels []string) (Result, error) {
	if err := config.Validate(); err != nil {
		return Result{}, err
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	for _, candidate := range finder.CandidatePaths(dataConfig, locationLevels) {
		if candidate.PathType != "local" {
			return Result{}, fmt.Errorf("%w %q: %s", ErrUnsupportedSource, candidate.PathType, candidate.Path)
		}
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
		return fetchSelected(ctx, config, candidate)
	}

	return Result{}, fmt.Errorf("%w for the supplied data identity and location levels", ErrDataNotFound)
}

func fetchSelected(ctx context.Context, config Config, source dataconfig.PathItem) (Result, error) {
	name := filepath.Base(source.Path)
	if name == "" || name == "." || name == ".." || name == "/" || name == "\\" {
		return Result{}, fmt.Errorf("invalid source basename %q", name)
	}
	targetDir, err := filepath.Abs(config.TargetDir)
	if err != nil {
		return Result{}, fmt.Errorf("resolve target directory %s: %w", config.TargetDir, err)
	}
	targetPath := filepath.Join(targetDir, name)

	log.Infof("fetch selected local source: %s", source.Path)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return Result{}, fmt.Errorf("create target directory %s: %w", targetDir, err)
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

	reused, reusedMethod, err := inspectTarget(targetPath, config.Method, source, sourceInfo)
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
