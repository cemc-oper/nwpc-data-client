# 查找本地数据文件

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
| `--member` | 集合成员编号，如 `000`、`014`，代入模板变量 `{{.Member}}` |
| `--show-types` | 列出可用数据类型并退出（配合 `--data-config-dir` 时列目录中的配置） |
| `--debug` | 输出调试日志，可看到每个候选路径的检查过程 |

## 搜索顺序

1. 按配置中 `paths` 列表的顺序遍历目录，`--location-level` 不匹配的
   条目跳过；
2. 每个目录依次拼接候选文件名（`file_names` 列表在前，`file_name`
   追加在最后）；
3. 返回第一个存在的文件；全部未命中时返回 `default`。

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
