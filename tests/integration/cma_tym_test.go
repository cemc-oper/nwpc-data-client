//go:build integration

package integration

import (
	"fmt"
	"testing"
)

func TestCMA_TYM(t *testing.T) {
	checkDate := yesterday(1)
	checkDateTime := dateHour(checkDate, "00")

	tests := []TestCase{
		// grib2/orig
		{
			Name:          "grib2/orig runtime",
			DataType:      "cma_tym/current/grib2/orig",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_TYM_POST_DATA/%s/data/output/grib2_orig/rmf.tcgra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig archive",
			DataType:      "cma_tym/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/TYM/Prod-grib/%s/ORIG/rmf.tcgra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig runtime with config",
			DataType:      "cma_tym/current/grib2/orig",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_TYM_POST_DATA/%s/data/output/grib2_orig/rmf.tcgra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig archive with config",
			DataType:      "cma_tym/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/TYM/Prod-grib/%s/ORIG/rmf.tcgra.%s003.grb2", checkDateTime, checkDateTime),
		},

		// bin/modelvar
		{
			Name:          "bin/modelvar runtime",
			DataType:      "cma_tym/current/bin/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/%s/modelvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/modelvar runtime with config",
			DataType:      "cma_tym/current/bin/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/%s/modelvar%s00300", checkDateTime, checkDateTime),
		},

		// bin/postvar
		{
			Name:          "bin/postvar runtime",
			DataType:      "cma_tym/current/bin/postvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/postvar archive",
			DataType:      "cma_tym/current/bin/postvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/TYM/Fcst-main/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/postvar runtime with config",
			DataType:      "cma_tym/current/bin/postvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/postvar archive with config",
			DataType:      "cma_tym/current/bin/postvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/TYM/Fcst-main/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runTestCase(t, tt)
		})
	}
}
