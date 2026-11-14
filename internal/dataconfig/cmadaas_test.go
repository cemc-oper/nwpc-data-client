package dataconfig

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCMADAASRepsCrossYearPaths(t *testing.T) {
	content, err := LoadConfigContent("", "local/cma_reps/current/grib2/orig")
	require.NoError(t, err)
	startTime := time.Date(2026, 12, 31, 18, 0, 0, 0, time.UTC)
	storageBase := filepath.Join(t.TempDir(), "cmadaas")
	config, err := ParseConfigContentWithStorageBase(content, startTime, 24*time.Hour, "001", storageBase)
	require.NoError(t, err)

	assert.Equal(t, []string{"mef.gra.001.2026123118024.grb2"}, config.FileNames)
	assert.Equal(t,
		"Z_NAFP_C_BABJ_20261231180000_P_NWPC-GRAPES-REPS-CN-02400-m001.grib2",
		config.FileName)

	var cmadaasPaths []PathItem
	for _, item := range config.Paths {
		if strings.HasPrefix(item.Path, storageBase+"/") {
			cmadaasPaths = append(cmadaasPaths, item)
		}
	}
	require.Len(t, cmadaasPaths, 2)
	assert.Equal(t, PathItem{
		PathType:  "local",
		LevelType: "archive",
		Path:      storageBase + "/ORIG-DATA/NAFP/CMA-REPS/2026/20261231",
	}, cmadaasPaths[0])
	assert.Equal(t, PathItem{
		PathType:  "local",
		LevelType: "archive",
		Path:      storageBase + "/ORIG-DATA/NAFP/CMA-REPS/2027/20270101",
	}, cmadaasPaths[1])
}

func TestParseConfigContentAllowsEmptyStorageBase(t *testing.T) {
	content := `default: NOTFOUND
file_name: data.grb2
paths:
  - type: local
    level: runtime
    path: /hpc/data
  - type: local
    level: archive
    path: '{{.StorageBase}}/DATA'
`
	config, err := ParseConfigContent(content, time.Now(), 0, "")
	require.NoError(t, err)
	require.Len(t, config.Paths, 2)
	assert.Equal(t, "/hpc/data", config.Paths[0].Path)
	assert.Equal(t, "/DATA", config.Paths[1].Path)
}

func TestParseConfigContentRejectsRelativeStorageBase(t *testing.T) {
	_, err := ParseConfigContentWithStorageBase("default: NOTFOUND\npaths: []\n", time.Now(), 0, "", "relative/mount")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "absolute local path")
}
