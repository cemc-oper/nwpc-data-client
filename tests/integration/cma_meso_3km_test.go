//go:build integration

package integration

import (
	"fmt"
	"testing"
)

func TestCMA_MESO_3KM(t *testing.T) {
	checkDate := yesterday(1)
	checkDateTime := dateHour(checkDate, "00")
	todayDate := today()
	todayDateTime := dateHour(todayDate, "00")

	tests := []TestCase{
		// grib2/orig.cold
		{
			Name:          "grib2/orig.cold runtime",
			DataType:      "cma_meso_3km/current/grib2/orig.cold",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_MESO_3KM_POST_DATA/%s/data/output/grib2_orig/rmf.hgra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig.cold archive",
			DataType:      "cma_meso_3km/current/grib2/orig.cold",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Prod-grib/%s/ORIG/rmf.hgra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig.cold runtime with config",
			DataType:      "cma_meso_3km/current/grib2/orig.cold",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_MESO_3KM_POST_DATA/%s/data/output/grib2_orig/rmf.hgra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig.cold archive with config",
			DataType:      "cma_meso_3km/current/grib2/orig.cold",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Prod-grib/%s/ORIG/rmf.hgra.%s003.grb2", checkDateTime, checkDateTime),
		},

		// bin/modelvar.cold
		{
			Name:          "bin/modelvar.cold runtime",
			DataType:      "cma_meso_3km/current/bin/modelvar.cold",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/modelvar%s00300", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/modelvar%s00300", todayDateTime)},
			},
		},
		{
			Name:          "bin/modelvar.cold runtime/archive",
			DataType:      "cma_meso_3km/current/bin/modelvar.cold",
			LocationLevel: "runtime/archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/%s/modelvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/modelvar.cold archive",
			DataType:      "cma_meso_3km/current/bin/modelvar.cold",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/%s/modelvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/modelvar.cold runtime with config",
			DataType:      "cma_meso_3km/current/bin/modelvar.cold",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/modelvar%s00300", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/modelvar%s00300", todayDateTime)},
			},
		},
		{
			Name:          "bin/modelvar.cold runtime/archive with config",
			DataType:      "cma_meso_3km/current/bin/modelvar.cold",
			LocationLevel: "runtime/archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/%s/modelvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/modelvar.cold archive with config",
			DataType:      "cma_meso_3km/current/bin/modelvar.cold",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/%s/modelvar%s00300", checkDateTime, checkDateTime),
		},

		// bin/postvar.cold
		{
			Name:          "bin/postvar.cold runtime",
			DataType:      "cma_meso_3km/current/bin/postvar.cold",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/postvar%s00300", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/postvar%s00300", todayDateTime)},
			},
		},
		{
			Name:          "bin/postvar.cold runtime/archive",
			DataType:      "cma_meso_3km/current/bin/postvar.cold",
			LocationLevel: "runtime/archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/postvar.cold archive",
			DataType:      "cma_meso_3km/current/bin/postvar.cold",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/postvar.cold runtime with config",
			DataType:      "cma_meso_3km/current/bin/postvar.cold",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/postvar%s00300", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/00/fcst/grapes_model/run/postvar%s00300", todayDateTime)},
			},
		},
		{
			Name:          "bin/postvar.cold runtime/archive with config",
			DataType:      "cma_meso_3km/current/bin/postvar.cold",
			LocationLevel: "runtime/archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/postvar.cold archive with config",
			DataType:      "cma_meso_3km/current/bin/postvar.cold",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/%s/postvar%s00300", checkDateTime, checkDateTime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runTestCase(t, tt)
		})
	}
}
