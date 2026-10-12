# 模板变量与函数

数据配置文件和 `check local` 的命令模板都使用 Go 模板（`{{...}}` 定界符）。
本页列出全部可用变量与函数（实现见 `common/config_template.go` 与
`common/time.go`）。

## 模板变量

变量来自 `ConfigTemplateVariable`，由起报时次、预报时效和集合成员计算得到：

| 变量 | 类型 | 说明 | 示例值 |
|---|---|---|---|
| `{{.StartTime}}` | `time.Time` | 起报时间 | — |
| `{{.ForecastTime}}` | `time.Duration` | 预报时效 | — |
| `{{.Year}}` | string | 起报年，4 位 | `2026` |
| `{{.Month}}` | string | 起报月，2 位补零 | `06` |
| `{{.Day}}` | string | 起报日，2 位补零 | `24` |
| `{{.Hour}}` | string | 起报时，2 位补零 | `00` |
| `{{.ForecastHour}}` | string | 预报时效小时数，3 位补零 | `024` |
| `{{.ForecastMinute}}` | string | 预报时效的分钟部分（总分钟对 60 取余），2 位补零 | `30` |
| `{{.Member}}` | string | 集合成员编号（`--member` 传入），非集合预报为空 | `014` |

`check local` 的命令模板额外提供：

| 变量 | 说明 |
|---|---|
| `{{.FilePath}}` | 检测到并已稳定的文件完整路径 |

## 模板函数

### 起报时间偏移

`generateStartTime` 对起报时间加减整数小时，返回新的 `time.Time`，
通常与时间分量函数通过管道组合使用：

```
{{generateStartTime .StartTime -3 | getYear}}
{{generateStartTime .StartTime 24 | getDay}}
```

### 时间分量

| 函数 | 签名 | 说明 |
|---|---|---|
| `getYear` | `func(time.Time) string` | 4 位年 |
| `getMonth` | `func(time.Time) string` | 2 位月 |
| `getDay` | `func(time.Time) string` | 2 位日 |
| `getHour` | `func(time.Time) string` | 2 位时 |

### 预报时效计算

`generateForecastTime` 将预报时效加上一个时长字符串（如 `"3h"`、
`"1h30m"`），返回新的 `time.Duration`：

```
{{generateForecastTime .ForecastTime "3h" | getForecastHour | printf "%03d"}}
```

### 预报时效分量

| 函数 | 签名 | 说明 |
|---|---|---|
| `getForecastHour` | `func(time.Duration) int` | 时效总小时数（取整） |
| `getForecastMinute` | `func(time.Duration) int` | 时效分钟部分（总分钟对 60 取余） |

两者返回整数，配合 Go 内置的 `printf` 控制格式：

```
{{.ForecastTime | getForecastHour | printf "%03d"}}
{{.ForecastTime | getForecastMinute | printf "%02d"}}
```

## 完整示例

起报时次 `2026062400`、预报时效 `1h30m` 时：

| 模板 | 渲染结果 |
|---|---|
| `{{.Year}}{{.Month}}{{.Day}}{{.Hour}}` | `2026062400` |
| `{{.ForecastHour}}` | `001` |
| `{{.ForecastMinute}}` | `30` |

管道组合示例：

```
{{generateStartTime .StartTime -12 | getYear}}        --> 2026（2026-06-23 12:00 的年份）
{{.ForecastTime | getForecastHour | printf "%03d"}}   --> 001
{{.ForecastTime | getForecastMinute | printf "%02d"}} --> 30
```
