# 构建与版本

## make 构建

`make` 调用 `scripts/version.sh` 计算版本字符串，并通过 `-ldflags`
注入 `common.Version`、`common.BuildTime`、`common.GitCommit`：

- 恰好检出标签 `v$(cat VERSION)`：发布版本（如 `0.11.0`）
- 其他提交：`<VERSION>-dev.<提交数>+g<sha>`
- 无 git（源码包）：`<VERSION>-dev`

`VERSION` 文件保存**下一个**发布版本号（无 `v` 前缀、无 `-dev` 后缀）；
发布工作流会校验推送的标签与其一致。

## GoReleaser 快照构建

```bash
goreleaser release --snapshot --clean               # 全平台，产物在 dist/
goreleaser build --snapshot --clean --single-target # 仅本机平台，最快
```

快照版本由 GoReleaser 依据最近 git 标签推导（`<last-tag>-dev+<sha>`），
与 `scripts/version.sh` 相互独立。无 GoReleaser 或无 git 的环境（如隔离
HPC）使用 `make`。

发布构建面向 linux 与 windows 的 amd64、arm64 四个平台
（配置见 `.goreleaser.yaml`）。
