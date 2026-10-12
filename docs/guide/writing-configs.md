# 编写数据配置

数据配置是一个 YAML 文件，**整体作为 Go 模板**先渲染再解析（当前格式
版本 v3 使用 `{{...}}` 定界符）。渲染时可用起报时次、预报时效、集合
成员等变量，完整清单见[模板变量与函数](../reference/template-variables.md)。

## 最小示例

```yaml
# nwpc-data-client data config
_config_version: 3

default: NOTFOUND
file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2
paths:
  - type: local
    level: runtime
    path: /g1/u/CMA_GFS_GMF/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}
```

起报时次 `2026062400`、预报时效 `24h` 渲染后，工具会检查：

```
/g1/u/CMA_GFS_GMF/2026062400/gmf.gra.2026062400024.grb2
```

## 各部分说明

**`_config_version`**（可选）：声明配置格式版本，当前版本为 3。
建议显式写出，避免依赖启发式探测。格式版本的演进与迁移见
[配置格式版本](config-versions.md)。

**`default`**：所有候选路径都未命中时的返回值，约定为 `NOTFOUND`。

**`file_name` / `file_names`**：候选文件名。两者可同时使用，
`file_names` 列表在前、`file_name` 追加在最后，按顺序尝试。
`file_names` 适合同一产品存在新旧两种命名的情况：

```yaml
file_names:
  - modelvar{{.Year}}{{.Month}}{{.Day}}{{.Hour}}_{{.ForecastHour}}
  - modelvar{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}
file_name: modelvar.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}.grb2
```

**`paths`**：候选目录列表，按顺序搜索。每个条目包含：

- `type`：目前仅支持 `local`；
- `level`：位置级别，`runtime` / `archive` / `storage`（可 `/` 组合），
  与命令行 `--location-level` 配合控制搜索范围；
- `path`：目录路径模板，不含文件名。

典型的"先查业务运行目录、再查归档目录"写法：

```yaml
paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/POST_DATA/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/output
  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}
```

## 常用模板技巧

起报时间偏移（如 4DVar 产品以上一时次为目录）：

```yaml
path: /data/4dvar/{{generateStartTime .StartTime -3 | getYear}}{{generateStartTime .StartTime -3 | getMonth}}{{generateStartTime .StartTime -3 | getDay}}
```

预报时效偏移与非默认位数：

```yaml
file_name: prod_{{generateForecastTime .ForecastTime "3h" | getForecastHour | printf "%03d"}}.grb2
```

分钟级时效：

```yaml
file_name: modelvar{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}{{.ForecastMinute}}
```

集合成员：

```yaml
file_name: gef.gra.{{.Member}}.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2
```

模板中也可以声明局部变量，把复杂表达式写一次复用：

```yaml
{{$startTime := generateStartTime .StartTime -3}}
file_name: modelvar{{getYear $startTime}}{{getMonth $startTime}}{{getDay $startTime}}{{getHour $startTime}}_{{.ForecastHour}}
```

## 校验

写完配置后用 `config lint` 校验（模板解析 + 结构校验）：

```bash
nwpc_data_client config lint --data-config-file=./my_config.yaml
```

字段的完整约束见[参考：数据配置格式](../reference/data-config-format.md)。
