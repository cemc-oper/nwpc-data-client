# Checker 配置

CheckerConfig 是 `check local` 的运行时配置（定义于
`internal/checker/config.go`），与数据路径配置（DataConfig）相互独立。
它通过 `--checker-config` 从 YAML 文件加载，并与 CLI 参数合并。

## 字段

| YAML 字段 | CLI 参数 | 默认值 | 说明 |
|---|---|---|---|
| `data_type` | `--data-type` | — | 数据类型。与 `data_config_file` 二选一 |
| `data_config_dir` | `--data-config-dir` | `""` | 配置目录，可选，不能单独使用 |
| `data_config_file` | `--data-config-file` | `""` | 单个配置文件。设置后忽略 `data_config_dir` 与 `data_type` |
| `location_levels` | `--location-level` | `""` | 逗号分隔的位置级别，空表示搜索全部 |
| `max_check_count` | `--max-check-count` | `2880` | 每个预报时效的最大检查次数（存在性检查与大小稳定检查共用同一预算） |
| `check_interval` | `--check-interval` | `5s` | 两次检查的间隔（Go duration 格式） |
| `check_interval_jitter` | `--check-interval-jitter` | `0.1` | 逐轮抖动比例：每轮实际间隔在 `interval×(1±j)` 内随机，取值 `[0, 0.5]`，`0` 关闭 |
| `max_concurrent_stable_checks` | `--max-concurrent-stable-checks` | `4` | 同时处于大小稳定观察的文件数上限（仅批量轮询模式生效） |
| `delay_time` | `--delay-time` | `10s` | 相邻预报时效检查任务之间的错开等待 |
| `debug` | `--debug` | `false` | 调试日志 |
| `check_mode` | `--check-mode` | `per_time` | 轮询模式：`per_time`（每时效一协程）或 `batched`（单循环 readdir 批量发现） |
| `arrival_order` | `--arrival-order` | `auto` | 到达顺序覆盖：`auto` 跟随数据配置的 `_arrival_order`，`sequential` 启用顺序门控，`unordered` 关闭（双向覆盖均有效；仅 batched 模式生效） |
| `forecast_times` | —（仅 YAML） | — | 预报时效字符串列表，如 `["0h", "6h", "24h"]`；缺省时从标准输入读取 |
| `execute_command` | `--execute-command` | `""` | 单条命令模板 |
| `execute_commands` | —（仅 YAML） | — | 多条命令模板列表，与 `execute_command` 互斥 |

## 合并与校验规则

YAML 与 CLI 按以下优先级合并（`MergeCheckerConfig`）：

1. 内置默认值；
2. YAML 中的非零值；
3. **显式给出**的 CLI 参数（未显式给出的参数不覆盖 YAML）。
   字段与 CLI 参数的对应关系由结构体的 `flag` 标签声明。

特殊规则：

- `forecast_times` 与 `execute_commands` 没有对应 CLI 参数，只能来自 YAML；
- 显式给出 `--execute-command` 时，YAML 中的 `execute_commands` 被丢弃
  （两者互斥）。

合并完成后校验（任一不满足即报错退出）：

- `data_type` 或 `data_config_file` 必须提供其一；
- `execute_command` 与 `execute_commands` 不得同时存在；
- `check_interval_jitter` 必须在 `[0, 0.5]` 内；
- `check_mode` 必须是 `per_time` 或 `batched`；
- `arrival_order` 必须是 `auto`、`sequential` 或 `unordered`。

## 抖动（jitter）

为避免同一时次批量启动的进程在整点对齐、同步冲击存储元数据服务，轮询
带有两级随机抖动：

- **启动抖动**：进程进入轮询前随机等待 `0 ~ check_interval`（全量程，
  不可配置），打散进程间相位；
- **逐轮抖动**：每轮间隔在 `check_interval×(1±j)` 内随机，`j` 即
  `check_interval_jitter`（默认 0.1，`0` 关闭）。抖动按比例随间隔联动，
  额外发现延迟的上界为 `j×check_interval`。

## 完整示例

```yaml
# checker.yaml
data_type: cma_gfs_gmf/current/grib2/orig
location_levels: runtime
max_check_count: 360
check_interval: 10s
delay_time: 5s
forecast_times:
  - 0h
  - 6h
  - 12h
  - 24h
execute_commands:
  - "echo file ready: {{.FilePath}}"
  - "echo cycle {{.Year}}{{.Month}}{{.Day}}{{.Hour}} + {{.ForecastHour}}"
```

使用：

```bash
nwpc_data_client check local 2026062400 --checker-config=./checker.yaml
```

命令模板中除 `{{.FilePath}}` 外，还可以使用全部
[模板变量与函数](template-variables.md)。
