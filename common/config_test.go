package common

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigContent(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		startTime    time.Time
		forecastTime time.Duration
		member       string
		expected     DataConfig
	}{
		{
			name: "Test 1",
			content: `# cma-gfs
# grib2 orig
default: NOTFOUND
file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/data/output/grib2_orig/

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
`,
			startTime:    time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC),
			forecastTime: 24 * time.Hour,
			member:       "",
			expected: DataConfig{
				Default:  "NOTFOUND",
				FileName: "gmf.gra.2026061700024.grb2",
				Paths: []PathItem{
					{
						"local",
						"runtime",
						"/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2026061700/data/output/grib2_orig/",
					},
					{
						"local",
						"archive",
						"/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/2026061700/ORIG",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataConfig, err := ParseConfigContent(tt.content, tt.startTime, tt.forecastTime, tt.member)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, dataConfig)
		})
	}
}

const testLocalConfigContent = `default: NOTFOUND
file_name: data.grb2
paths:
  - type: local
    level: runtime
    path: /some/dir
`

func TestLoadConfigContentFromConfigDir(t *testing.T) {
	configDir := t.TempDir()
	dataType := "test_system/current/grib2/orig"
	configPath := filepath.Join(configDir, filepath.FromSlash(dataType+".yaml"))
	assert.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0755))
	assert.NoError(t, os.WriteFile(configPath, []byte(testLocalConfigContent), 0644))

	content, err := LoadConfigContent(configDir, dataType)

	assert.NoError(t, err)
	assert.Equal(t, testLocalConfigContent, content)
}

func TestLoadConfigContentRejectsPathEscape(t *testing.T) {
	configDir := t.TempDir()

	_, err := LoadConfigContent(configDir, "../escape")

	assert.Error(t, err)
}

func TestLoadConfigContentMissingFileDoesNotFallBackToEmbedded(t *testing.T) {
	// The data type exists in the embedded configs, but a non-empty config
	// dir is a full replacement: a missing file is an error, never a
	// fallback to the embedded configs.
	configDir := t.TempDir()
	dataType := "cma_gfs_gmf/current/grib2/orig"

	_, err := LoadConfigContent(configDir, dataType)

	assert.Error(t, err)
}

func TestLoadConfigContentEmbeddedHit(t *testing.T) {
	content, err := LoadConfigContent("", "local/cma_gfs_gmf/current/grib2/orig")

	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestLoadConfigContentEmbeddedMiss(t *testing.T) {
	_, err := LoadConfigContent("", "local/not_exist_system/current/grib2/orig")

	assert.Error(t, err)
}
