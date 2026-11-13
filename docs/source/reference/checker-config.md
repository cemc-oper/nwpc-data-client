# Checker 配置

CheckerConfig 是 `check local` 的运行时配置（定义于
`internal/checker/config.go`），与数据路径配置（DataConfig）相互独立。
它通过 `--checker-config` 从 YAML 文件加载，并与 CLI 参数合并。

## 字段

| YAML 字段 | CLI 参数 | 默认值 | 说明 |
|---|---|---|---|
| `data_type` | `--data-type` | — | 数据类型。与 `data_config_file` 至少提供一个 |
| `data_config_dir` | `--data-config-dir` | `""` | 完整配置树目录，可选，不能单独使用；不是内嵌配置的增量覆盖层 |
| `data_config_file` | `--data-config-file` | `""` | 单个配置文件。设置后忽略 `data_config_dir` 与 `data_type` |
| `location_levels` | `--location-level` | `""` | 逗号分隔的位置级别，空表示搜索全部 |
| `max_check_count` | `--max-check-count` | `2880` | 每个预报时效的最大检查次数，必须大于 0；存在性检查与大小稳定检查共用同一预算 |
| `check_interval` | `--check-interval` | `5s` | 两次检查的间隔（Go duration 格式），必须大于 0 |
| `check_interval_jitter` | `--check-interval-jitter` | `0.1` | 逐轮抖动比例：每轮实际间隔在 `interval×(1±j)` 内随机，取值 `[0, 0.5]`，`0` 关闭 |
| `max_concurrent_stable_checks` | `--max-concurrent-stable-checks` | `4` | 同时处于大小稳定观察的文件数上限，必须大于 0（仅批量轮询模式生效） |
| `delay_time` | `--delay-time` | `10s` | 相邻预报时效检查任务之间的错开等待（Go duration 格式） |
| `debug` | `--debug` | `false` | 调试日志 |
| `check_mode` | `--check-mode` | `per_time` | 轮询模式：`per_time`（每时效一协程）或 `batched`（单循环 readdir 批量发现） |
| `arrival_order` | `--arrival-order` | `auto` | 到达顺序覆盖：`auto` 跟随数据配置的 `_arrival_order`，`sequential` 启用顺序门控，`unordered` 关闭（仅 batched 模式生效） |
| `forecast_times` | —（仅 YAML） | — | 预报时效字符串列表，如 `["0h", "6h", "24h"]`；空列表或缺省时由 `check local` 从标准输入读取 |
| `execute_command` | `--execute-command` | `""` | 单条命令模板 |
| `execute_commands` | —（仅 YAML） | — | 多条命令模板列表，与 `execute_command` 互斥 |
| `command_timeout` | `--command-timeout` | `30m` | 每条命令的超时配置（Go duration 格式），`0` 表示禁用，负数非法 |
| `report_file` | `--report-file` | `""` | JSON 运行报告路径配置，空表示不生成 |

`dry_run` 只提供 CLI 参数 `--dry-run`，不是可持久化的 YAML 字段。

> **运行时边界**
>
> `command_timeout` 已接入每条命令的实际执行，默认 30 分钟，`0` 禁用。
> `report_file` 会生成 versioned JSON 运行报告；空值禁用，`--dry-run` 在任何情况下
> 都不会写报告。

## 合并与校验规则

YAML 与 CLI 按以下优先级合并（`MergeCheckerConfig`）：

1. 内置默认值；
2. YAML 中**显式出现**的字段，包括 `0`、`false` 和空字符串；
3. **显式给出**的 CLI 参数（未显式给出的参数不覆盖 YAML）。
   字段与 CLI 参数的对应关系由结构体的 `flag` 标签声明。

字段是否出现与字段值是否有效是两件事。例如，`check_interval_jitter: 0`
会关闭抖动，`command_timeout: 0` 会禁用命令超时，而 `max_check_count: 0`
会进入校验并被拒绝，不会回退到默认值 2880。显式空的 `check_interval`
同样是非法 duration。

特殊规则：

- `forecast_times` 与 `execute_commands` 没有对应 CLI 参数，只能来自 YAML；
- 显式给出 `--execute-command` 时，YAML 中的 `execute_commands` 被丢弃
  （两者互斥）。

加载与合并完成后校验（任一不满足即报错退出）：

- YAML 使用严格解码，未知字段报错并包含文件路径；
- `data_type` 或 `data_config_file` 必须提供其一，`data_config_dir` 不能单独使用；
- `check_interval`、`delay_time` 和 `command_timeout` 必须是合法 Go duration，
  其中检查间隔必须大于 0，命令超时不得为负数；
- `max_check_count` 和 `max_concurrent_stable_checks` 必须大于 0；
- `check_interval_jitter` 必须在 `[0, 0.5]` 内；
- `check_mode` 必须是 `per_time` 或 `batched`；
- `arrival_order` 必须是 `auto`、`sequential` 或 `unordered`；
- `forecast_times` 逐项必须可解析；静态 lint 允许空列表，因为运行时可以从标准输入读取；
- `execute_command` 与 `execute_commands` 不得同时存在，所有命令模板必须可解析。

## 部署前校验

提供两层互补的检查，详见[校验、错误与运行报告](../guide/check-local/validation-and-reports.md)：

- `check lint`：只做 checker YAML 的严格加载、默认值合并和静态字段校验；
  不需要起报时次、完整 data config 或非空 `forecast_times`；
- `check local --dry-run`：使用真实配置合并和 data config 展开生成运行计划，
  验证起报时次、预报时效、候选路径和命令模板，但不访问候选业务文件。

## JSON 报告结构

`report_file` 非空时生成以下 versioned 结构；完整状态和耗时口径见
[JSON 运行报告](../guide/check-local/validation-and-reports.md#json-运行报告)。

```json
{
  "report_version": 1,
  "start_time": "2026102600",
  "data_type": "cma_gfs_gmf/current/grib2/orig",
  "check_mode": "batched",
  "arrival_order": "unordered",
  "started_at": "2026-10-26T00:00:00Z",
  "finished_at": "2026-10-26T00:01:00Z",
  "status": "success",
  "summary": {
    "total": 1,
    "succeeded": 1,
    "failed": 0,
    "cancelled": 0,
    "pending": 0
  },
  "results": [
    {
      "forecast_time": "000h",
      "status": "success",
      "file_path": "/path/to/file",
      "check_count": 3,
      "wait_seconds": 20.0,
      "command_seconds": 0.1,
      "error": ""
    }
  ]
}
```

报告通过同目录临时文件和重命名原子发布；写失败会非零退出。文件命名、保留和清理
由调度侧负责，client 不做日志轮转，也不把报告写入 stdout。

## 抖动（jitter）

为避免同一时次批量启动的进程在整点对齐、同步冲击存储元数据服务，轮询
带有两级随机抖动：

- **启动抖动**：进程进入轮询前随机等待 `0 ~ check_interval`（全量程，
  不可配置），打散进程间相位；
- **逐轮抖动**：每轮间隔在 `check_interval×(1±j)` 内随机，`j` 即
  `check_interval_jitter`（默认 0.1，`0` 关闭）。抖动按比例随间隔联动，
  额外发现延迟的上界为 `j×check_interval`。

dry-run 不执行启动抖动或逐轮抖动。

## 完整示例

```yaml
# checker.yaml
data_type: cma_gfs_gmf/current/grib2/orig
location_levels: runtime
max_check_count: 360
check_interval: 10s
check_interval_jitter: 0
delay_time: 5s
forecast_times:
  - 0h
  - 6h
  - 12h
  - 24h
execute_commands:
  - 'echo "file ready: {{.FilePath}}"'
  - "echo cycle {{.Year}}{{.Month}}{{.Day}}{{.Hour}} + {{.ForecastHour}}"
```

使用：

```bash
nwpc_data_client check local 2026062400 --checker-config=./checker.yaml
```

命令模板中除 `{{.FilePath}}` 外，还可以使用全部
[模板变量与函数](template-variables.md)。
