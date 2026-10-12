# 内置数据类型清单

内嵌配置即 `conf/local/` 目录树，数据类型为去掉 `.yaml` 后缀的相对路径。
运行 `nwpc_data_client local --show-types` 可列出当前二进制实际内置的
数据类型。

命名约定：

- `_ctl` 后缀：控制预报产品；
- `.cold` 后缀：冷启动起报的产品（检索路径中的 warm/cold 目录不同；
  个别配置的文件名也略有差异。`cma_meso_1km/current/grib2/orig.cold`
  例外，与 warm 配置内容相同，仅保留数据类型名）；
- 集合预报（`cma_geps`、`cma_reps`）需配合 `--member` 使用。

下表"文件名模式"为渲染后的样式（`FF`/`FFF` 为预报时效，`MMM` 为集合
成员），准确模板以 `conf/local/` 下对应配置文件为准。

## CMA-GFS 全球模式（gmf 区域）

| 数据类型 | 文件名模式 |
|---|---|
| `cma_gfs_gmf/current/grib2/orig` | `gmf.gra.YYYYMMDDHHFFF.grb2` |
| `cma_gfs_gmf/current/grib2/ne` | `ne_gmf.gra.YYYYMMDDHHFFF.grb2` |
| `cma_gfs_gmf/current/grib2/modelvar` | `modelvarYYYYMMDDHHFFF.grb2` |
| `cma_gfs_gmf/current/bin/modelvar` | `modelvarYYYYMMDDHH_FFF` |
| `cma_gfs_gmf/current/bin/postvar` | `postvarYYYYMMDDHH_FFF` |

## CMA-GFS 全球模式（gda 区域）

| 数据类型 | 文件名模式 |
|---|---|
| `cma_gfs_gda/current/grib2/orig` | `gda.gra.YYYYMMDDHHFFF.grb2` |
| `cma_gfs_gda/current/grib2/modelvar` | `modelvarYYYYMMDDHHFFF.grb2` |
| `cma_gfs_gda/current/bin/modelvar` | `modelvarYYYYMMDDHH_FFF` |
| `cma_gfs_gda/current/bin/postvar` | `postvarYYYYMMDDHH_FFF` |

## CMA-MESO 1km

| 数据类型 | 文件名模式 |
|---|---|
| `cma_meso_1km/current/grib2/orig` | `rmf.hgra.YYYYMMDDHHFFF.grb2` |
| `cma_meso_1km/current/grib2/orig.cold` | `rmf.hgra.YYYYMMDDHHFFF.grb2` |
| `cma_meso_1km/current/bin/modelvar` | `modelvarYYYYMMDDHHFFFMM` |
| `cma_meso_1km/current/bin/modelvar.cold` | `modelvarYYYYMMDDHHFFFMM` |
| `cma_meso_1km/current/bin/modelvar_ctl` | `modelvar.ctl_YYYYMMDDHH00000` |
| `cma_meso_1km/current/bin/modelvar_ctl.cold` | `modelvar.ctl_YYYYMMDDHH00000` |

## CMA-MESO 3km

| 数据类型 | 文件名模式 |
|---|---|
| `cma_meso_3km/current/grib2/orig` | `rmf.hgra.YYYYMMDDHHFFF.grb2` |
| `cma_meso_3km/current/grib2/orig.cold` | `rmf.hgra.YYYYMMDDHHFFF.grb2` |
| `cma_meso_3km/current/bin/modelvar` | `modelvarYYYYMMDDHHFFFMM` |
| `cma_meso_3km/current/bin/modelvar.cold` | `modelvarYYYYMMDDHHFFFMM` |
| `cma_meso_3km/current/bin/modelvar_ctl` | `modelvar.ctl_YYYYMMDDHH00000` |
| `cma_meso_3km/current/bin/modelvar_ctl.cold` | `modelvar.ctl_YYYYMMDDHH00000` |
| `cma_meso_3km/current/bin/postvar` | `postvarYYYYMMDDHHFFFMM` |
| `cma_meso_3km/current/bin/postvar.cold` | `postvarYYYYMMDDHHFFFMM` |
| `cma_meso_3km/current/bin/postvar_ctl` | `post.ctl_YYYYMMDDHHFFF00` |
| `cma_meso_3km/current/bin/postvar_ctl.cold` | `postvar.ctl_YYYYMMDDHH00000` |

## 集合预报与台风模式

| 数据类型 | 文件名模式 |
|---|---|
| `cma_geps/current/grib2/orig` | `gef.gra.MMM.YYYYMMDDHHFFF.grb2` |
| `cma_reps/current/grib2/orig` | `mef.gra.MMM.YYYYMMDDHHFFF.grb2` |
| `cma_tym/current/grib2/orig` | `rmf.tcgra.YYYYMMDDHHFFF.grb2` |
| `cma_tym/current/bin/modelvar` | `modelvarYYYYMMDDHHFFFMM` |
| `cma_tym/current/bin/modelvar_ctl` | `model.ctl_YYYYMMDDHH00000` |
| `cma_tym/current/bin/postvar` | `postvarYYYYMMDDHHFFFMM` |
| `cma_tym/current/bin/postvar_ctl` | `post.ctl_YYYYMMDDHH00000` |

> 新增或修改配置后请运行 `make generate` 再生成内嵌代码，并同步更新本页。
