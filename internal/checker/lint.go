package checker

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// LintCheckerFile statically validates one checker YAML file. It applies the
// normal defaults and YAML field-presence rules without CLI flags, data config
// loading, stdin, data directory access, command execution, or report creation.
func LintCheckerFile(filePath string) error {
	fileConfig, err := LoadCheckerConfig(filePath)
	if err != nil {
		return err
	}
	if _, err := mergeCheckerConfig(fileConfig, CheckerConfig{}, nil); err != nil {
		return fmt.Errorf("validate checker config file %s failed: %w", filePath, err)
	}
	return nil
}

// LintCheckerPaths validates files and recursively discovers .yaml files in
// directories. It returns the number of discovered files and all traversal and
// validation errors joined together. An empty directory is valid.
func LintCheckerPaths(paths []string) (int, error) {
	files, lintErrors := collectCheckerFiles(paths)
	for _, filePath := range files {
		if err := LintCheckerFile(filePath); err != nil {
			lintErrors = append(lintErrors, err)
		}
	}
	return len(files), errors.Join(lintErrors...)
}

// collectCheckerFiles expands files and directories in argument order. A file
// supplied directly is validated regardless of its extension; directory
// traversal selects only .yaml files. Duplicate paths are validated once.
func collectCheckerFiles(paths []string) ([]string, []error) {
	var files []string
	var lintErrors []error
	seen := make(map[string]bool)
	addFile := func(filePath string) {
		filePath = filepath.Clean(filePath)
		if seen[filePath] {
			return
		}
		seen[filePath] = true
		files = append(files, filePath)
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			lintErrors = append(lintErrors, fmt.Errorf("access checker path %s failed: %w", path, err))
			continue
		}
		if !info.IsDir() {
			addFile(path)
			continue
		}

		err = filepath.WalkDir(path, func(filePath string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
				addFile(filePath)
			}
			return nil
		})
		if err != nil {
			lintErrors = append(lintErrors, fmt.Errorf("walk checker config directory %s failed: %w", path, err))
		}
	}
	return files, lintErrors
}
