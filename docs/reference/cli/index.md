# 命令行参考

本目录页面由 `tools/gen-cli-docs` 基于 cobra 命令定义**自动生成**，
与 `--help` 输出保持一致，请勿手工编辑。

修改任何命令定义（`cmd/` 下的 `Use` / `Short` / `Long` / 参数说明）后，
请运行 `make docs-cli` 重新生成并提交本目录的变更；CI 会校验生成结果
与仓库内容一致，不一致则构建失败。

```{eval-rst}
.. toctree::
   :maxdepth: 1

   nwpc_data_client
   nwpc_data_client_local
   nwpc_data_client_check
   nwpc_data_client_check_local
   nwpc_data_client_config
   nwpc_data_client_config_lint
   nwpc_data_client_config_migrate
   nwpc_data_client_config_export
   nwpc_data_client_version
```
