//go:build integration

package integration

import (
	"fmt"
	"testing"
)

func TestCMA_GFS_GDA(t *testing.T) {
	checkDate := yesterday(1)
	checkDateTime := dateHour(checkDate, "00")
	todayDate := today()
	todayDateTime := dateHour(todayDate, "00")

	checkDate2 := yesterday(2)
	checkDateTime2 := dateHour(checkDate2, "00")

	hour4DV := "21"
	checkDateTime4DV := dateHour(checkDate2, hour4DV)
	todayDateTime4DV := dateHour(todayDate, hour4DV)

	tests := []TestCase{
		// grib2/orig
		{
			Name:          "grib2/orig runtime",
			DataType:      "cma_gfs_gda/current/grib2/orig",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "0h",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/%s/data/output/grib2_orig/gda.gra.%s000.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig archive",
			DataType:      "cma_gfs_gda/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "0h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/%s/ORIG/gda.gra.%s000.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig runtime with config",
			DataType:      "cma_gfs_gda/current/grib2/orig",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "0h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/%s/data/output/grib2_orig/gda.gra.%s000.grb2", checkDateTime, checkDateTime),
		},
		{
			Name:          "grib2/orig archive with config",
			DataType:      "cma_gfs_gda/current/grib2/orig",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "0h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/%s/ORIG/gda.gra.%s000.grb2", checkDateTime, checkDateTime),
		},

		// grib2/modelvar
		{
			Name:          "grib2/modelvar runtime",
			DataType:      "cma_gfs_gda/current/grib2/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime2,
			ForecastTime:  "0h",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/%s/data/output/grib2_orig/modelvar%s000.grb2", checkDateTime2, checkDateTime2),
		},
		{
			Name:          "grib2/modelvar archive",
			DataType:      "cma_gfs_gda/current/grib2/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime2,
			ForecastTime:  "0h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/%s/MODELVAR/modelvar%s000.grb2", checkDateTime2, checkDateTime2),
		},
		{
			Name:          "grib2/modelvar runtime with config",
			DataType:      "cma_gfs_gda/current/grib2/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime2,
			ForecastTime:  "0h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/%s/data/output/grib2_orig/modelvar%s000.grb2", checkDateTime2, checkDateTime2),
		},
		{
			Name:          "grib2/modelvar archive with config",
			DataType:      "cma_gfs_gda/current/grib2/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime2,
			ForecastTime:  "0h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/%s/MODELVAR/modelvar%s000.grb2", checkDateTime2, checkDateTime2),
		},

		// bin/modelvar (4DVAR, check_date_4dvar=-2 day, hour=21)
		{
			Name:          "bin/modelvar runtime",
			DataType:      "cma_gfs_gda/current/bin/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/modelvar%s_006", checkDateTime4DV),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/modelvar%s_006", todayDateTime4DV)},
			},
		},
		{
			Name:          "bin/modelvar archive",
			DataType:      "cma_gfs_gda/current/bin/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/%s00/modelvar%s_006", checkDate, checkDateTime4DV),
		},
		{
			Name:          "bin/modelvar runtime with config",
			DataType:      "cma_gfs_gda/current/bin/modelvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/modelvar%s_006", checkDateTime4DV),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/modelvar%s_006", todayDateTime4DV)},
			},
		},
		{
			Name:          "bin/modelvar archive with config",
			DataType:      "cma_gfs_gda/current/bin/modelvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/%s00/modelvar%s_006", checkDate, checkDateTime4DV),
		},

		// bin/postvar
		{
			Name:          "bin/postvar runtime",
			DataType:      "cma_gfs_gda/current/bin/postvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/postvar%s_006", checkDateTime4DV),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/postvar%s_006", todayDateTime4DV)},
			},
		},
		{
			Name:          "bin/postvar archive",
			DataType:      "cma_gfs_gda/current/bin/postvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/%s00/postvar%s_006", checkDate, checkDateTime4DV),
		},
		{
			Name:          "bin/postvar runtime with config",
			DataType:      "cma_gfs_gda/current/bin/postvar",
			LocationLevel: "runtime",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/postvar%s_006", checkDateTime4DV),
			Fallbacks: []Fallback{
				{StartTime: todayDateTime, Path: fmt.Sprintf("/g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/00/output/postvar%s_006", todayDateTime4DV)},
			},
		},
		{
			Name:          "bin/postvar archive with config",
			DataType:      "cma_gfs_gda/current/bin/postvar",
			LocationLevel: "archive",
			StartTime:     checkDateTime,
			ForecastTime:  "3h",
			ConfigDir:     "local",
			ExpectedPath:  fmt.Sprintf("/g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/%s00/postvar%s_006", checkDate, checkDateTime4DV),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runTestCase(t, tt)
		})
	}
}
