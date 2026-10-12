# 数据类型

数据类型是**层级相对路径**，而不是扁平名称。内嵌配置按
`<系统>/<版本>/<格式>/<名称>.yaml` 组织在 `conf/local/` 下，例如：

```
conf/local/cma_gfs_gmf/current/grib2/orig.yaml
```

对应的数据类型为：

```
cma_gfs_gmf/current/grib2/orig
```

使用 `--data-type` 单独指定时，加载名为 `local/<data-type>` 的内嵌配置；
与 `--data-config-dir` 联用时加载 `<data-config-dir>/<data-type>.yaml`。

> 本文档正在编写中：补充各业务系统的数据类型说明与目录约定。
