// this is a auto generated file.
package config

var EmbeddedConfigs = [][2]string{
	{`hpc/gda_grapes_gfs/bin/modelvar`, `# gda grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GDAS/FCST_results
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
`},
	{`hpc/gda_grapes_gfs/bin/postvar`, `# gda grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GDAS/FCST_results
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`hpc/gda_grapes_gfs/grib2/modelvar`, `# gda grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_V2.2_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR`},
	{`hpc/gda_grapes_gfs/grib2/orig`, `# gda grapes gfs
#   orig grib2

default: NOTFOUND

file_name: gda.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_V2.2_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
`},
	{`hpc/gmf_grapes_gfs/bin/modelvar`, `# gmf grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GMFS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GMFS/FCST_results
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`hpc/gmf_grapes_gfs/bin/postvar`, `# gmf grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GMFS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GMFS/FCST_results
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`hpc/gmf_grapes_gfs/grib2/modelvar`, `# gmf grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  # run time dir
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  # archive dir
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR

  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR
`},
	{`hpc/gmf_grapes_gfs/grib2/ne`, `# gmf grapes gfs
#   grib2 ne

default: NOTFOUND

file_name: ne_gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_ne/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/CMACAST
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/CMACAST
`},
	{`hpc/gmf_grapes_gfs/grib2/orig`, `# gmf grapes gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
  - type: storage
    path: /sstorage1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG
`},
	{`local/gda_grapes_gfs/bin/modelvar`, `# gda grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GDAS/FCST_results
  - type: local
    path: Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/`},
	{`local/gda_grapes_gfs/bin/postvar`, `# gda grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GDAS/FCST_results
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`local/gda_grapes_gfs/grib2/modelvar`, `# gda grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR`},
	{`local/gda_grapes_gfs/grib2/orig`, `# gda grapes gfs
#   orig grib2

default: NOTFOUND

file_name: gda.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GDA_GRAPES_GFS_POST/gfs_togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG`},
	{`local/gmf_grapes_gfs/bin/modelvar`, `# gmf grapes gfs
#   modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GMFS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GMFS/FCST_results
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`local/gmf_grapes_gfs/bin/postvar`, `# gmf grapes gfs
#   postvar

default: NOTFOUND

file_name: postvar{.Year}{.Month}{.Day}{.Hour}_{.Forecast}

paths:
  - type: local
    path: /g2/nwp/GRAPES_GFS/MODEL/data/NWP_GMFS/{.Hour4DV}/output
  - type: local
    path: /g2/nwp/GRAPES_GFS/DATA/DATABAK/NWP_GMFS/FCST_results
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Fcst-long/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}`},
	{`local/gmf_grapes_gfs/grib2/modelvar`, `# gmf grapes gfs
#   grib2 modelvar

default: NOTFOUND

file_name: modelvar{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  # run time dir
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}

  # archive dir
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/MODELVAR`},
	{`local/gmf_grapes_gfs/grib2/ne`, `# gmf grapes gfs
#   grib2 ne

default: NOTFOUND

file_name: ne_gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_ne/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/CMACAST`},
	{`local/gmf_grapes_gfs/grib2/orig`, `# gmf grapes gfs
#   grib2 orig

default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.Forecast}.grb2

paths:
  - type: local
    path: /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/{.Year}{.Month}{.Day}{.Hour}
  - type: local
    path: /g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}/ORIG`},
}
