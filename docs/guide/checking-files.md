# 等待文件并执行命令

`nwpc_data_client check local <start_time>` 针对一组预报时效轮询本地
数据路径；每个文件生成且大小稳定后，执行对应的命令。

## 工作流程

1. 起报时次为**位置参数**（`YYYYMMDDHH`）；
2. 预报时效列表来自 `--checker-config` 的 `forecast_times` 字段
   （推荐），或标准输入中以空白分隔的记号；
3. 每个预报时效启动一个并发检查任务，相邻任务之间错开
   `--delay-time`（默认 10s）；
4. 每个任务按 `local` 相同的配置解析规则查找文件：文件不存在则按
   `--check-interval`（默认 5s）间隔重试；
5. 文件出现后进入大小稳定检查：大小**非零且连续两次检查相同**即判定
   稳定（存在性检查与大小检查共用 `--max-check-count` 预算，默认 2880 次）；
6. 稳定后执行命令模板；命令失败会记录错误并终止（`log.Fatal`）。

## 使用 checker 配置文件（推荐）

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
  - 24h
execute_command: "echo file ready: {{.FilePath}}"
```

```bash
nwpc_data_client check local 2026062400 --checker-config=./checker.yaml
```

全部字段及与 CLI 参数的合并规则见[参考：Checker 配置](../reference/checker-config.md)。

## 从标准输入提供预报时效

```bash
echo "0h 3h 6h" | nwpc_data_client check local 2026062400 \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --execute-command="echo ready: {{.FilePath}}"
```

`forecast_times` 与标准输入不能同时为空，否则报错退出；
YAML 中的 `forecast_times` 优先于标准输入。

## 命令模板

命令是 Go 模板，除全部[配置模板变量](../reference/template-variables.md)
外还可使用 `{{.FilePath}}`（稳定文件的路径）：

```yaml
execute_commands:
  - "python3 process.py --input {{.FilePath}}"
  - "echo {{.Year}}{{.Month}}{{.Day}}{{.Hour}} + {{.ForecastHour}} done"
```

- `--execute-command`（或 YAML `execute_command`）指定**一条**命令；
- **多条**命令只能通过 YAML 的 `execute_commands` 列表提供；
- 两者互斥：显式给出 `--execute-command` 时 YAML 的 `execute_commands`
  被丢弃；
- 命令按列表顺序依次执行，共享进程的标准输入输出。

## 常用调优

```bash
nwpc_data_client check local 2026062400 --checker-config=./checker.yaml \
    --check-interval=30s --max-check-count=120
```

- 长时效文件增大 `--check-interval` 降低轮询开销；
- `--max-check-count` × `--check-interval` 大致是单个时效的最长等待
  时间（如默认 2880 × 5s = 4 小时）；超时未稳定该任务报错，整体以
  非零码退出；
- `--delay-time` 避免所有时效在同一时刻集中访问存储。
