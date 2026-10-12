# 数据配置格式

数据配置（DataConfig）的完整字段参考。格式约定与 JSON Schema 保持一致：

- 当前版本：`schemas/config/data-config.schema.json`
- v1 版本：`schemas/config/data-config-v1.schema.json`

顶层字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `_config_version` | int | 配置格式版本声明（可选，保留字段） |
| `default` | string | 未找到文件时的返回默认值 |
| `file_name` | string | 文件名模板 |
| `file_names` | list | 候选文件名模板列表 |
| `paths` | list | 候选目录条目，含 `type` / `level` / `path` |

> 本文档正在编写中：对照 `common.ValidateDataConfig` 补全字段约束说明。
