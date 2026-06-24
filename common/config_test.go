package common

import (
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
