# 使用自定义配置

nwpc-data-client 支持三种配置来源，按优先级从高到低：

| 来源 | 参数组合 | 说明 |
|---|---|---|
| 单个配置文件 | `--data-config-file=<file>` | 直接加载该文件，优先级最高 |
| 配置目录 | `--data-config-dir=<dir>` + `--data-type` | 加载 `<dir>/<data-type>.yaml` |
| 内嵌配置 | 仅 `--data-type` | 编译进二进制的配置，默认方式 |

## 配置目录是全量替换

`--data-config-dir` 是**全量替换**而非覆盖叠加：

- 加载 `<dir>/<data-type>.yaml`，不会附加 `local/` 前缀；
- 目录中缺失的配置直接报错，**不会**回退到内嵌配置；
- `--show-types` 也只列出目录中的配置。

因此配置目录必须自包含：部署时用到的每个数据类型都要有对应文件。

## 在内嵌配置基础上定制

需要以内嵌配置为基线做本地扩展时，先导出再修改：

```bash
# 导出内嵌配置为 ./myconf/local/... 目录树（与 conf/ 字节一致）
nwpc_data_client config export ./myconf

# 按需修改、增删配置后，作为配置目录使用
nwpc_data_client local --data-config-dir=./myconf/local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 --forecast-time=24h
```

`config export` 默认保留已存在的文件，加 `--force` 强制覆盖。

## 典型场景

**场景一：单机调试单个配置**

把某个配置文件复制出来改路径模板，用 `--data-config-file` 直接验证，
不影响其他任何配置。

**场景二：为某个业务系统部署定制配置集**

`config export` 导出基线 → 修改路径指向本地存储 → 分发整个目录 →
各节点用 `--data-config-dir` 引用。

**场景三：checker 任务**

`check local` 同样支持以上三种方式，通过 `--checker-config` YAML 中的
`data_config_dir` / `data_config_file` / `data_type` 字段或对应 CLI
参数指定。

## 校验自定义配置

分发前用 `config lint` 校验（支持各格式版本）：

```bash
nwpc_data_client config lint ./myconf/local
nwpc_data_client config lint --data-config-file=./my_config.yaml
```
