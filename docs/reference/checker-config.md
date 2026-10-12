# Checker 配置

`check local` 的运行时配置（CheckerConfig，定义于 `cmd/checker_config.go`），
由 `--checker-config` YAML 文件加载，并与 CLI 参数合并：

- 每个字段带 `flag` 结构标签，映射到对应 CLI 参数名。
- 显式给出的 CLI 参数覆盖 YAML 中的值。
- 合并后校验：`data_type` 或 `data_config_file` 必须提供其一；
  `execute_command` 与 `execute_commands` 互斥。

主要字段：

| 字段 | 对应 CLI 参数 | 说明 |
|---|---|---|
| `data_type` | `--data-type` | 数据类型 |
| `data_config_file` | `--data-config-file` | 单个配置文件 |
| `forecast_times` | — | 预报时效列表（仅 YAML） |
| `execute_command` | `--execute-command` | 单条命令模板 |
| `execute_commands` | — | 多条命令模板（仅 YAML） |
| `delay_time` | `--delay-time` | 各时效检查前的错开等待 |

> 本文档正在编写中：补全字段表与完整 YAML 示例。
