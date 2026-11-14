# 脚本接入与故障排查

本页汇总三个数据命令在 Shell 和业务调度中的输出约定，以及常见问题的定位方法。

## 在脚本中使用 local

`local` 找到文件时输出路径并以退出码 0 返回；**未找到时输出数据配置中的
`default` 值（通常为 `NOTFOUND`），退出码仍为 0**。因此不能只依靠 `set -e`
判断文件是否存在：

```bash
input_file="$(
    nwpc_data_client local \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --start-time=2026062400 \
        --forecast-time=24h
)" || exit 1

if [ "$input_file" = "NOTFOUND" ]; then
    printf '%s\n' "input file is not available" >&2
    exit 1
fi
```

自定义配置如果使用其他 `default` 值，脚本应判断该值，而不是固定写成 `NOTFOUND`。
可以使用 `--debug` 查看每个候选路径的检查过程。

## 获取 fetch 的输出路径

`fetch` 未找到文件、源未稳定、目标冲突或传输失败时退出码为 1；成功时 stdout
只包含最终绝对路径，日志和错误写 stderr：

```bash
if ! input_file="$(
    nwpc_data_client fetch \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --start-time=2026062400 \
        --forecast-time=24h \
        --target-dir=/work/input
)"; then
    printf '%s\n' "fetch failed" >&2
    exit 1
fi
```

不同候选来源的 basename 可能不同，后续程序应使用 `input_file`，不要自行拼接目标
文件名。需要固定或按起报时次命名的文件名时，优先使用
[`--target-name`](fetch/index.md#自定义目标文件名)，而不是 fetch 之后再改名。
外部信号取消的退出码为 3，可与普通运行失败区分。

## 验证 checker 任务

部署前按需要先执行静态 lint，再执行带起报时次的 dry-run：

```bash
nwpc_data_client check lint ./checker.yaml
```

```bash
nwpc_data_client check local 2026062400 \
    --checker-config=./checker.yaml \
    --dry-run
```

`check lint` 不加载数据配置；dry-run 会展开候选路径和命令模板，但不检查文件是否
存在、不执行命令，也不写运行报告。两者的具体边界见
[校验、错误与运行报告](check-local/validation-and-reports.md)。

业务运行建议设置 `report_file`，用最终 JSON 报告判断逐时效状态。普通检查或命令
失败退出码为 1，可处理的信号取消为 3；报告写失败也会使 CLI 非零退出。

## 找不到文件或得到 NOTFOUND

按以下顺序检查：

1. 用 `--debug` 确认实际展开并检查了哪些候选路径；
2. 核对起报时次、预报时效、集合成员和 `--location-level`；
3. 使用 `local --show-types` 确认当前内嵌或自定义目录中确实存在该数据类型；
4. 对 CMADAAS 等挂载来源，确认 `--storage-base` 指向已挂载且可读的绝对路径；
5. 直接检查候选目录，区分文件尚未生成、挂载不可用和路径模板不符合当前业务布局。

挂载不存在或不可读在 `local` 中可能表现为最终 `NOTFOUND`，不一定是产品尚未生成。
CMADAAS 的完整配置方法见[CMADAAS 挂载归档](cmadaas.md)。

## 自定义目录中提示数据类型缺失

`--data-config-dir` 是配置树的**全量替换**，缺失文件不会回退到内嵌配置。先导出
完整基线再修改：

```bash
nwpc_data_client config export ./myconf
```

使用时传入 `./myconf/local`，而不是额外的 `local/local` 子目录。详见
[使用自定义配置](data-config/custom-configs.md)。

## StorageBase 或 cmafs 配置与旧 client 不兼容

包含 `{{.StorageBase}}` 或 cmafs 候选的配置需要支持这些能力的 client。应先升级
目标机器的安装根，再发布配置；不要只替换配置文件。版本兼容要求分别见
[CMADAAS 挂载归档](cmadaas.md)和[数据配置格式](../reference/data-config-format.md)。

## fetch 报告源为空或尚未稳定

源文件必须是非空普通文件，并在稳定间隔内保持大小和 mtime 不变。正在写入的产品
不会触发复制，也不会自动改选其他候选。

如果本地候选需要等待产品生成，使用 `check local` 等待并确认稳定后再执行后续命令；
`fetch` 本身是一次性操作，`check local` 不轮询 cmafs。稳定检查不是内容格式校验，
详细语义见[源文件检查与目标发布](fetch/reliability.md)。

## fetch 报告目标冲突

fetch 不覆盖已有目标。目标类型、源路径或大小/mtime 不符合复用规则时，应确认已有
文件是否属于同一次业务结果，或改用新的目标目录，不要假定重复执行会强制覆盖。

同目录锁文件会保留，不能删除锁文件后推断没有其他等待者。异常退出留下的唯一命名
临时文件应按业务运行目录的清理策略处理。完整规则见
[重复执行与已有目标](fetch/reliability.md)。

## cmafs 连接、认证或主机密钥失败

- 确认运行节点安装了 `ssh` 和 `rsync`，现有密钥或 agent 支持非交互登录；
- 提供可信的 `known_hosts_file`，不要通过关闭主机密钥校验绕过错误；
- 检查 backend 的主机、端口、用户和超时字段；
- 注意 `config lint --backend-config` 只做离线校验，不能证明网络和认证可用；
- 权限、认证、主机密钥、网络和传输错误会直接失败，不会继续其他候选。

完整配置和获取流程见[从 cmafs 二级存储获取](fetch/cmafs.md)。

## check local 等待或命令超时

- 等待超时：核对 `max_check_count × check_interval` 是否覆盖产品最晚到达时间，
  并用 `--debug` 和运行报告定位具体时效；
- 批量任务：确认 `check_mode` 和 `arrival_order` 符合数据实际到达方式；
- 命令超时：调整 `command_timeout`，`0` 表示禁用；
- 信号取消或超时后仍有子进程：命令取消不承诺清理 shell 派生的整个子孙进程树，
  长驻命令需要单独验证退出行为。

轮询调优见[轮询模式与到达顺序](check-local/polling-modes.md)，取消、超时和报告
状态见[校验、错误与运行报告](check-local/validation-and-reports.md)。
