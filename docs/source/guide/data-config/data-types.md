# 数据类型

数据类型是**层级相对路径**，而不是扁平名称。它唯一确定一个数据配置
文件，是 `local` 与 `check local` 命令定位数据的核心参数。

## 组织方式

仓库中的配置按 `<系统>/<版本>/<格式>/<名称>.yaml` 组织在 `conf/local/`
目录下：

```
conf/local/cma_gfs_gmf/current/grib2/orig.yaml
```

对应的数据类型为去掉 `.yaml` 后缀的相对路径：

```
cma_gfs_gmf/current/grib2/orig
```

四层目录的含义：

| 层级 | 示例 | 说明 |
|---|---|---|
| 系统 | `cma_gfs_gmf`、`cma_meso_3km` | 业务模式系统及区域 |
| 版本 | `current` | 系统版本，`current` 表示现行版本 |
| 格式 | `grib2`、`bin` | 数据文件格式 |
| 名称 | `orig`、`modelvar` | 产品名称 |

## 解析规则

数据类型如何映射到配置文件，取决于配置来源：

- 仅 `--data-type`：加载名为 `local/<data-type>` 的**内嵌**配置；
- `--data-type` + `--data-config-dir=<dir>`：加载 `<dir>/<data-type>.yaml`
  （注意不再添加 `local/` 前缀）；
- `--data-config-file`：直接指定文件，不使用数据类型。

数据类型中包含 `..` 以逃出配置目录的写法会被拒绝。

## 支持 CMADAAS 挂载归档的数据类型

以下现有 data type 同时包含原 HPC 路径和 CMADAAS 挂载路径，不需要使用新的
数据类型名称。CMADAAS 挂载根由 `--storage-base` 或 checker YAML 的 `storage_base`
显式指定，不固定为某台机器的目录：

| 系统 | data type |
|---|---|
| CMA-GFS-GMF | `cma_gfs_gmf/current/grib2/orig` |
| CMA-GEPS | `cma_geps/current/grib2/orig` |
| CMA-REPS | `cma_reps/current/grib2/orig` |
| CMA-TYM | `cma_tym/current/grib2/orig` |
| CMA-MESO 3KM | `cma_meso_3km/current/grib2/orig` |
| CMA-MESO 1KM | `cma_meso_1km/current/grib2/orig` |

支持范围仅限上述 `current/grib2/orig` 配置，不包括 `orig.cold`、其他产品或部署仓库
中的版本化业务配置。CMADAAS 文件与 HPC 文件的命名不同，但共用同一数据身份，按配置
中的候选顺序查找。使用示例和挂载前提见[local 指南](../local/index.md)，落位方法见
[fetch 指南](../fetch/index.md)。

## 查看可用数据类型

```bash
# 内嵌配置
nwpc_data_client local --show-types

# 自定义配置目录
nwpc_data_client local --data-config-dir=./config --show-types
```

内置数据类型完整清单见[参考：内置数据类型清单](../../reference/embedded-data-types.md)。
