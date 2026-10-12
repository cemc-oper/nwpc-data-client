# 测试

## 单元测试

```bash
make test        # 即 go test ./...
```

## 集成测试

位于 `tests/integration/`，带 `integration` 构建标签，通过 `os/exec`
调用 `bin/nwpc_data_client` 检查 CMA HPC 上的真实业务文件；数据不存在时
自动 skip。`go test ./...` 不会运行它们。

```bash
make test-integration   # 自动先构建二进制
```

辅助逻辑识别环境变量 `NWPC_DATA_CLIENT_PROGRAM` 与
`NWPC_DATA_CLIENT_CONFIG_DIR`，缺省为 `./bin/nwpc_data_client` 与 `./conf`。

## 端到端测试

```bash
make test-e2e
```

> 本文档正在编写中：补充测试编写约定。
