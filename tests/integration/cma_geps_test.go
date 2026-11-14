//go:build integration

package integration

import (
	"fmt"
	"testing"
)

func TestCMA_GEPS(t *testing.T) {
	checkDateTime := dateHour(yesterday(1), "12")
	expectedPath := fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GEPS/Prod-grib/%s/grib2.0p25/gef.gra.001.%s024.grb2", checkDateTime, checkDateTime)

	tests := []TestCase{
		{
			Name:          "grib2/orig archive",
			DataType:      "cma_geps/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "24h",
			Member:        "001",
			ExpectedPath:  expectedPath,
		},
		{
			Name:          "grib2/orig archive with config",
			DataType:      "cma_geps/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "24h",
			Member:        "001",
			ConfigDir:     "local",
			ExpectedPath:  expectedPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runTestCase(t, tt)
		})
	}
}
