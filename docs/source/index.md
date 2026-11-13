# nwpc-data-client 文档

nwpc-data-client 是面向 CEMC（NWPC/CMA）业务模式的数据文件查找工具。
它根据起报时次和预报时效展开 YAML 配置模板得到文件路径，也可以轮询等待
文件生成并执行后续命令。

## 核心特色

::::{grid} 1 1 2 2
:gutter: 3

:::{grid-item-card} 自动解析数据路径
:link: guide/local/index
:link-type: doc

将起报时次、预报时效和集合成员代入配置模板，按优先级检查运行、归档等多个存储位置，返回第一个实际存在的文件。
:::

:::{grid-item-card} 等待文件并触发任务
:link: guide/check-local/index
:link-type: doc

批量等待一组预报时效，确认文件生成且大小稳定后执行模板化命令，轻松串接绘图、传输和后处理任务。
:::

:::{grid-item-card} 配置驱动，开箱即用
:link: guide/data-config/index
:link-type: doc

内置 CEMC 主要业务系统的数据配置，以层级数据类型统一管理；同时支持导出基线、定制单个配置或部署完整配置目录。
:::

:::{grid-item-card} 面向业务运行的工具链
:link: guide/check-local/validation-and-reports
:link-type: doc

提供静态校验、运行前 dry-run、JSON 运行报告、信号取消和命令超时，并支持逐时效与批量两种轮询模式。
:::

::::

```{eval-rst}
.. toctree::
   :maxdepth: 2
   :caption: 开始
   :hidden:

   start/index
```

```{eval-rst}
.. toctree::
   :maxdepth: 2
   :caption: 指南
   :hidden:

   guide/index
```

```{eval-rst}
.. toctree::
   :maxdepth: 2
   :caption: 开发
   :hidden:

   development/index
```

```{eval-rst}
.. toctree::
   :maxdepth: 2
   :caption: 参考
   :hidden:

   reference/index
```
