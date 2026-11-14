# 从 cmafs 二级存储获取

`fetch` 可以在本地候选未命中时，通过 SSH/rsync 获取数据配置中的 cmafs 文件。
`local` 和 `check local` 不访问 cmafs，会跳过这类候选。

## 前置条件

运行节点需要：

- 已安装系统 `ssh` 和 `rsync`，Go 二进制不会捆绑这两个程序；
- 能通过现有 SSH 密钥或 agent 非交互登录二级存储主机；
- 已准备可信的主机密钥。client 强制 `StrictHostKeyChecking=yes`，不会自动接受
  未知主机，也不进入交互密码提示；
- 数据配置包含需要获取的 cmafs 候选。

不需要 `module load cmafs/cmafs`。认证、主机密钥和全部字段约束见
[参考：机器后端配置](../../reference/backend-config.md)。

## 配置本地与 cmafs 候选

同一业务路径可以同时声明直接挂载和 cmafs 后备路径，例如：

```yaml
paths:
  - type: local
    level: storage
    path: /sstorage/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
  - type: cmafs
    level: storage
    path: /sstorage/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{{.Year}}{{.Month}}{{.Day}}{{.Hour}}/ORIG
```

有直接 NFS 挂载的节点命中第一个 local 候选，不访问远端；没有该挂载时，只有候选
明确不存在才回退到 cmafs。权限、认证、主机密钥、网络和传输错误都会直接失败，
不会伪装成“未找到”并继续其他候选。

cmafs 的 `path` 必须是 POSIX 绝对路径且不能包含换行。完整配置规则见
[参考：数据配置格式](../../reference/data-config-format.md)。

## 编写并校验机器后端

机器连接参数与数据位置分开配置。例如：

```yaml
# backend.yaml
version: 1

cmafs:
  host: storage.example.org
  port: 22
  user: ""                    # 空值使用当前用户
  ssh_path: ssh
  rsync_path: rsync
  identity_file: ""           # 默认使用已有 SSH 配置和密钥
  known_hosts_file: /absolute/path/to/known_hosts
  connect_timeout: 15s
  transfer_timeout: 2h
```

`known_hosts_file` 非空时必须是绝对路径。部分系统的 SSH 配置可能把
`UserKnownHostsFile` 设为 `/dev/null`，此时必须显式提供可信的主机密钥文件。

分发前可以离线校验 backend YAML：

```bash
nwpc_data_client config lint --backend-config=./backend.yaml
```

该命令只检查配置，不测试网络连接或认证。

## 执行获取

```bash
nwpc_data_client fetch \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --backend-config=./backend.yaml \
    --start-time=2026111300 \
    --forecast-time=0h \
    --target-dir=/work/input
```

调试单个数据配置时，可用 `--data-config-file=./my_config.yaml` 替换
`--data-type`。只有实际回退到 cmafs 候选时才需要可用的 backend；前面的本地候选
命中后，不读取远端数据。

下载前后都会核对远端大小和 mtime，本地临时文件也要通过相同元数据校验，然后才
无覆盖发布。cmafs 不支持显式 `--method=symlink`，`auto` 始终下载为普通文件。
临时文件、目标复用和冲突规则见[源文件检查与目标发布](reliability.md)。

在已有 `/sstorage` 直接挂载的节点上，应优先使用 `type: local, level: storage`
候选，不必配置 backend。当前支持位置和机器差异以
[机器后端配置](../../reference/backend-config.md)中的说明为准。

连接、认证和主机密钥问题的排查见[脚本接入与故障排查](../troubleshooting.md)。
