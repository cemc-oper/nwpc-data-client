# fetch 命令：将数据落位到目标目录

`nwpc_data_client fetch` 使用与 `local` 相同的数据配置和候选顺序，选中第一个
已经存在的普通文件，再将本地或 cmafs 二级存储文件复制到用户指定的目录；本地源
还可以使用符号链接。成功时 stdout 只输出最终绝对路径，便于直接传给后续程序。

它是一次性操作，**不会等待尚未生成的文件**；需要轮询等待本地文件时使用
[`check local`](../check-local/index.md)。

```{eval-rst}
.. toctree::
   :maxdepth: 1

   cmafs
   reliability
```

## 基本用法

```bash
nwpc_data_client fetch \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h \
    --target-dir=/work/input
```

成功后输出例如：

```text
/work/input/gmf.gra.2026062400024.grb2
```

目标文件名默认继承实际源文件的 basename，目标目录不存在时会自动创建。可以将结果
保存为 Shell 变量：

```bash
input_file="$(
    nwpc_data_client fetch \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --start-time=2026062400 \
        --forecast-time=24h \
        --target-dir=/work/input
)"
```

日志和错误写 stderr，不会混入上述 stdout 路径。不同存储来源的文件名可能不同，
后续程序应使用返回路径，不要自行拼接或假定文件名。

## 自定义目标文件名

需要固定或按起报时次命名时，使用 `--target-name`，不必在 fetch 之后再改名：

```bash
nwpc_data_client fetch \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h \
    --target-dir=/work/input \
    --target-name='gfs_{{.Year}}{{.Month}}{{.Day}}{{.Hour}}_{{.ForecastHour}}.grb2'
```

成功后输出 `/work/input/gfs_2026062400_024.grb2`。

参数值是 Go 模板，可用变量和函数与数据配置相同（`{{.Year}}`、`{{.ForecastHour}}`、
`{{.Member}}`、`generateStartTime` 等，见[模板变量与函数](../../reference/template-variables.md)），
另提供 `{{.SourceName}}` 表示实际选中源的 basename，例如 `pre_{{.SourceName}}`。
不含 `{{` 的纯字符串直接作为文件名使用。

渲染结果必须是**单个文件名**，不允许包含路径分隔符、`.`、`..` 或空串，因此文件始终
落在 `--target-dir` 内；非法名字会在创建目录和复制之前报错。目标复用、冲突检测与
无覆盖发布规则与默认命名完全一致，只是作用在自定义名字上。

## copy、symlink 与 auto

| `--method` | 行为 |
|---|---|
| `copy`（默认） | 复制为独立普通文件，保留基本权限和修改时间 |
| `symlink` | 创建指向本地源绝对路径的符号链接 |
| `auto` | runtime/archive 本地源使用 symlink；storage 本地源和 cmafs 使用 copy |

storage 使用 copy 是为了避免生成仅在特定存储节点可见的链接。显式
`--method=symlink` 时，调用方需要自行确认下游能够访问源路径，并接受源文件被
清理或修改的影响。cmafs 不支持显式 symlink。

## 选择数据来源

- 普通本地或直接挂载目录不需要额外 backend。
- CMADAAS 挂载目录按本地源处理，通过 `--storage-base` 指定挂载根；完整示例见
  [CMADAAS 挂载归档](../cmadaas.md)。
- 没有直接挂载二级存储时，可以通过机器 backend 使用 SSH/rsync 获取 cmafs 文件，
  见[从 cmafs 二级存储获取](cmafs.md)。

所有来源都要求源文件非空并通过稳定性检查。目标以临时文件和无覆盖发布方式写入，
重复执行可以复用匹配的目标，但不会自动覆盖冲突文件。详细规则见
[源文件检查与目标发布](reliability.md)。

## 参数与退出码

数据配置参数与 `local` 相同：`--data-type`、`--data-config-dir`、
`--data-config-file`、`--location-level`、`--storage-base`、`--start-time`、
`--forecast-time` 和 `--member`。fetch 特有参数：

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--target-dir` | 无，必选 | 最终文件所在目录 |
| `--target-name` | 空，取源 basename | 目标文件名，支持 Go 模板，只能是文件名 |
| `--backend-config` | 空 | 机器后端 YAML；只有实际需要 cmafs 时必需 |
| `--method` | `copy` | `copy`、`symlink` 或 `auto` |
| `--stable-interval` | `1s` | 两次源稳定性检查的间隔 |
| `--debug` | `false` | 输出调试日志到 stderr |

- 成功：退出码 0，stdout 输出最终绝对路径；
- 未找到、配置错误、源未稳定、目标冲突或复制失败：退出码 1；
- 收到 SIGHUP、SIGINT、SIGTERM 或 SIGQUIT：退出码 3，并清理本进程临时文件。

完整参数见[命令行参考](../../reference/cli/nwpc_data_client_fetch.md)，脚本中的输出
检查和常见错误处理见[脚本接入与故障排查](../troubleshooting.md)。
