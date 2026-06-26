# NWPC Data Client

[![CI](https://github.com/cemc-oper/nwpc-data-client/actions/workflows/ci.yml/badge.svg)](https://github.com/cemc-oper/nwpc-data-client/actions/workflows/ci.yml)

面向 CMA 超算平台开发的 CEMC/CMA 数值天气预报业务系统数据查找命令行工具。

## 功能

- 在一组目录中查找数值预报业务系统数据文件路径
- 批量监视数据文件是否生成，当文件就绪后执行命令

## 安装

下载最新发布版本，或从源码构建。

### 本地构建 (Linux / macOS / WSL / Git Bash)

使用 `make` 将在 `bin/` 目录生成构建二进制文件 `nwpc_data_client`：

```bash
make
```

### 发布构建 (GoReleaser)

使用 GoReleaser 实现交叉编译，需要安装 GoReleaser 工具：

```bash
# 快照构建 (不需要 tag)
goreleaser build --snapshot --clean

# 发布构建 (需要 git tag)
goreleaser release --clean
```

### Windows 交叉编译 (PowerShell)

Windows 如果未安装 GoReleaser，可使用备用脚本实现交叉编译：

```powershell
# Linux amd64
.\build_cross.ps1 -Arch amd64

# Linux arm64
.\build_cross.ps1 -Arch arm64
```

## 快速开始

使用项目内置的数据配置查找 CMA-HPC2023 上的数值预报业务系统数据路径。
下面代码查找 CMA-GFS 预报系统 2026 年 6 月 22 日 12 时次 24 小时的基础 GRIB2 数据文件。

```bash
./bin/nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2026062212 \
    --forecast-time 24h
```

输出：

```text
/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2026062212/data/output/grib2_orig/gmf.gra.2026062212024.grb2
```

支持批量监视数据文件是否生成，并在文件生成时执行特定命令。
下面代码会持续检查 CMA-GFS 预报系统 2026 年 6 月 23 日 12 时次 000 小时和 024 小时的基础 GRIB2 数据文件，
并在找到时将文件拷贝到归档目录中。

```bash
echo "000h 024h" |
    nwpc_data_client check local 2026062312 \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --location-level runtime \
    --execute-command 'cp {{.FilePath}} /g7/JOB_TMP/wangdp/'
```

输出类似：

```text
INFO[0000] got check task for 2026062312 + 000h         
INFO[0000] got check task for 2026062312 + 024h         
INFO[0000] sleeping before check...0s                    forecast_time=000h
INFO[0000] checking begin...                             forecast_time=000h
INFO[0000] sleeping before check...10s                   forecast_time=024h
INFO[0000] checking... 0/2880                            forecast_time=000h
INFO[0000] checking exist...success: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2026062312/data/output/grib2_orig/gmf.gra.2026062312000.grb2  forecast_time=000h
INFO[0000] checking size... 0/2880                       forecast_time=000h
INFO[0000] checking size...changed 0/2880                forecast_time=000h
INFO[0005] checking size...success 1/2880                forecast_time=000h
INFO[0005] file is available, run command...             forecast_time=000h
INFO[0005] running command <cp /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2026062312/data/output/grib2_orig/gmf.gra.2026062312000.grb2 /g7/JOB_TMP/wangdp/> ...  forecast_time=000h
INFO[0005] run command success                           forecast_time=000h
INFO[0010] checking begin...                             forecast_time=024h
INFO[0010] checking... 0/2880                            forecast_time=024h
INFO[0010] checking exist...success: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2026062312/data/output/grib2_orig/gmf.gra.2026062312024.grb2  forecast_time=024h
INFO[0010] checking size... 0/2880                       forecast_time=024h
INFO[0010] checking size...changed 0/2880                forecast_time=024h
INFO[0015] checking size...success 1/2880                forecast_time=024h
INFO[0015] file is available, run command...             forecast_time=024h
INFO[0015] running command <cp /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2026062312/data/output/grib2_orig/gmf.gra.2026062312024.grb2 /g7/JOB_TMP/wangdp/> ...  forecast_time=024h
INFO[0015] run command success                           forecast_time=024h
INFO[0015] exiting
```

## 数据查找

`nwpc_data_client local` 使用 YAML 配置文件在 CMA-HPC2023 上查找本地文件。

### 用法

```bash
nwpc_data_client local \
    --data-type <data/type> \
    --start-time <YYYYMMDDHH> \
    --forecast-time <FFFh> \
    [--member <MMM>]
```

`nwpc_data_client local` 支持的所有参数：

| 标志                  | 说明                                                   | 默认值                           |
|----------------------|------------------------------------------------------|-------------------------------|
| `--data-config-dir`  | 配置文件目录，与 `--data-type` 配合使用。                         | 默认为空，使用内置到程序中的配置              |
| `--data-type`        | 用于在配置目录中定位配置文件的数据类型。                                 | 在未设置 `--data-config-file` 时必填 |
| `--data-config-file` | 单个配置文件路径，设置后将忽略 `--data-config-dir` 和 `--data-type`。 | -                             |
| `--location-level`   | 逗号分隔的存储层级，例如 `runtime,archive`。                      | 所有层级                          |
| `--start-time`       | 起报时间，`YYYYMMDDHH`。                                   | 必填                            |
| `--forecast-time`    | 预报时效，Go 时长格式，例如 `FFFh` 或 `FFFhMMm`。                  | 必填                            |
| `--member`           | 集合成员，格式为字符串，`MMM`。例如 `000` 或 `014`。                  | -                             |
| `--show-types`       | 列出可用的数据类型并退出。可以和 `--data-config-dir` 结合使用            | `false`                       |
| `--debug`            | 启用调试日志。                                              | `false`                       |

### 示例

#### 使用内置配置查找

查找 CMA-GFS 预报系统 2026 年 6 月 22 日 12 时次 24 小时基础 GRIB2 数据产品文件：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2026062212 \
    --forecast-time 24h
```

#### 使用自定义配置查找

使用 `--data-config-dir` 从自定义目录加载 YAML 配置。
CMA-GFS 业务系统后处理任务中使用如下的代码查找模式积分输出的 modelvar 文件。

```bash
nwpc_data_client local \
  --data-config-dir=/g1/u/op_post/OPER/UTIL/nwpc_data_client/config/local_hpc2023 \
  --data-type=cma_gfs_gmf/v4.2.3/oper/bin/modelvar \
  --location-level=all \
  --start-time=2026062300 \
  --forecast-time=024h
```

#### 使用单个配置文件

直接使用单个数据配置文件，并将忽略 `--data-config-dir` 和 `--data-type` 两个参数：

```bash
nwpc_data_client local \
    --data-config-file /path/to/config.yaml \
    --start-time 2026062300 \
    --forecast-time 48h
```

#### 列出可用的数据类型

列出二进制程序中内置的所有数据类型：

```bash
nwpc_data_client local --show-types
```

列出指定配置目录中所有可用的数据类型：

```bash
nwpc_data_client local --show-types \
  --data-config-dir /g1/u/op_post/OPER/UTIL/nwpc_data_client/config/local_hpc2023 
```

## 数据监视

`nwpc_data_client check local` 对一组预报时效对应的本地数据路径进行轮询，并可在找到每个文件并且文件大小不再变化后可选地执行命令。
通常用于业务工作流中等待模式输出，然后触发下游处理。

`start_time` 以位置参数形式传入，格式为 `YYYYMMDDHH`。

### 用法

`nwpc_data_client check local` 支持的所有参数：

| 标志                   | 说明                                                   | 默认值     |
|----------------------|------------------------------------------------------|---------|
| `--checker-config`   | 数据监视配置文件路径。包含以下所有命令行参数，如果设置命令行参数，将覆盖配置文件中的值。         | -       |
| `--data-config-dir`  | 配置文件目录，与 `--data-type` 配合使用。                         | 默认为空，使用内置到程序中的配置              |
| `--data-type`        | 用于在配置目录中定位配置文件的数据类型。                                 | 在未设置 `--data-config-file` 时必填 |
| `--data-config-file` | 单个配置文件路径，设置后将忽略 `--data-config-dir` 和 `--data-type`。 | -                             |
| `--location-level`   | 逗号分隔的存储层级，例如 `runtime,archive`。                      | 所有层级                          |
| `--max-check-count`  | 单个预报时效的最大检查轮数。                                       | `2880`  |
| `--check-interval`   | 轮询间隔，Go 时长格式（`30s`、`1m` 等）。                          | `5s`    |
| `--delay-time`       | 预报时效之间的交错延迟，Go 时长格式。                                 | `10s`   |
| `--execute-command`  | 文件稳定后执行的命令，Go 模板字符串。                                 | -       |
| `--debug`            | 启用调试日志。                                              | `false` |

### 轮询行为

对于每个预报时效，检查器将执行以下步骤：

1. 等待一个短时间的交错延迟（`--delay-time`，默认 `10s`）。
2. 重复检查解析后的路径，直到文件存在。
3. 继续轮询，直到文件大小不再变化，表示文件已写入完成。
4. 如果配置了 `--execute-command` 或者在配置文件中设置了命令列表，则运行这些命令。

如果在达到 `--max-check-count` 之前文件未找到或未稳定，检查器将以错误退出。

### 通过标准输入传入预报时效

当未提供 `--checker-config` 时，预报时效从 `stdin` 读取，每行一个，或以空白字符分隔，格式为 `FFFh` 或 `FFFhMMm`（例如 `000h`、`024h`、`003h10m`）。

下面命令来自 CMA-GFS 产品制作系统的数据检查任务，等待基础 GRIB2 数据输出，找到数据后通知 ecFlow，用于触发后续的绘图任务。

```bash
seq 0 3 240 | awk '{b=$1"h"; print b}' |
    nwpc_data_checker local \
      --data-config-dir=/g1/u/op_post/OPER/UTIL/nwpc_data_client/config/local_hpc2023 \
      --location-level=all \
      --data-type=cma_gfs_gmf/v4.2.3/oper/grib2/orig \
      --max-check-count=2880 \
      --check-interval=5s \
      --delay-time 1s \
      --execute-command 'ecflow_client --event grib2checked_{{.Forecast}}' \
      2026062600
```

`--execute-command` 模板支持与配置文件相同的变量（见[模板变量](#模板变量)），并额外支持 `.FilePath`，其值为解析后的文件路径。

### 检查器运行时配置

对于包含长命令字符串或多个下游步骤的工作流，可通过 `--checker-config`（或 `-c`）使用 YAML 运行时配置文件。
除位置参数 `start_time` 外的所有命令行参数都可以移入该文件，命令行上给出的任何参数都会覆盖文件中的值。

```yaml
# checker-config.yaml
data_type: cma_gfs_gmf/current/grib2/orig
location_levels: runtime
max_check_count: 2880
check_interval: 5s
delay_time: 0s
debug: false

forecast_times:
  - 000h
  - 024h
  - 048h

# 只能使用 execute_command（单个命令）或 execute_commands（列表）其中之一。
execute_commands:
  - 'echo "found {{.FilePath}}"'
  - '/app/postprocess.sh {{.FilePath}}'
```

运行方式：

```bash
nwpc_data_client check local 2025052900 --checker-config checker-config.yaml
```

如果配置文件中省略了 `forecast_times`，预报时效仍将按照原有行为从 `stdin` 读取。

`execute_command` 字段接受单个 Go 模板字符串。
`execute_commands` 字段接受模板字符串列表，每个文件稳定后按顺序执行。
如果任何命令失败，检查器将以错误退出。

> **注意**
> 
> 对于同时存在于检查器运行时配置和 CLI 标志中的任何字段，**CLI 标志始终覆盖 YAML 文件**。

> **YAML 引号提示**
>
> 某些 YAML 解析器（包括 `nwpc_data_client check` 使用的解析器）难以处理包含冒号加空格的双引号标量，
> 例如 `"echo file path: {{.FilePath}}"`。
> 为避免解析错误，请使用单引号包裹命令字符串，例如 `'echo "file path: {{.FilePath}}"'`。


## 数据配置

数据类型由 YAML 配置文件描述。
当未提供 `--data-config-dir` 时，`nwpc_data_client` 使用构建时嵌入二进制中的配置（见[重新生成内置配置](#重新生成内置配置)）。
实际应用时推荐在自定义目录中定义数据配置文件。

### 文件格式

配置文件是 YAML 格式文档，并支持 Go 模板语法，包含以下字段：

```yaml
default: NOTFOUND

file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2
file_names: []

paths:
  - type: local
    level: runtime
    path: /some/dir/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/

  - type: local
    level: archive
    path: /another/dir/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/
```

| 字段 | 说明 |
|------|------|
| `default` | 未找到文件时返回的值。大多数配置中默认为 `NOTFOUND`。 |
| `file_name` | 单个文件名模板。`file_name` 或 `file_names` 必须设置其一。 |
| `file_names` | 文件名模板列表。它们按顺序尝试，位于可选的 `file_name` 之后。适用于同一份数据可能以多个名称存储的场景。 |
| `paths` | 要搜索的目录条目列表。按顺序检查；第一个存在的文件获胜。 |

`paths` 中的每个条目包含：

| 字段 | 说明 |
|------|------|
| `type` | 路径类型。对于 `nwpc_data_client local` 和 `nwpc_data_client check local`，该值为 `local`。 |
| `level` | 存储层级，例如 `runtime`、`archive` 或 `all`。使用 `--location-level` 可将搜索限制到特定层级。空层级或 `all` 可匹配任何过滤器。 |
| `path` | 目录模板。解析后的文件名会拼接在该目录后。 |

### 模板变量

配置文件作为 Go 模板执行，使用 `{{` 和 `}}` 作为分隔符。以下变量可用：

| 变量 | 类型 | 示例 | 说明 |
|------|------|------|------|
| `.StartTime` | `time.Time` | - | 解析后的起报时间。 |
| `.ForecastTime` | `time.Duration` | - | 解析后的预报时效。 |
| `.Year` | string | `2025` | 起报时间的四位年份。 |
| `.Month` | string | `05` | 起报时间的零填充月份。 |
| `.Day` | string | `29` | 起报时间的零填充日期。 |
| `.Hour` | string | `00` | 起报时间的零填充小时。 |
| `.ForecastHour` | string | `024` | 零填充的预报小时，三位数。 |
| `.ForecastMinute` | string | `00` | 零填充的预报分钟，两位数。 |
| `.Member` | string | - | 集合成员，适用时。 |

### 模板辅助函数

模板辅助函数也可用于更复杂的表达式：

```yaml
file_name: gmf.gra.{{generateStartTime .StartTime}}{{getForecastHour .ForecastTime}}.grb2
```

| 函数 | 说明 |
|------|------|
| `generateStartTime` | 将起报时间格式化为 `YYYYMMDDHH`。 |
| `getYear` | 提取年份。 |
| `getMonth` | 提取月份。 |
| `getDay` | 提取日期。 |
| `getHour` | 提取小时。 |
| `generateForecastTime` | 将预报时效格式化为 `FFFhMMm`。 |
| `getForecastHour` | 提取预报小时。 |
| `getForecastMinute` | 提取预报分钟。 |

### 解析流程

1. 从内置配置、自定义配置目录或单个配置文件加载 YAML 内容。
2. `nwpc_data_client` 使用起报时间、预报时效和可选的集合成员将 YAML 解析为 Go 模板。
3. 将解析后的目录模板和文件名模板拼接。
4. 按顺序检查每个结果路径；返回第一个存在的文件，如果没有则返回 `default`。

### 示例

```yaml
# conf/local/cma_gfs_gmf/current/grib2/orig.yaml
default: NOTFOUND

file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/data/output/grib2_orig/

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
```

对于起报时间 `2025052900` 和预报时效 `024h`，第一个路径解析为：

```text
/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2025052900/data/output/grib2_orig/gmf.gra.2025052900024.grb2
```

如果第一个路径中不存在该文件，工具将按照顺序查找第二个路径：

```text
/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/2025052900/ORIG/gmf.gra.2025052900024.grb2
```

## 开发

### 运行测试

运行单元测试：

```bash
make test
```

这将运行 `go test ./...`。

集成测试位于 `tests/integration/` 下，使用 `integration` 构建标签（`//go:build integration`）。
集成测试需要在 CMA-HPC2023-SC1 上运行，在运行时会检查数据文件是否存在，在数据不可用时跳过对应测试条目。
`make test-integration` 会自动构建二进制文件并运行集成测试：

```bash
make test-integration
```

运行单个集成测试：

```bash
go test -tags=integration -run TestCMA_GFS_GMF ./tests/integration/...
```

集成测试辅助程序会首先使用环境变量 `NWPC_DATA_CLIENT_PROGRAM` 和 `NWPC_DATA_CLIENT_CONFIG_DIR` 中设置的值，
否则会使用默认值 `./bin/nwpc_data_client` 和 `./conf`。

### 重新生成内置配置

`conf/` 下的配置可以通过生成的 `common/config/config.autogen.go` 文件嵌入到编译后的可执行程序中。

```bash
make generate
```

上述命令会在 `common/config/generate` 中运行 `go generate`。

> **注意**
> 
> 在 `conf/` 下添加或编辑 YAML 配置后，需要重新运行生成新的 `common/config/config.autogen.go` 文件。

## 许可证

Copyright &copy; 2019-2026 developers at cemc-oper.

`nwpc-data-client` 基于 [The MIT License](https://opensource.org/licenses/MIT) 授权。
