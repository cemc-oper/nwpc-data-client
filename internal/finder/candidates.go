package finder

import (
	"path/filepath"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
)

// CandidatePaths expands a rendered data config into full candidate paths
// without checking whether they exist. Path order is paths first, then
// file_names, with file_name last, matching the historical local lookup.
func CandidatePaths(config dataconfig.DataConfig, locationLevels []string) []dataconfig.PathItem {
	fileNames := make([]string, 0, len(config.FileNames)+1)
	fileNames = append(fileNames, config.FileNames...)
	if len(config.FileName) > 0 {
		fileNames = append(fileNames, config.FileName)
	}

	var candidates []dataconfig.PathItem
	for _, item := range config.Paths {
		if !InLocationLevelTypes(locationLevels, item.LevelType) {
			continue
		}
		for _, fileName := range fileNames {
			candidates = append(candidates, dataconfig.PathItem{
				Path:      filepath.Join(item.Path, fileName),
				PathType:  item.PathType,
				LevelType: item.LevelType,
			})
		}
	}
	return candidates
}
