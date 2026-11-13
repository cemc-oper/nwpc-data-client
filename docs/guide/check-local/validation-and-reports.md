# 校验、错误与运行报告

## 静态校验 check lint

`check lint` 在合并或部署前检查 checker YAML 本身。它可以接收一个或多个文件、目录，目录会递归查找 `.yaml` 文件并聚合全部错误：

```bash
# 校验单份配置。
nwpc_data_client check lint ./checker.yaml

# 递归校验多个目录，也可以混合文件和目录。
nwpc_data_client check lint ./checker/system_a ./checker/system_b
```

校验内容包括严格 YAML 字段、默认值合并、duration、正整数、预报时效格式、data config 组合和命令模板语法。任一文件失败时非零退出；空目录成功。

静态 lint 明确不做以下事情：

- 不需要起报时次，不加载 data config；
- 不要求 `forecast_times` 非空，因为真实运行时可以从标准输入读取；
- 不读取标准输入，不访问业务数据目录；
- 不执行命令，不创建运行报告。

因此它能发现字段拼写错误和静态值错误，但不能发现 data type 缺失、data config 模板展开失败等问题；这些由 dry-run 检查。

## 运行前验证 check local --dry-run

`--dry-run` 使用与真实 `check local` 相同的配置合并和运行计划构造：

```bash
nwpc_data_client check local 2026062400 \
    --checker-config=./checker.yaml \
    --data-config-dir=/path/to/complete/config \
    --dry-run
```

输出包括最终 `check_mode`、有效 arrival order、时效数量、逐时效候选路径和原始命令模板。候选路径只是按配置展开的计划，**不表示文件已经存在，也不表示真实查找已经选中该路径**。命令模板使用哑 `FilePath` 验证 Go 模板变量和函数，但不会作为 shell 命令执行。

dry-run 会读取 checker/data 配置，并在 YAML 未提供时效时读取标准输入；时效优先级与真实运行一致。`--data-config-dir` 仍是完整配置树的全量替代，缺失 data type 不会回退到内嵌配置。任一时效的 data config 或命令模板展开失败都会非零退出。

它不会执行以下运行时动作：

- 不 `stat` 候选文件，不 `readdir` 候选目录；
- 不执行启动抖动、逐时效延迟或轮询等待；
- 不启动命令执行器，不执行命令；
- 即使设置了 `report_file`，也不写运行报告。

可用它在部署前验证完整配置组合；如果只需要检查 YAML 字段而还没有完整配置树或起报时次，使用更轻量的 `check lint`。

## 错误与取消

运行错误由 CLI 最外层写到标准错误，不附加命令 usage。普通检查或命令失败退出码为 1；收到 SIGHUP、SIGINT、SIGTERM 或 SIGQUIT 时退出码保持为 3。

取消覆盖 checker 自身的启动抖动、逐时效延迟、逐轮轮询等待、稳定文件等待以及正在执行的外部命令。fail-fast 或取消后不再启动新的检查和命令，已启动的协程收尾后才返回。

每条命令默认有 30 分钟超时，可通过 `command_timeout` / `--command-timeout` 调整，`0` 表示禁用。命令超时作为普通运行错误退出 1，外部信号取消仍退出 3；两种错误可区分。

> 命令取消和超时只保证命令依赖实际支持的直接进程终止，不承诺跨平台清理 shell 派生的整个子孙进程树。业务命令若会派生长驻子进程，需要单独验证其退出行为。

## JSON 运行报告

设置 `report_file` / `--report-file` 后，成功、检查失败、命令失败、命令超时和可处理的信号取消都会写出最终 JSON 报告；默认空值不生成文件，dry-run 也从不写报告。

```yaml
report_file: /path/to/runtime/check-2026062400.json
```

报告以 `report_version: 1` 标识 schema，包含运行状态、汇总及逐时效结果。总状态为 `success`、`failed` 或 `cancelled`；逐时效状态包括 `pending`、`success`、`check_failed`、`command_failed` 和 `cancelled`。fail-fast 时尚未开始的时效保留为 `pending`，已开始但未完成的时效为 `cancelled`。

`check_count` 是实际文件存在性与大小检查次数；`wait_seconds` 从首次实际检查计时，不包含 per-time 的启动错峰；`command_seconds` 单独记录命令执行耗时。报告不包含命令 stdout/stderr 或环境变量。

报告在目标目录内先写临时文件，再原子重命名，读取者不会看到半份 JSON。报告写失败会使 CLI 非零退出，即使检查本身已经成功。报告文件的命名、按起报时次区分、保留周期和清理由业务调度侧负责。

完整 JSON 结构示例见[参考：Checker 配置](../../reference/checker-config.md#json-报告结构)。
