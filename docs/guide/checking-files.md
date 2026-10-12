# 等待文件并执行命令

`nwpc_data_client check local <start_time>` 针对一组预报时效轮询本地数据
路径，每个文件大小稳定（两次检查之间不再变化）后执行模板化命令。

- 起报时次为**位置参数**。
- 预报时效来自 `--checker-config` 的 `forecast_times` 字段（推荐），
  或标准输入中以空白分隔的列表。
- 每个预报时效在错开的 `--delay-time` 后并发检查。
- 命令模板中除配置模板变量外还可使用 `{{.FilePath}}`。
- `--execute-command` 设置单条命令；多条命令（`execute_commands`）
  只能通过 checker 配置 YAML 提供，两者互斥。

checker 配置字段与 CLI 参数的覆盖关系见
[Checker 配置参考](../reference/checker-config.md)。

> 本文档正在编写中：补充完整 checker 配置示例与运行流程说明。
