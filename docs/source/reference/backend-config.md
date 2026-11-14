# 机器后端配置

`fetch --backend-config` 指定机器级后端 YAML。它描述 cmafs 二级存储的 SSH/rsync
连接方式，与 data config 中的数据位置分离。只有实际需要 cmafs 候选时才要求后端；
纯本地 fetch 不需要该文件。完整的候选配置、认证准备和获取步骤见
[指南：从 cmafs 二级存储获取](../guide/fetch/cmafs.md)。

配置解析是离线的。可用以下命令单独校验，不会连接服务器：

```bash
nwpc_data_client config lint --backend-config=./backend.yaml
```

## 示例

SC1 登录节点使用的配置：

```yaml
version: 1

cmafs:
  host: 10.40.150.126
  port: 22
  user: ""                    # 空值使用当前用户
  ssh_path: ssh
  rsync_path: rsync
  identity_file: ""           # 默认使用已有 SSH 配置和密钥
  known_hosts_file: /g1/u/wangdp/.ssh/known_hosts
  connect_timeout: 15s
  transfer_timeout: 2h
```

`version` 是 backend 自身格式版本，当前必须为 1；它不是 data config 的
`_config_version`。配置可以没有 `cmafs` 段，表示未配置远端后端。

## 字段

| 字段 | 默认值 | 说明 |
|---|---|---|
| `cmafs.host` | 无，cmafs 段存在时必选 | 二级存储 SSH 主机名或地址 |
| `cmafs.port` | `22` | SSH 端口，范围 1–65535 |
| `cmafs.user` | 空 | 远端用户；空为执行 client 的当前用户 |
| `cmafs.ssh_path` | `ssh` | 本机 SSH 可执行文件 |
| `cmafs.rsync_path` | `rsync` | 本机 rsync 可执行文件 |
| `cmafs.identity_file` | 空 | 可选私钥路径，非空时必须为绝对路径 |
| `cmafs.known_hosts_file` | 空 | 可选主机密钥文件，非空时必须为绝对路径 |
| `cmafs.connect_timeout` | `15s` | 单次 SSH Stat 的连接/执行上限 |
| `cmafs.transfer_timeout` | `2h` | 单文件 rsync 传输总上限 |

两个 timeout 都必须是大于零的 Go duration。YAML 使用严格字段校验，未知字段、
重复字段、多文档和不支持的 version 都会报错。host、user 和可执行文件字段不能
包含空白、控制字符或用于注入额外选项的前导 `-`。

## 认证与主机密钥

- 不设置密码字段，不保存私钥内容，也不在命令行传密码。
- 使用 BatchMode 和现有 SSH 密钥/agent，不进入交互密码提示。
- 强制 `StrictHostKeyChecking=yes`，不自动接受未知主机。
- SC1 的系统 SSH 配置把 `UserKnownHostsFile` 设为 `/dev/null`，因此必须像示例一样
  显式设置 `known_hosts_file`，不能依赖空值。
- 其他连接定制优先使用用户的 `~/.ssh/config`；第一版不提供任意 SSH 参数透传。

运行节点还需预装 `ssh` 和 `rsync`，Go 二进制不会捆绑这两个程序。

## 获取流程

cmafs 源使用远端 GNU stat 获取普通文件类型、大小和 mtime。只有明确的不存在才继续
下一个候选；权限、认证、主机密钥、网络和超时均直接失败。

下载使用 `rsync --protect-args -t` 到目标目录内的 client 临时文件。client 检查本地
大小和 mtime，再次读取远端元数据确认传输期间没有变化，最后无覆盖发布。显式
`--method=symlink` 不能用于 cmafs；auto 对 cmafs 始终复制。

SC1 上，登录节点可使用此后端；datatrans 节点已有 `/sstorage` 本地 NFS 挂载，应优先
配置 `type: local, level: storage` 直接复制，不需要 backend。serial 上虽可加载 cmafs
module，但实测底层 SSH 不可用，不属于当前支持位置。
