//go:build integration

package integration

import (
	"fmt"
	"testing"
)

func TestCMA_GFS_GMF(t *testing.T) {
	checkDate := yesterday(1)
	checkDateTime := dateHour(checkDate, "00")
	todayDate := today()
	todayDateTime := dateHour(todayDate, "00")

	checkDate2 := yesterday(2)
	checkDateTime2 := dateHour(checkDate2, "00")

	tests := []TestCase{
		// grib2/orig
		{
			Name:          "grib2/orig runtime",
			DataType:      "cma_gfs_gmf/current/grib2/orig",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/%s/data/output/grib2_orig/gmf.gra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig archive",
			DataType:      "cma_gfs_gmf/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/%s/ORIG/gmf.gra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig runtime with config",
			DataType:      "cma_gfs_gmf/current/grib2/orig",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/%s/data/output/grib2_orig/gmf.gra.%s003.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig archive with config",
			DataType:      "cma_gfs_gmf/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/%s/ORIG/gmf.gra.%s003.grb2", checkDateTime, checkDateTime),
		},

		// grib2/modelvar
		{
			Name:          "grib2/modelvar runtime",
			DataType:      "cma_gfs_gmf/current/grib2/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime2,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/%s/data/output/grib2_orig/modelvar%s003.grb2", checkDateTime2, checkDateTime2),
		},
		{
			Name:          "grib2/modelvar archive",
			DataType:      "cma_gfs_gmf/current/grib2/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime2,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/%s/MODELVAR/modelvar%s003.grb2", checkDateTime2, checkDateTime2),
		},
		{
			Name:          "grib2/modelvar runtime with config",
			DataType:      "cma_gfs_gmf/current/grib2/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime2,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/%s/data/output/grib2_orig/modelvar%s003.grb2", checkDateTime2, checkDateTime2),
		},
		{
			Name:          "grib2/modelvar archive with config",
			DataType:      "cma_gfs_gmf/current/grib2/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime2,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/%s/MODELVAR/modelvar%s003.grb2", checkDateTime2, checkDateTime2),
		},

		// bin/modelvar
		{
			Name:          "bin/modelvar runtime",
			DataType:      "cma_gfs_gmf/current/bin/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/modelvar%s_003", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/modelvar%s_003", todayDateTime)},
			},
		},
		{
			Name:          "bin/modelvar archive",
			DataType:      "cma_gfs_gmf/current/bin/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Fcst-long/%s/modelvar%s_003", checkDateTime, checkDateTime),
		},
		{
			Name:          "bin/modelvar runtime with config",
			DataType:      "cma_gfs_gmf/current/bin/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/modelvar%s_003", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/modelvar%s_003", todayDateTime)},
			},
		},
		{
			Name:          "bin/modelvar archive with config",
			DataType:      "cma_gfs_gmf/current/bin/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Fcst-long/%s/modelvar%s_003", checkDateTime, checkDateTime),
		},

		// bin/postvar
		{
			Name:          "bin/postvar runtime",
			DataType:      "cma_gfs_gmf/current/bin/postvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/postvar%s_003", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/postvar%s_003", todayDateTime)},
			},
		},
		{
			Name:          "bin/postvar runtime with config",
			DataType:      "cma_gfs_gmf/current/bin/postvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/postvar%s_003", checkDateTime),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/00/output/postvar%s_003", todayDateTime)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runTestCase(t, tt)
		})
	}
}
