// this is a auto generated file.
package config

var EmbeddedConfigs = [][2]string{
	{`hpc/grapes_gfs_gda/bin/modelvar`, `# gda grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: local
    level: runtime
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/{.Hour4DV}/output

  - type: local
    level: runtime/archive
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GDAS/FCST_results

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
`},
	{`hpc/grapes_gfs_gda/bin/postvar`, `# gda grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: local
    level: runtime
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/{.Hour4DV}/output

  - type: local
    level: runtime/archive
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GDAS/FCST_results

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`hpc/grapes_gfs_gda/grib2/modelvar`, `# gda grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_V2.2_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR`},
	{`hpc/grapes_gfs_gda/grib2/orig`, `# gda grapes gfs
#   orig grib2

default: NOTFOUND

file_name: gda.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_V2.2_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
`},
	{`hpc/grapes_gfs_gmf/bin/modelvar`, `# gmf grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: local
    level: runtime
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GMFS/{.Hour4DV}/output

  - type: local
    level: runtime/archive
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GMFS/FCST_results

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`hpc/grapes_gfs_gmf/bin/postvar`, `# gmf grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: local
    level: runtime
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GMFS/{.Hour4DV}/output

  - type: local
    level: runtime/archive
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GMFS/FCST_results

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`hpc/grapes_gfs_gmf/grib2/modelvar`, `# gmf grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  # run time dir
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  # archive dir
  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR
`},
	{`hpc/grapes_gfs_gmf/grib2/ne`, `# gmf grapes gfs
#   grib2 ne

default: NOTFOUND

file_name: ne_gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_ne/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/CMACAST

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/CMACAST
`},
	{`hpc/grapes_gfs_gmf/grib2/orig`, `# gmf grapes gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG

  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
`},
	{`local/cma_gfs_gda/current/bin/modelvar`, `# CMA-GFS gda
#   modelvar

default: NOTFOUND

file_name: modelvar{generateStartTime .StartTime -3 | getYear}{generateStartTime .StartTime -3 | getMonth}{generateStartTime .StartTime -3 | getDay}{generateStartTime .StartTime -3 | getHour}_{.Forecast}

paths:
  - type: local
    level: runtime
    path: /g0/nwp/CMA-GFS4.0_DATA/MODEL/data/NWP_GDAS/{.Hour}/output

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_gfs_gda/current/bin/postvar`, `# CMA-GFS gda
#   postvar

default: NOTFOUND

file_name: postvar{generateStartTime .StartTime -3 | getYear}{generateStartTime .StartTime -3 | getMonth}{generateStartTime .StartTime -3 | getDay}{generateStartTime .StartTime -3 | getHour}_{.Forecast}

paths:
  - type: local
    level: runtime
    path: /g0/nwp/CMA-GFS4.0_DATA/MODEL/data/NWP_GDAS/{.Hour}/output

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year}{.Month}{.Day}{.Hour}
`},
	{`local/cma_gfs_gda/current/grib2/modelvar`, `# CMA-GFS gda
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g0/nwp_pd/NWP_CMA_GFS_GDA_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/MODELVAR`},
	{`local/cma_gfs_gda/current/grib2/orig`, `# CMA-GFS gda
#   grib2 orig

default: NOTFOUND

file_name: gda.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g0/nwp_pd/NWP_CMA_GFS_GDA_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG
`},
	{`local/cma_gfs_gmf/current/bin/modelvar`, `# cma-gfs gmf
#   grib2 modelvar

default: NOTFOUND

file_names:
  - modelvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}
  # - modelvar{generateStartTime .StartTime -3 | getYear}{generateStartTime .StartTime -3 | getMonth}{generateStartTime .StartTime -3 | getDay}{generateStartTime .StartTime -3 | getHour}_{generateForecastTime .ForecastTime "3h" | getForecastHour | printf "%03d"}

paths:
  - type: local
    level: runtime
    path: /g0/nwp/CMA-GFS4.0_DATA/MODEL/NWP_GMFS/{.Hour}/output

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_gfs_gmf/current/grib2/modelvar`, `# cma-gfs gmf
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  # run time dir
  - type: local
    level: runtime
    path: /g0/nwp_pd/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig/

  # archive dir
  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/MODELVAR`},
	{`local/cma_gfs_gmf/current/grib2/ne`, `# cma-gfs gmf
#   grib2 ne

default: NOTFOUND

file_name: ne_gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g0/nwp_pd/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_ne/

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/CMACAST`},
	{`local/cma_gfs_gmf/current/grib2/orig`, `# cma-gfs gmf
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g0/nwp_pd/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig/

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_meso_3km/current/bin/modelvar.cold`, `# CMA-MESO 3KM
#   modelvar cold

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/cold/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/modelvar`, `# CMA-MESO 3KM
#   modelvar warm

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/warm/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/DATABAK/warm/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar.cold`, `# CMA-MESO 3KM
#   postvar cold

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/cold/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar`, `# CMA-MESO 3KM
#   postvar warm

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/warm/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/DATABAK/warm/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar_ctl.cold`, `# CMA-MESO 3KM
#   postvar ctl cold

default: NOTFOUND

file_name: post.ctl_{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/cold/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/DATABAK/cold/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/bin/postvar_ctl`, `# CMA-MESO 3KM
#   postvar ctl

default: NOTFOUND

file_name: post.ctl_{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/warm/{.Hour}/fcst/grapes_model/run

  - type: local
    level: runtime/archive
    path: /g0/nwp/NWP_RMFS_DATA/cma_meso_3km/DATABAK/warm/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_meso_3km/current/grib2/orig.cold`, `# CMA-MESO 3KM
#   grib2 orig cold

default: NOTFOUND

file_name: rmf.hgra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g0/nwp_pd/NWP_GRAPES_MESO_3KM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}_cold/togrib2/output

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_meso_3km/current/grib2/orig`, `# CMA-MESO 3KM
#   grib2 orig

default: NOTFOUND

file_name: rmf.hgra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g0/nwp_pd/NWP_GRAPES_MESO_3KM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/togrib2/output

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`local/cma_tym/current/bin/modelvar`, `# CMA-TYM
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/run

  - type: local
    level: runtime_archive
    path: /g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/bin/modelvar_ctl`, `# CMA-TYM
#   modelvar

default: NOTFOUND

file_name: model.ctl_{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/run

  - type: local
    level: runtime_archive
    path: /g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/bin/postvar`, `# CMA-TYM
#   postvar

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}00

paths:
  - type: local
    level: runtime
    path: /g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/bin/postvar_ctl`, `# CMA-TYM
#   postvar ctl

default: NOTFOUND

file_name: post.ctl_{.Year}{.Month}{.Day}{.Hour}

paths:
  - type: local
    level: runtime
    path: /g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Fcst-main/{.Year}{.Month}{.Day}{.Hour}`},
	{`local/cma_tym/current/grib2/orig`, `# CMA-TYM
#   grib2 orig

default: NOTFOUND

file_name: rmf.tcgra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_GRAPES_TYM_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/rundir/output/orig_grib2

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG`},
	{`storage/grapes_gfs_gda/bin/modelvar`, `# gda grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
`},
	{`storage/grapes_gfs_gda/bin/postvar`, `# gda grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`storage/grapes_gfs_gda/grib2/modelvar`, `# gda grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR`},
	{`storage/grapes_gfs_gda/grib2/orig`, `# gda grapes gfs
#   orig grib2

default: NOTFOUND

file_name: gda.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
`},
	{`storage/grapes_gfs_gmf/bin/modelvar`, `# gmf grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`storage/grapes_gfs_gmf/bin/postvar`, `# gmf grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`storage/grapes_gfs_gmf/grib2/modelvar`, `# gmf grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR
`},
	{`storage/grapes_gfs_gmf/grib2/ne`, `# gmf grapes gfs
#   grib2 ne

default: NOTFOUND

file_name: ne_gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/CMACAST
`},
	{`storage/grapes_gfs_gmf/grib2/orig`, `# gmf grapes gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
`},
}
