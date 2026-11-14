# 数据配置格式

数据配置（DataConfig）是描述数据文件位置的 YAML 文档，整体作为 Go 模板
渲染后再解析。本页是渲染后结构的完整字段参考，与
`internal/dataconfig.ValidateDataConfig` 及 JSON Schema（`schemas/config/` 下的
`data-config.schema.json`）保持一致。

## 顶层字段

| 字段 | 类型 | 必选 | 说明 |
|---|---|---|---|
| `_config_version` | int | 否 | 保留字段，声明配置**格式**版本（当前为 3）。只用于版本校验，不参与路径解析 |
| `_arrival_order` | string | 否 | 保留字段，声明预报时效的到达顺序：`sequential`（按时效顺序产出）或 `unordered`（缺省）。只用于 `check local` batched 模式的顺序门控，不参与路径解析 |
| `default` | string | 否 | 所有候选路径都未找到文件时的返回值，通常为 `NOTFOUND` |
| `file_name` | string | 否 | 单个候选文件名（模板） |
| `file_names` | list&lt;string&gt; | 否 | 候选文件名（模板）列表 |
| `paths` | list | **是** | 候选目录条目列表，不能为空 |

`file_name` 与 `file_names` 至少应提供一个。两者同时存在时合并为一个
候选列表：**`file_names` 在前，`file_name` 追加在最后**。

## paths 条目

| 字段 | 类型 | 必选 | 说明 |
|---|---|---|---|
| `type` | string | **是** | 路径类型：`local` 为本地路径，`cmafs` 为通过机器 backend 访问的二级存储路径 |
| `level` | string | **是** | 位置级别：`runtime`、`archive`、`storage`，或用 `/` 连接的组合（如 `runtime/archive`） |
| `path` | string | **是** | 目录路径（模板），不含文件名 |

## 搜索语义

`local` 命令按以下规则查找文件：

1. 按 `paths` 列表顺序遍历目录条目；`--location-level` 非空时跳过
   `level` 不匹配的条目。
2. 每个目录依次与每个候选文件名拼接（`file_names` 在前，`file_name` 在后）。
3. 返回第一个实际存在的文件的完整路径。
4. 全部未命中时返回 `default` 的值。

`local` 和 `check local` 只检查 `type: local` 条目，并跳过 cmafs 条目。
`fetch` 同时处理 local 与 cmafs：本地候选明确不存在时才继续；需要 cmafs 时通过
`--backend-config` 指定的机器后端执行只读 Stat 和 rsync 下载。backend 的权限、
认证、网络或传输错误不会伪装成候选未命中。含 cmafs 条目的配置只能部署给支持该
类型的 client；必须先升级安装根的 client，再发布这些配置。fetch 未找到时返回错误，不输出
`default` 作为成功路径。

cmafs 的 `path` 必须是 POSIX 绝对路径且不能包含换行。同一业务路径可以先写
`type: local, level: storage`，再写 `type: cmafs`：有直接 NFS 挂载的节点命中前者，
没有挂载的登录节点回退到后者。机器连接参数见[后端配置](backend-config.md)。

## 校验规则

`config lint` 与 `internal/dataconfig.ValidateDataConfig` 对渲染后的配置执行：

- 严格 YAML 解码：出现以上未列出的字段直接报错；
- `paths` 必须存在且非空；
- 每个 `paths` 条目的 `type` / `level` / `path` 均必填；
- `type` 接受 `local` 或 `cmafs`；`level` 必须匹配
  `runtime` / `archive` / `storage` 或其 `/` 组合；
- cmafs 的 `path` 必须是绝对路径且不得包含换行；
- `_arrival_order` 只接受 `sequential` / `unordered`。

## 到达顺序与顺序门控

`_arrival_order: sequential` 声明该数据的预报时效**按时效顺序产出**
（模式积分产品的固有属性）。在 `check local` 的 batched 模式下，它启用
顺序门控：每一轮只有最早未到的时效参与目录扫描，它的命令执行完毕后
才放开下一个时效；无序数据（缺省 `unordered`）不做门控，全部待等时效
每轮一起检查。

运行时可用 `--arrival-order` / checker YAML 的 `arrival_order` 覆盖数据
配置（双向有效：`sequential→unordered` 是紧急开关，`unordered→sequential`
是试验手段），缺省 `auto` 跟随数据配置。详见
[Checker 配置](checker-config.md)。

per_time 模式不做门控（各时效协程相互独立）。

## 示例

```yaml
# nwpc-data-client data config
_config_version: 3

default: NOTFOUND
file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2
file_names: []
paths:
  - type: local
    level: runtime
    path: /g1/u/CMA_GFS_GMF/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}
  - type: local
    level: archive
    path: /archive/CMA_GFS_GMF/{{.Year}}/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}
  # datatrans 等节点直接挂载业务二级存储时使用。
  - type: local
    level: storage
    path: /sstorage/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
  # 登录节点没有直接挂载时通过 backend 获取同一文件。
  - type: cmafs
    level: storage
    path: /sstorage/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
```
