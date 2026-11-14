# CMADAAS 挂载归档

CMADAAS 数据已经挂载到当前机器时，可以直接作为本地文件使用。数据配置中的对应
条目是 `type: local, level: archive`，不使用 CMADAAS API，也不需要 cmafs backend。
`local`、`fetch` 和 `check local` 都支持这种来源。

目前只有六个 `current/grib2/orig` 内嵌配置包含 CMADAAS 候选，具体清单见
[数据类型指南](data-config/data-types.md)。支持范围不包括 `orig.cold`、其他产品或
部署仓库中的版本化业务配置。

## 指定挂载根

挂载根不写死在内嵌配置中。CMADAAS 路径以 `{{.StorageBase}}/DATA/...` 或
`{{.StorageBase}}/ORIG-DATA/...` 开头，调用方通过 `--storage-base` 指定当前机器的
实际挂载根。非空值必须是绝对路径。

以下示例均假设挂载根为 `/CMADAAS`；其他机器可传入 `/mnt/cmadaas` 等实际路径。
`local` 和 `fetch` 使用 `--storage-base=/CMADAAS`，`check local` 除同名 CLI 参数外，
也可在 checker YAML 中设置 `storage_base: /CMADAAS`。

不引用 `{{.StorageBase}}` 的配置不受该参数影响。如果未设置，空字符串会直接参与
模板渲染，例如 CMADAAS 路径可能成为 `/DATA/...`；该路径不会被人为过滤，若实际存在
也会按普通 local 候选参与查找。因此使用这些内嵌配置时，应显式提供适合当前机器的
挂载根。

模板变量的完整说明见[参考：模板变量与函数](../reference/template-variables.md)。

## 候选顺序与文件名

CMADAAS local 候选位于 cmafs 之前，命中后不会访问 cmafs。使用
`--location-level=archive` 只搜索归档 local 条目，跳过 storage 和 cmafs；因此纯
CMADAAS 查询即使未命中，也不会因为机器未配置 cmafs backend 而报后端错误。

原 HPC archive 路径仍排在 CMADAAS 之前，两者都存在时返回原 HPC 文件。路径优先级
高于文件名优先级，准确顺序以所用数据配置为准。

HPC 文件通常采用 `gmf.gra...` 等名称，CMADAAS 文件通常采用
`Z_NAFP_C_BABJ_...` 命名。调用方应直接使用 `local` 或 `fetch` 返回的实际路径，
不能假定最终 basename。

## 查找或等待文件

只查找现有文件：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026111300 --forecast-time=24h \
    --storage-base=/CMADAAS --location-level=archive
```

等待一组时效时，在 checker YAML 中同时设置数据类型和挂载根：

```yaml
# checker.yaml
data_type: cma_gfs_gmf/current/grib2/orig
location_levels: archive
storage_base: /CMADAAS
forecast_times: ["0h", "6h", "24h"]
execute_command: 'echo "ready: {{.FilePath}}"'
```

```bash
nwpc_data_client check local 2026111300 --checker-config=./checker.yaml
```

`local` 的通用规则见[查找本地数据文件](local/index.md)，等待任务的配置见
[check local 基本用法](check-local/basic-usage.md)。

## 落位到目标目录

从 CMADAAS 获取文件使用普通本地源链路，不经过 SSH/rsync：

```bash
nwpc_data_client fetch \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026111300 --forecast-time=24h \
    --storage-base=/CMADAAS --location-level=archive \
    --target-dir=/work/input --method=copy
```

默认 `copy` 生成独立普通文件；显式 `symlink` 创建指向实际挂载根下源文件的链接，
下游必须能访问同一挂载路径。`auto` 对 archive 本地源选择 symlink，因此需要独立副本
时应显式使用 `copy`，不要依赖 auto。

复制、稳定检查、原子发布和目标复用规则与其他本地源相同，见
[fetch 指南](fetch/index.md)和[源文件检查与目标发布](fetch/reliability.md)。

## 集合预报与部署注意事项

GEPS、REPS 需要同时传入三位成员编号，例如 `--member=001`。REPS 配置包含起报当日
和次日两个目录候选，以兼容当前 06/18 时次的归档布局。

运行前必须确认传入的挂载根已挂载且当前用户可读。`local` 沿用普通 Stat 错误处理，
挂载不存在或不可读时可能最终输出 `NOTFOUND`；可使用 `--debug` 查看候选检查过程，
并直接检查挂载目录进行诊断。

引用 `{{.StorageBase}}` 的配置只能由支持该变量的 client 使用；旧 client 会在整文件
模板渲染时报错。部署时必须先升级 client，再发布这些配置。

常见挂载和路径问题见[脚本接入与故障排查](troubleshooting.md)。
