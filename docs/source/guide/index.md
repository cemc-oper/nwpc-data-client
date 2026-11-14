# 指南

面向日常使用者的任务式文档，按命令、存储来源和数据配置组织。

## 按任务选择命令

| 需求 | 命令 | 是否等待文件生成 | 是否获取 cmafs | 成功输出 |
|---|---|---|---|---|
| 只取得现有文件路径 | `local` | 否 | 否 | 源文件路径；未找到时输出配置中的 `default` |
| 复制或链接到工作目录 | `fetch` | 否 | 是 | 最终目标文件的绝对路径 |
| 等待多个时效并执行任务 | `check local` | 是 | 否 | 命令输出及可选 JSON 运行报告 |

`local` 未找到文件时仍以退出码 0 返回，脚本需要额外判断 `default` 值。
`fetch` 是一次性落位操作；需要等待产品生成时，使用 `check local`。三者共用
数据类型和数据配置，但对 cmafs 候选的处理不同，详见
[数据配置](data-config/index.md)。

```{eval-rst}
.. toctree::
   :maxdepth: 2

   local/index
   fetch/index
   check-local/index
   data-config/index
   cmadaas
   troubleshooting
```

## 专题指南

- [CMADAAS 挂载归档](cmadaas.md)：在 `local`、`fetch` 和 `check local` 中使用
  已挂载的 CMADAAS 数据。
- [脚本接入与故障排查](troubleshooting.md)：处理输出与退出码，并诊断路径、
  挂载、传输、目标冲突和等待超时问题。
