# 安装

## 下载预编译二进制

从 [GitHub Releases](https://github.com/cemc-oper/nwpc-data-client/releases)
下载对应平台的 `nwpc_data_client`，支持 linux / windows 的 amd64 与 arm64。
解压后放入 `PATH` 即可使用。

## 源码构建

需要 Go 工具链（版本见 `go.mod`）。在仓库根目录
（`repo/nwpc-data-client/`）执行：

```bash
make
```

二进制输出到 `./bin/nwpc_data_client`，版本号、构建时间、git commit
通过 `-ldflags` 注入。无 git 的环境（如源码包、隔离 HPC）同样可用
`make` 构建，版本号会带 `-dev` 后缀。

本地快照构建（需要安装 GoReleaser）：

```bash
goreleaser build --snapshot --clean --single-target   # 本机平台
goreleaser release --snapshot --clean                 # 全平台交叉编译，产物在 dist/
```

## 验证

```bash
nwpc_data_client version
```

输出版本号、构建时间与 git commit 即安装成功。

## 内嵌配置

二进制内置了 CEMC 主要业务系统的数据配置（`conf/local/` 目录树编译
嵌入），开箱即用；需要自定义时见[指南：使用自定义配置](../guide/custom-configs.md)。
