# 编写数据配置

数据配置是 YAML 文件，整体作为 Go 模板（`{{...}}` 定界符）先渲染再解析：

```yaml
default: NOTFOUND
file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2
file_names: []
paths:
  - type: local
    level: runtime
    path: /some/dir/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/
```

常用模板变量：`{{.Year}}`、`{{.Month}}`、`{{.Day}}`、`{{.Hour}}`、
`{{.ForecastHour}}`、`{{.ForecastMinute}}`、`{{.Member}}`、
`{{.StartTime}}`、`{{.ForecastTime}}`。
完整清单见[模板变量与函数](../reference/template-variables.md)。

路径条目的 `level` 字段与 `--location-level` 参数配合，控制搜索范围。

> 本文档正在编写中：补充 file_name / file_names 的区别、多路径搜索、
> 模板函数用法与完整示例。
