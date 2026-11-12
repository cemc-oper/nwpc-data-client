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
make docs       # 构建 Sphinx 文档（需先按下文准备好 Python 环境）
```

文档使用 Sphinx + MyST，依赖声明在 `docs/requirements.txt`，
本地、CI 与 Read the Docs 统一用 [uv](https://docs.astral.sh/uv/) 安装
（不使用 lockfile）。首次准备：

```bash
uv venv --python 3.14 .venv   # 与 CI / Read the Docs 的 Python 版本保持一致
uv pip install --python .venv -r docs/requirements.txt
source .venv/bin/activate
```

之后构建：

```bash
make docs   # 先重新生成 CLI 参考页，再构建到 docs/_build/html
```

只改文档、不需要重新生成 CLI 参考页时，可直接调用 Sphinx
（`-W --keep-going` 与 CI 一致，警告视为错误）：

```bash
sphinx-build -W --keep-going -b html docs docs/_build/html
```

预览可执行 `python -m http.server 8000 -d docs/_build/html`，
浏览器打开 `http://localhost:8000`。

Windows 环境请使用 Git Bash 或兼容的 make；HPC 等无 git 环境下
`make` 仍可构建（版本号带 `-dev` 后缀），但 GoReleaser 快照构建依赖
git 标签，不可用。
