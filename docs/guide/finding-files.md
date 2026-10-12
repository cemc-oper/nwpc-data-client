# 查找本地数据文件

`nwpc_data_client local` 按配置模板计算文件路径，并在多个候选目录中
返回第一个存在的文件。

必选参数：

- `--data-type`：数据类型（除非使用 `--data-config-file`）
- `--start-time`：起报时次，`YYYYMMDDHH`
- `--forecast-time`：预报时效，如 `24h`

可选参数：

- `--location-level`：只搜索指定级别的路径条目（如 `runtime`、`runtime,archive`）
- `--member`：集合成员编号，如 `000`、`014`
- `--show-types`：列出可用数据类型并退出
- `--debug`：输出调试日志

> 本文档正在编写中：补充多目录搜索顺序、default 返回值、典型示例。
