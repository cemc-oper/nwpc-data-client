# 快速开始

## 查找一个文件

使用内嵌配置查找 CMA-GFS 模式 gmf 区域的 grib2 产品：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h
```

- 文件存在时输出其绝对路径；
- 未找到时输出配置中的 `default` 值（通常为 `NOTFOUND`）。

查看全部可用的内嵌数据类型：

```bash
nwpc_data_client local --show-types
```

加 `--debug` 可以看到每个候选路径的检查过程，排查"为什么没找到"
时很有用。

## 等待文件并执行命令

`check local` 对一组预报时效轮询等待，文件稳定后执行命令：

```bash
echo "0h 3h 6h" | nwpc_data_client check local 2026062400 \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --execute-command="echo ready: {{.FilePath}}"
```

文件大小连续两次检查不再变化即判定稳定，随后执行命令模板，
`{{.FilePath}}` 会被替换为实际文件路径。

## 下一步

- 查找文件的完整参数与搜索规则：[指南：查找本地数据文件](../guide/finding-files.md)
- 等待任务的配置文件写法：[指南：等待文件并执行命令](../guide/checking-files.md)
- 编写自己的数据配置：[指南：编写数据配置](../guide/writing-configs.md)
