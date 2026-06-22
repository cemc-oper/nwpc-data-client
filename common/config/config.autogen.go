// this is a auto generated file.
package config

var EmbeddedConfigs = [][2]string{
	{`local/cma_geps/current/grib2/orig`, `# cma-reps
#   grib2 orig

default: NOTFOUND

file_name: gef.gra.{.Member}.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GEPS/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/grib2`},
	{`local/cma_gfs_gda/current/bin/modelvar`, `# cma-gfs
#   modelvar

default: NOTFOUND

  { $gdaStartTime := generateStartTime .StartTime -3 }
  { $gdaForecastTime := generateForecastTime .ForecastTime "3h" }

file_name: "modelvar{ getYear $gdaStartTime }{getMonth $gdaStartTime }{ getDay $gdaStartTime }{ getHour $gdaStartTime}_{ getForecastHour $gdaForecastTime | printf "%03d"}"

paths:
  - type: local
    level: runtime
    path: /g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/{.Hour}/output

  - type: local
    level: archive
    path:  /g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_gfs_gda/current/bin/postvar`, `# CMA-GFS gda
#   postvar

default: NOTFOUND

  { $gdaStartTime := generateStartTime .StartTime -3 }
  { $gdaForecastTime := generateForecastTime .ForecastTime "3h" }

file_name: "postvar{ getYear $gdaStartTime}{ getMonth $gdaStartTime }{ getDay $gdaStartTime }{ getHour $gdaStartTime}_{ getForecastHour $gdaForecastTime | printf "%03d"}"

paths:
  - type: local
    level: runtime
    path: /g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_EN4DVAR/{.Hour}/output

  - type: local
    level: archive
    path:  /g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/{.Year}{.Month}{.Day}{.Hour}
`},
	{`local/cma_gfs_gda/current/grib2/modelvar`, `# CMA-GFS gda
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/MODELVAR`},
	{`local/cma_gfs_gda/current/grib2/orig`, `# CMA-GFS gda
#   grib2 orig

default: NOTFOUND

file_name: gda.gra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG
`},
	{`local/cma_gfs_gmf/current/bin/modelvar`, `# cma-gfs
#   modelvar
#   version >= v0.8.0

default: NOTFOUND

file_names:
  - modelvar{.Year}{.Month}{.Day}{.Hour}_{.ForecastHour}
  # - modelvar{generateStartTime .StartTime -3 | getYear}{generateStartTime .StartTime -3 | getMonth}{generateStartTime .StartTime -3 | getDay}{generateStartTime .StartTime -3 | getHour}_{generateForecastTime .ForecastTime "3h" | getForecastHour | printf "%03d"}

paths:
  - type: local
    level: runtime
    path: /g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/{.Hour}/output

  - type: local
    level: archive
    path:  /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Fcst-long/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_gfs_gmf/current/bin/postvar`, `# cma-gfs
#   postvar
#   version >= v0.8.0


default: NOTFOUND

file_names:
  - postvar{.Year}{.Month}{.Day}{.Hour}_{.ForecastHour}

paths:
  - type: local
    level: runtime
    path: /g2/op_gfs/CMA-GFS/CMA-GFS4.2.3_DATA/MODEL/GRAPES_GMFS/{.Hour}/output

#  - type: local
#    level: archive
#    path:  /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Fcst-long/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_gfs_gmf/current/grib2/modelvar`, `# cma-gfs gmf
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  # run time dir
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig/

  # archive dir
  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/MODELVAR`},
	{`local/cma_gfs_gmf/current/grib2/ne`, `# cma-gfs gmf
#   grib2 ne

default: NOTFOUND

file_name: ne_gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_ne/

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/CMACAST`},
	{`local/cma_gfs_gmf/current/grib2/orig`, `# cma-gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig/

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_meso_1km/current/bin/modelvar.cold`, `# CMA-MESO 1KM
#   modelvar cold

default: NOTFOUND

file_names:
  - modelvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_1km/cold/{.Hour}/fcst/grapes_model/run`},
	{`local/cma_meso_1km/current/bin/modelvar`, `# CMA-MESO 1KM
#   modelvar warm

default: NOTFOUND

file_names:
  - modelvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_1km/warm/{.Hour}/fcst/grapes_model/run
`},
	{`local/cma_meso_1km/current/bin/modelvar_ctl.cold`, `# CMA-MESO 1KM
#   modelvar_ctl cold

default: NOTFOUND

file_names:
  - modelvar.ctl_{.Year}{.Month}{.Day}{.Hour}00000

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_1km/cold/{.Hour}/fcst/grapes_model/run`},
	{`local/cma_meso_1km/current/bin/modelvar_ctl`, `# CMA-MESO 1KM
#   modelvar_ctl warm

default: NOTFOUND

file_names:
  - modelvar.ctl_{.Year}{.Month}{.Day}{.Hour}00000

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_1km/warm/{.Hour}/fcst/grapes_model/run
`},
	{`local/cma_meso_1km/current/grib2/orig.cold`, `# CMA-MESO 1KM
#   grib2_orig cold

default: NOTFOUND

file_name: rmf.hgra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_MESO_1KM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_1KM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_meso_1km/current/grib2/orig`, `# CMA-MESO 1KM
#   grib2_orig warm

default: NOTFOUND

file_name: rmf.hgra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_MESO_1KM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_1KM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_meso_3km/current/bin/modelvar.cold`, `# CMA-MESO 3KM
#   modelvar cold

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/modelvar`, `# CMA-MESO 3KM
#   modelvar warm

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/warm/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/warm/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/modelvar_ctl.cold`, `# CMA-MESO 3KM
#   modelvar_ctl cold

default: NOTFOUND

file_name: modelvar.ctl_{.Year}{.Month}{.Day}{.Hour}00000

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/modelvar_ctl`, `# CMA-MESO 3KM
#   modelvar_ctl warm

default: NOTFOUND

file_name: modelvar.ctl_{.Year}{.Month}{.Day}{.Hour}00000

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/warm/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/warm/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar.cold`, `# CMA-MESO 3KM
#   postvar cold

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar`, `# CMA-MESO 3KM
#   postvar warm

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/warm/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/warm/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar_ctl.cold`, `# CMA-MESO 3KM
#   postvar_ctl cold

default: NOTFOUND

file_name: postvar.ctl_{.Year}{.Month}{.Day}{.Hour}00000

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar_ctl`, `# CMA-MESO 3KM
#   postvar

default: NOTFOUND

file_name: post.ctl_{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/warm/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/warm/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/grib2/orig.cold`, `# CMA-MESO 3KM
#   grib2 orig cold

default: NOTFOUND

file_name: rmf.hgra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_MESO_3KM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_meso_3km/current/grib2/orig`, `# CMA-MESO 3KM
#   grib2 orig

default: NOTFOUND

file_name: rmf.hgra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_MESO_3KM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/MESO_3KM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_reps/current/grib2/orig`, `# cma-reps
#   grib2 orig

default: NOTFOUND

file_name: mef.gra.{.Member}.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/REPS/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/grib2`},
	{`local/cma_tym/current/bin/modelvar`, `# CMA-TYM
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  # - type: local
  #   level: archive
  #   path: /g3/COMMONDATA/OPER/NWPC/GRAPES_TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/bin/modelvar_ctl`, `# CMA-TYM
#   modelvar_ctl

default: NOTFOUND

file_name: model.ctl_{.Year}{.Month}{.Day}{.Hour}00000

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  # - type: local
  #   level: archive
  #   path: /g3/COMMONDATA/OPER/CEMC/TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/bin/postvar`, `# CMA-TYM
#   postvar

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}{.ForecastMinute}

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/bin/postvar_ctl`, `# CMA-TYM
#   postvar_ctl

default: NOTFOUND

file_name: post.ctl_{.Year}{.Month}{.Day}{.Hour}00000

paths:
  - type: local
    level: runtime
    path: /g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_V4_DATA/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/grib2/orig`, `# CMA-TYM
#   grib2 orig

default: NOTFOUND

file_name: rmf.tcgra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_TYM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/TYM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG
`},
}
