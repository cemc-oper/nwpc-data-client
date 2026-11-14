# fetch 命令：将本地或二级存储数据落位到目标目录

`nwpc_data_client fetch` 使用与 `local` 相同的数据配置和候选顺序，选中第一个
已经存在的普通文件，再将本地或 cmafs 二级存储文件复制到用户指定的目录；本地源
还可以使用符号链接。成功时 stdout 只输出最终绝对路径，便于直接传给后续程序。

它是一次性操作，**不会等待尚未生成的文件**；需要轮询等待时仍使用
`check local`。

## 基本用法

```bash
nwpc_data_client fetch \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time=2026062400 \
    --forecast-time=24h \
    --target-dir=/work/input
```

成功后输出例如：

```text
/work/input/gmf.gra.2026062400024.grb2
```

目标文件名继承源文件 basename。目标目录不存在时会自动创建。可以将结果保存为
Shell 变量：

```bash
input_file="$(
    nwpc_data_client fetch \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --start-time=2026062400 \
        --forecast-time=24h \
        --target-dir=/work/input
)"
```

日志和错误写 stderr，不会混入上述 stdout 路径。

## copy、symlink 与 auto

| `--method` | 行为 |
|---|---|
| `copy`（默认） | 复制为独立普通文件，保留基本权限和修改时间 |
| `symlink` | 创建指向本地源绝对路径的符号链接 |
| `auto` | runtime/archive 本地源使用 symlink；storage 本地源和 cmafs 使用 copy |

storage 使用 copy 是为了避免生成仅在特定存储节点可见的链接。显式
`--method=symlink` 时，调用方需要自行确认下游能够访问源路径，并接受源文件被
清理或修改的影响。第一版不提供硬链接模式。

## 从 cmafs 二级存储获取

数据配置可以同时声明直接挂载和 cmafs 后备路径，例如 GFS_GMF ORIG：

```yaml
paths:
  - type: local
    level: storage
    path: /sstorage/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
  - type: cmafs
    level: storage
    path: /sstorage/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
```

在 datatrans 上，第一个 local 候选直接命中 NFS 文件；在登录节点没有该挂载时，
候选明确不存在才回退到 cmafs。此时提供机器 backend：

```bash
nwpc_data_client fetch \
    --data-config-file=./my_config.yaml \
    --backend-config=./backend.yaml \
    --start-time=2026111300 \
    --forecast-time=0h \
    --target-dir=/work/input
```

backend 使用系统 SSH/rsync，不需要 `module load cmafs/cmafs`。只有远端明确不存在
才继续其他候选；权限、认证、主机密钥、网络和传输错误都会直接失败。cmafs 不支持
显式 symlink，auto 始终下载为普通文件。

下载前后都会核对远端大小和 mtime，本地临时文件也要通过相同元数据校验，然后才
无覆盖发布。字段、认证要求及 SC1 的 known_hosts 特殊配置见
[机器后端配置](../../reference/backend-config.md)。

## 源文件就绪检查

本地源文件必须：

1. 是普通文件且大小非零；
2. 在 `--stable-interval`（默认 `1s`）间隔的两次检查中，大小和修改时间不变；
3. copy 完成后仍保持相同大小和修改时间。

复制使用同一个已打开源文件描述符完成检查和读取，避免检查后被替换成另一个
文件。两次稳定观察只是就绪启发式，不进行 GRIB、NetCDF 等内容格式校验。

候选路径只在明确不存在时继续下一个；权限错误、目录、空文件或正在变化的文件
都会直接失败，不会改选其他数据源。

## 原子发布、重复执行与并发

复制流程为：

```text
目标目录内唯一临时文件 → 复制与校验 → 无覆盖发布 → 删除临时名称
```

正式目标不会暴露仍在写入的半份文件。同一目标使用同目录锁文件进行跨进程互斥；
锁文件会保留，不能直接删除后假设没有其他等待者。进程异常退出可能留下唯一命名的
临时文件，清理由业务运行目录的维护策略负责。

已有目标按以下规则处理：

- copy 可复用大小非零且大小、mtime（Unix 秒）与源一致的普通文件；
- symlink 可复用指向完全相同源路径的符号链接；
- auto 可以复用合法普通文件；只有非 storage 的本地源可以复用同向 symlink，cmafs 不复用远端链接；
- 方法与目标类型不符、内容元数据不符、空文件或目录均报目标冲突；
- 任何情况下都不覆盖已有目标。

复用是操作幂等判断，不是密码学完整性校验。由于需要重新读取源元数据，源文件已经
删除时不能仅凭目标存在就判定为同一结果。

## 参数与退出码

数据配置参数与 `local` 相同：`--data-type`、`--data-config-dir`、
`--data-config-file`、`--location-level`、`--start-time`、`--forecast-time`
和 `--member`。fetch 特有参数：

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--target-dir` | 无，必选 | 最终文件所在目录 |
| `--backend-config` | 空 | 机器后端 YAML；只有实际需要 cmafs 时必需 |
| `--method` | `copy` | `copy`、`symlink` 或 `auto` |
| `--stable-interval` | `1s` | 两次源稳定性检查的间隔 |
| `--debug` | `false` | 输出调试日志到 stderr |

- 成功：退出码 0，stdout 输出最终绝对路径；
- 未找到、配置错误、源未稳定、目标冲突或复制失败：退出码 1；
- 收到 SIGHUP、SIGINT、SIGTERM 或 SIGQUIT：退出码 3，并清理本进程临时文件。

完整参数见[命令行参考](../../reference/cli/nwpc_data_client_fetch.md)。
