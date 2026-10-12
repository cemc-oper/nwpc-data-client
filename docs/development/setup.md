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

文档相关：

```bash
make docs-cli   # 重新生成 CLI 参考页（修改 cmd/ 下命令定义后必须执行）
make docs       # 构建 Sphinx 文档（需 pip install -r docs/requirements.txt）
```

Windows 环境请使用 Git Bash 或兼容的 make；HPC 等无 git 环境下
`make` 仍可构建（版本号带 `-dev` 后缀），但 GoReleaser 快照构建依赖
git 标签，不可用。
