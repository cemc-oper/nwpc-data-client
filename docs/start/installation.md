# 安装

## 下载预编译二进制

从 GitHub Releases 下载对应平台的二进制（支持 linux / windows 的 amd64 与 arm64）。

## 源码构建

需要 Go 工具链。在仓库根目录（`repo/nwpc-data-client/`）执行：

```bash
make
```

二进制输出到 `./bin/nwpc_data_client`。

也可以使用 GoReleaser 做本地快照构建：

```bash
goreleaser build --snapshot --clean --single-target
```

> 本文档正在编写中：补充各平台安装细节与校验方法。
