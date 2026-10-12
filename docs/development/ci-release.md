# CI 与发版

GitHub Actions 工作流位于 `.github/workflows/`：

- `ci.yml`：push / PR 到 main、master、develop 时运行
  `go vet`、gofmt 检查、`make`、`make test`、e2e 测试（含 Windows 作业）。
- `release.yml`：推送 `v*` 标签时运行 GoReleaser 发布。
  标签必须与 `VERSION` 文件一致，否则发布失败。
- `docs.yml`：校验 CLI 参考页与命令定义同步（`make docs-cli` 后
  `git diff --exit-code`），并以严格模式（`-W`）构建 Sphinx 文档。
  文档本身由 **Read the Docs** 构建与托管（配置见仓库根目录的
  `.readthedocs.yaml`）。

> 本文档正在编写中：补充发版操作步骤。
