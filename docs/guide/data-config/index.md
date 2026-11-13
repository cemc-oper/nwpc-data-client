# 数据配置

数据配置用 YAML 描述数据文件的目录、命名规则和候选位置。`local` 与 `check local` 共用同一套配置解析规则，数据类型则用于定位具体配置文件。

```{eval-rst}
.. toctree::
   :maxdepth: 1

   data-types
   custom-configs
   writing-configs
   config-versions
```

- [数据类型](data-types.md)：理解层级数据类型与配置文件之间的映射。
- [使用自定义配置](custom-configs.md)：选择配置来源，导出并部署定制配置目录。
- [编写数据配置](writing-configs.md)：学习路径、文件名和 Go 模板的写法。
- [数据配置格式版本与迁移](config-versions.md)：了解 v1/v2/v3 的差异以及校验、迁移方法。

完整字段约束见[参考：数据配置格式](../../reference/data-config-format.md)，模板语法见[参考：模板变量与函数](../../reference/template-variables.md)。
