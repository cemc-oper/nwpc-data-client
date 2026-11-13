# 轮询模式与到达顺序

`check local` 提供 `per_time` 和 `batched` 两种轮询模式。时效较少时使用默认的 `per_time` 即可；需要等待大量时效时，可以使用 `batched` 减少重复目录访问。

## 逐时效轮询（per_time）

默认的 `per_time` 模式为每个预报时效启动一个协程独立轮询。相邻任务通过 `--delay-time` 错开启动，避免所有时效在同一时刻集中访问存储。每个任务独立执行文件查找、大小稳定检查和命令执行。

## 批量轮询（batched）

时效很多时（如逐 10 分钟、上百个时效），逐时效轮询会产生大量重复的目录访问。批量轮询模式把所有待等时效收敛到**单一轮询循环**：每轮对每个目录只做一次 `readdir`，在内存中按 `paths` 顺序完成优先级匹配。

启用方式（YAML 与 CLI 均可）：

```yaml
# checker.yaml
check_mode: batched
data_type: cma_gfs_gmf/current/grib2/orig
forecast_times: ["0h", "3h", "6h"]
execute_command: "echo ready: {{.FilePath}}"
```

```bash
nwpc_data_client check local 2026062400 --check-mode=batched ...
```

### 与 per_time 的差异

- **目录访问量**：每轮每目录一次 `readdir`，与时效数量无关；`per_time`
  为每时效每轮一次或多次 `stat`；
- **发现时刻**：各时效的发现对齐到统一轮询周期边界，不再有逐时效的
  `--delay-time` 错峰——**`--delay-time` 在 batched 模式下不适用**（无
  逐协程错峰，设置后会被忽略）；进程内错相由统一循环承担，进程间错相
  由 [jitter](../../reference/checker-config.md#抖动jitter) 承担；
- **稳定观察并发上限**：同时处于大小稳定观察的文件数受
  `max_concurrent_stable_checks`（默认 4）约束，超出的命中时效排到下一
  轮；`per_time` 模式各协程天然独立，不适用此上限；
- **检查次数预算**：`max_check_count` 仍按时效独立计数，语义不变；
  超时未等到同样以非零码退出；
- **路径优先级**：`paths` 顺序优先（runtime 优先于 archive）在内存匹配
  中完全保持，目录扫描顺序不影响结果；
- **顺序门控**：数据配置声明 `_arrival_order: sequential` 时（本仓库
  内置的模式积分产品均已声明），每轮只检查最早未到的时效，其命令执行
  完后才放开下一个——目录扫描数收敛到 1。可用 `--arrival-order` 双向
  覆盖。

## 如何选择

适用场景：大批量时效、业务化部署。时效少（几个）时两种模式差别不大，保持默认 `per_time` 即可。

两种模式的完整参数及默认值见[参考：Checker 配置](../../reference/checker-config.md)。
