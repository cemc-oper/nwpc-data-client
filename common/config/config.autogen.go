// this is a auto generated file.
package config

var EmbeddedConfigs = [][2]string{
	{`hpc/grapes_gfs_gda/bin/modelvar`, `# gmf grapes gfs
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
	{`hpc/grapes_gfs_gda/bin/postvar`, `# gmf grapes gfs
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
	{`hpc/grapes_gfs_gda/grib2/modelvar`, `# gmf grapes gfs
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
	{`hpc/grapes_gfs_gda/grib2/ne`, `# gmf grapes gfs
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
	{`hpc/grapes_gfs_gda/grib2/orig`, `# gmf grapes gfs
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
	{`hpc/grapes_gfs_gmf/bin/modelvar`, `# gda grapes gfs
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
	{`hpc/grapes_gfs_gmf/bin/postvar`, `# gda grapes gfs
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
	{`hpc/grapes_gfs_gmf/grib2/modelvar`, `# gda grapes gfs
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
	{`hpc/grapes_gfs_gmf/grib2/orig`, `# gda grapes gfs
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
	{`local/grapes_gfs_gda/bin/modelvar`, `# gda grapes gfs
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
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`local/grapes_gfs_gda/bin/postvar`, `# gda grapes gfs
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
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`local/grapes_gfs_gda/grib2/modelvar`, `# gda grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR`},
	{`local/grapes_gfs_gda/grib2/orig`, `# gda grapes gfs
#   orig grib2

default: NOTFOUND

file_name: gda.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG`},
	{`local/grapes_gfs_gmf/bin/modelvar`, `# gmf grapes gfs
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
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`local/grapes_gfs_gmf/bin/postvar`, `# gmf grapes gfs
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
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`local/grapes_gfs_gmf/grib2/modelvar`, `# gmf grapes gfs
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
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR`},
	{`local/grapes_gfs_gmf/grib2/ne`, `# gmf grapes gfs
#   grib2 ne

default: NOTFOUND

file_name: ne_gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_ne/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/CMACAST`},
	{`local/grapes_gfs_gmf/grib2/orig`, `# gmf grapes gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  - type: local
    level: archive
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG`},
}
