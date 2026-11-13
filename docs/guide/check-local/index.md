# check local 命令

`nwpc_data_client check local` 针对一组预报时效轮询本地数据路径，在文件生成且大小稳定后执行对应命令，适合连接模式产品与后处理、绘图、传输等业务任务。

本章介绍从编写任务到部署运行的完整流程：

```{eval-rst}
.. toctree::
   :maxdepth: 1

   basic-usage
   polling-modes
   validation-and-reports
```

- [基本用法与命令模板](basic-usage.md)：时效输入、checker 配置、命令模板和业务示例。
- [轮询模式与到达顺序](polling-modes.md)：逐时效轮询、批量轮询和顺序到达门控。
- [校验、错误与运行报告](validation-and-reports.md)：部署前验证、取消与超时，以及 JSON 报告。

所有命令行参数见[参考：check local](../../reference/cli/nwpc_data_client_check_local.md)，checker YAML 字段见[参考：Checker 配置](../../reference/checker-config.md)。
