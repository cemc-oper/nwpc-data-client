# 数据配置格式版本与迁移指南

nwpc-data-client 的数据配置文件历史上存在三种模板格式。三者的 YAML 结构完全相同（`default` / `file_name` / `file_names` / `paths`），区别只在模板语法。**配置格式版本与 client 发布版本解耦**，采用单调递增整数（1、2、3）。

| 配置格式版本 | 对应 client 版本 | 模板范围 | 定界符 | 变量名示例 |
|---|---|---|---|---|
| v1 | ≤ 0.8 | 仅 `file_name` / `file_names` / `paths[].path` 字段级解析 | `{...}` | `{.Year}`、`{.Forecast}` |
| v2 | 0.9 | 整个 YAML 文本先模板后解析 | `{...}` | `{.Year}`、`{.ForecastHour}` |
| v3 | ≥ 0.10 | 整个 YAML 文本先模板后解析 | `{{...}}` | `{{.Year}}`、`{{.ForecastHour}}` |

当前二进制只接受 v3 格式（无模板标记的纯文本配置除外）。用错版本的配置会在加载时直接报错，并提示检测到的版本、要求的版本和迁移方法。

## 版本声明

每个配置文件顶部用保留字段 `_config_version` 声明格式版本（必须是静态文本，不允许含模板动作）：

```yaml
# nwpc-data-client data config
_config_version: 3

default: NOTFOUND
file_name: gmf.gra.{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.ForecastHour}}.grb2
```

没有该字段的老文件按启发式规则探测：含 `{{` 判为 v3；含 `{.Forecast}` 或 4DV/1HR 偏移变量判为 v1；含 `{.变量}` / `{函数 ...}` 判为 v2。显式字段与嗅探结果冲突时报错。

## 校验：`config lint`

```bash
# 递归校验目录下所有 *.yaml
nwpc_data_client config lint conf/local

# 校验单个文件
nwpc_data_client config lint --data-config-file=./my_config.yaml
```

lint 按检测到的格式版本选择对应的定界符、函数表和变量结构体，用真实的 Go template 解析器执行模板校验，再对渲染结果做结构校验。CI 中对 `conf/` 全量执行。

## 迁移：`config migrate`

```bash
# 原地迁移单个文件到 v3（默认目标版本）
nwpc_data_client config migrate --data-config-file=./my_config.yaml

# 迁移整个目录树到另一个目录（保持目录结构，源文件不动）
nwpc_data_client config migrate ./old_conf/local --output-dir=./new_conf/local

# 只迁移到 v2
nwpc_data_client config migrate --data-config-file=./my_config.yaml --to-version=2
```

行为要点：

- 迁移前先用 lint 校验源文件。**源文件本身不合法时拒绝迁移并给出明确错误**（例如 v1 配置在 `paths[].path` 中使用了函数——v1 的路径模板没有注册函数，该配置在旧 client 下本就不可用）。
- 迁移结果自动写入 `_config_version` 字段，并通过 lint 验证后才会写出；验证失败不产出文件。
- 已在目标版本的文件自动跳过。
- v2 → v3 只重写以变量（`{.X}` / `{$x}`）、已知函数或模板关键字（`if`/`range`/`end` 等）开头的 `{...}` 动作，路径中的合法花括号（如 shell brace 展开 `{1,2}`）不受影响。

### v1 → v2 变量改名映射

| v1 | v2 |
|---|---|
| `{.Forecast}` | `{.ForecastHour}` |
| `{.Year4DV}`（Month/Day/Hour 同理） | `{generateStartTime .StartTime -3 \| getYear}` |
| `{.Year1HR}`（Month/Day/Hour 同理） | `{generateStartTime .StartTime -1 \| getYear}` |
| `{.Forecast1HR}` | `{generateForecastTime .ForecastTime "1h" \| getForecastHour \| printf "%03d"}` |

v1 与 v2 共有的变量（`{.Year}` 等）和 8 个模板函数不变。注意 v1 的函数只在 `file_name` / `file_names` 中注册；迁移把 4DV/1HR 偏移变量展开为函数调用后，这些调用可以出现在 v2 配置的任何位置（v2 起为整文件模板）。
