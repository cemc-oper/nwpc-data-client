# 内嵌配置

`conf/` 下的 YAML 配置可以嵌入二进制，生成代码为
`common/config/config.autogen.go`：

```bash
make generate
```

该命令在 `common/config/generate` 下运行 `go generate`，递归遍历 `conf/`，
把每个 `.yaml` 以其无扩展名的斜杠分隔相对路径为键嵌入（如
`local/cma_gfs_gmf/current/grib2/orig`）。

**注意**：新增或修改 `conf/` 下的 YAML 后必须运行 `make generate` 并提交
`config.autogen.go`，否则内嵌配置与 `conf/` 不一致。

> 本文档正在编写中：补充内嵌配置的加载优先级说明。
