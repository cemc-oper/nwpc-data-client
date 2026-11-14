package dataconfig

import (
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
	config, err := ParseConfigContentWithStorageBase(content, startTime, 24*time.Hour, "001", "/mnt/cmadaas")
	require.NoError(t, err)

	assert.Equal(t, []string{"mef.gra.001.2026123118024.grb2"}, config.FileNames)
	assert.Equal(t,
		"Z_NAFP_C_BABJ_20261231180000_P_NWPC-GRAPES-REPS-CN-02400-m001.grib2",
		config.FileName)

	var cmadaasPaths []PathItem
	for _, item := range config.Paths {
		if strings.HasPrefix(item.Path, "/mnt/cmadaas/") {
			cmadaasPaths = append(cmadaasPaths, item)
		}
	}
	require.Len(t, cmadaasPaths, 2)
	assert.Equal(t, PathItem{
		PathType:  "local",
		LevelType: "archive",
		Path:      "/mnt/cmadaas/ORIG-DATA/NAFP/CMA-REPS/2026/20261231",
	}, cmadaasPaths[0])
	assert.Equal(t, PathItem{
		PathType:  "local",
		LevelType: "archive",
		Path:      "/mnt/cmadaas/ORIG-DATA/NAFP/CMA-REPS/2027/20270101",
	}, cmadaasPaths[1])
}

func TestParseConfigContentRequiresStorageBaseWhenUsed(t *testing.T) {
	content := `default: NOTFOUND
file_name: data.grb2
paths:
  - type: local
    level: archive
    path: '{{.StorageBase}}/DATA'
`
	_, err := ParseConfigContent(content, time.Now(), 0, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "storage base is required")
}

func TestParseConfigContentRejectsRelativeStorageBase(t *testing.T) {
	_, err := ParseConfigContentWithStorageBase("default: NOTFOUND\npaths: []\n", time.Now(), 0, "", "relative/mount")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "absolute local path")
}
