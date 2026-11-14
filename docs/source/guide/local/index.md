# local 命令：查找本地数据文件

`nwpc_data_client local` 按配置模板计算候选文件路径，依次检查并输出
第一个实际存在的文件路径。

## 基本用法

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h
```

- 找到文件：输出完整路径，退出码 0；
- 未找到：输出配置中的 `default` 值（通常为 `NOTFOUND`）。

## 参数

必选（除非使用 `--data-config-file` 或 `--show-types`）：

| 参数 | 说明 |
|---|---|
| `--data-type` | 数据类型，层级相对路径，如 `cma_gfs_gmf/current/grib2/orig` |
| `--start-time` | 起报时次，`YYYYMMDDHH`，如 `2026062400` |
| `--forecast-time` | 预报时效，Go duration 格式，如 `0h`、`24h`、`1h30m` |

可选：

| 参数 | 说明 |
|---|---|
| `--data-config-dir` | 自定义配置目录，与 `--data-type` 联用，替代内嵌配置 |
| `--data-config-file` | 直接指定单个配置文件，优先级最高 |
| `--location-level` | 逗号分隔的位置级别（如 `runtime,archive`），只搜索对应级别的路径条目；空为全部 |
| `--storage-base` | `{{.StorageBase}}` 使用的本地存储挂载根；非空时必须是绝对路径，配置引用该变量时必填 |
| `--member` | 集合成员编号，如 `000`、`014`，代入模板变量 `{{.Member}}` |
| `--show-types` | 列出可用数据类型并退出（配合 `--data-config-dir` 时列目录中的配置） |
| `--debug` | 输出调试日志，可看到每个候选路径的检查过程 |

## 搜索顺序

1. 按配置中 `paths` 列表的顺序遍历目录，`--location-level` 不匹配的
   条目跳过；
2. 每个目录依次拼接候选文件名（`file_names` 列表在前，`file_name`
   追加在最后）；
3. 返回第一个存在的文件；全部未命中时返回 `default`。

## CMADAAS 挂载目录

六个 `current/grib2/orig` 内嵌配置支持直接读取 CMADAAS 挂载目录；具体 data type
见[数据类型指南](../data-config/data-types.md)。挂载目录在配置中是
`type: local, level: archive`，不使用 CMADaaS API，也不需要 cmafs backend。

**挂载根不写死在配置中。** CMADAAS 路径以 `{{.StorageBase}}/DATA/...` 或
`{{.StorageBase}}/ORIG-DATA/...` 开头，调用方通过 `--storage-base` 指定当前机器的
实际挂载根。例如当前 WSL 使用 `/CMADAAS`：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026111300 --forecast-time=24h \
    --storage-base=/CMADAAS --location-level=archive
```

其他机器可传入自己的绝对挂载路径，例如 `--storage-base=/mnt/cmadaas`。这六份配置
实际引用该变量，因此未提供 `--storage-base` 时会明确报错，不会尝试 `/DATA` 等错误
路径。`check local` 除同名 CLI 参数外，也可在 checker YAML 中设置：

```yaml
storage_base: /CMADAAS
```

`--location-level=archive` 只搜索归档 local 条目，跳过 storage 和 cmafs。原 HPC
archive 路径仍排在 CMADAAS 之前；两者都存在时返回原 HPC 文件。HPC 与 CMADAAS 的
文件名不同，调用方应直接使用 stdout 返回的实际路径。

GEPS、REPS 需要同时传入三位成员编号，例如 `--member=001`。REPS 配置包含起报当日
和次日两个目录候选，以兼容当前 06/18 时次的归档布局。

运行前必须确认传入的挂载根已挂载且当前用户可读。`local` 沿用现有 Stat 错误处理，
挂载不存在或不可读时可能最终输出 `NOTFOUND`；可使用 `--debug` 查看候选检查过程，
并直接检查挂载目录进行诊断。

引用 `{{.StorageBase}}` 的配置只能由支持该变量的 client 使用；旧 client 会在整文件
模板渲染时报错。部署时必须先升级 client，再发布这些配置。

## 示例

只搜索 runtime 级别的目录：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 --forecast-time=24h \
    --location-level=runtime
```

集合预报成员：

```bash
nwpc_data_client local \
    --data-type=cma_geps/current/grib2/orig \
    --start-time=2026062400 --forecast-time=24h \
    --member=014
```

使用自定义配置目录：

```bash
nwpc_data_client local \
    --data-config-dir=./config --data-type=cma_meso_1km/current/bin/modelvar \
    --start-time=2026062400 --forecast-time=24h
```

使用单个配置文件：

```bash
nwpc_data_client local \
    --data-config-file=./my_config.yaml \
    --start-time=2026062400 --forecast-time=3h
```
