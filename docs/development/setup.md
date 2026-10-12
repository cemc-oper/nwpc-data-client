# 开发环境搭建

仓库采用 workspace 布局：工作区根目录的 `CLAUDE.md` 提供总体说明，
实际的 Git 仓库与 Go module 位于 `repo/nwpc-data-client/`，
所有构建、测试命令都在该目录下执行。

要求：

- Go（版本见 `go.mod`）
- 可选：GoReleaser（快照构建与交叉编译）

常用命令：

```bash
make          # 构建到 ./bin
make test     # 单元测试
go vet ./...  # 静态检查（CI 强制）
gofmt -l .    # 格式检查（CI 强制）
```

> 本文档正在编写中：补充 Windows / HPC 环境下的注意事项。
