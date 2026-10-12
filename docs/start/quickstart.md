# 快速开始

以下示例使用内嵌配置查找 CMA-GFS 模式区域的 grib2 文件：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h
```

文件存在时输出其绝对路径；未找到时输出配置中的 `default` 值（通常为 `NOTFOUND`）。

查看全部可用的内嵌数据类型：

```bash
nwpc_data_client local --show-types
```

> 本文档正在编写中：补充 check local 的入门示例。
