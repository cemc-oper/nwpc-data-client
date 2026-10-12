# 模板变量与函数

数据配置与命令模板可用的变量（`ConfigTemplateVariable`）：

| 变量 | 说明 |
|---|---|
| `{{.StartTime}}` | 起报时次（time.Time） |
| `{{.ForecastTime}}` | 预报时效（time.Duration） |
| `{{.Year}}` `{{.Month}}` `{{.Day}}` `{{.Hour}}` | 起报时间的年/月/日/时（两位补零） |
| `{{.ForecastHour}}` | 预报时效小时数（三位补零，如 `024`） |
| `{{.ForecastMinute}}` | 预报时效分钟部分 |
| `{{.Member}}` | 集合成员编号 |

模板函数（定义于 `common/time.go`）：

- `generateStartTime` / `getYear` 等时间分量函数
- `generateForecastTime` / `getForecastHour` 等预报时效函数

`check local` 的命令模板额外提供 `{{.FilePath}}`（检测到的文件路径）。

> 本文档正在编写中：补全函数签名、参数与示例。
