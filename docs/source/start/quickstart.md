# 快速开始

## 查找文件

`local` 使用起报时次和预报时效展开数据配置，依次检查候选路径，并输出第一个实际存在的文件路径。

### CMA-GFS 确定性预报

使用内嵌配置查找 CMA-GFS 模式 gmf 区域的基础 GRIB2 产品：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h
```

### CMA-GEPS 集合预报

查找 CMA-GEPS 的 GRIB2 产品时，还需要通过 `--member` 指定集合成员。例如查找 014 号成员：

```bash
nwpc_data_client local \
    --data-type=cma_geps/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h \
    --member=014
```

成员编号会代入配置中的 `{{.Member}}` 模板变量，参与生成文件名。

- 文件存在时输出其绝对路径，退出码为 0；
- 未找到时输出配置中的 `default` 值（通常为 `NOTFOUND`）。

查看全部可用的内嵌数据类型：

```bash
nwpc_data_client local --show-types
```

加 `--debug` 可以看到每个候选路径的检查过程，排查“为什么没找到”时很有用。

## 将文件复制或链接到指定目录

`fetch` 在 `local` 的查找结果上增加数据落位。下面把选中的文件复制到工作目录，
stdout 返回最终绝对路径：

```bash
nwpc_data_client fetch \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h \
    --target-dir=/work/input
```

默认方法为 `copy`；`--method=symlink` 创建本地源链接，`--method=auto` 对
runtime/archive 使用链接、对 storage 使用复制。源文件需要通过非空和稳定性检查，
已有目标不会被覆盖。fetch 是一次性操作，不会等待尚未生成的数据。

## 等待文件并执行命令

`check local` 可以等待一组预报时效对应的文件，并在每个文件就绪后执行命令。
下面的示例等待 0、3、6 小时三个时效：

```bash
echo "0h 3h 6h" | nwpc_data_client check local 2026062400 \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --execute-command="echo ready: {{.FilePath}}"
```

其中：

- 起报时次 `2026062400` 是位置参数，标准输入提供以空白分隔的预报时效；
- 默认的 `per_time` 模式为每个时效启动独立检查任务，任务之间默认错开 10 秒；
- 文件尚未生成时，默认每 5 秒重试一次，每个时效最多检查 2880 次；
- 文件出现后，只有大小**非零且连续两次检查相同**才判定为稳定，避免在文件仍在写入时触发后续任务；
- 文件稳定后执行命令模板，`{{.FilePath}}` 会被替换为实际文件路径；检查或命令失败时，CLI 以非零码退出。

业务任务推荐使用 `--checker-config` 提供 YAML 配置，将预报时效、轮询参数和命令集中管理。大量时效还可以使用 `batched` 批量轮询模式，减少重复目录访问。

## 下一步

- 查找文件的完整参数与搜索规则：[指南：local 命令](../guide/local/index.md)
- 复制、链接与目标冲突规则：[指南：fetch 命令](../guide/fetch/index.md)
- 等待任务的配置与命令写法：[指南：check local 基本用法](../guide/check-local/basic-usage.md)
- 选择轮询模式：[指南：轮询模式与到达顺序](../guide/check-local/polling-modes.md)
- 编写自己的数据配置：[指南：编写数据配置](../guide/data-config/writing-configs.md)
- 使用 CMADAAS 挂载数据：[指南：CMADAAS 挂载归档](../guide/cmadaas.md)
- 在脚本中判断输出并排查错误：[指南：脚本接入与故障排查](../guide/troubleshooting.md)
