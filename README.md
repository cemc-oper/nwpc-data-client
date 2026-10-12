# NWPC Data Client

[![Maturity-Incubating](https://img.shields.io/badge/Maturity-Incubating-56B4E9)](https://github.com/cemc-oper/nwpc-data-client)
[![GitHub Release](https://img.shields.io/github/v/release/cemc-oper/nwpc-data-client)](https://github.com/cemc-oper/nwpc-data-client/releases)
[![CI](https://github.com/cemc-oper/nwpc-data-client/actions/workflows/ci.yml/badge.svg)](https://github.com/cemc-oper/nwpc-data-client/actions/workflows/ci.yml)
[![Documentation Status](https://readthedocs.org/projects/nwpc-data-client/badge/?version=latest)](https://nwpc-data-client.readthedocs.io/)

面向 CMA 超算平台开发的 CEMC/CMA 数值天气预报业务系统数据查找命令行工具。

## 功能

- 在一组目录中查找数值预报业务系统数据文件路径
- 批量监视数据文件是否生成，当文件就绪后执行命令

## 文档

完整文档托管在 Read the Docs：<https://nwpc-data-client.readthedocs.io/>

- **开始**：简介、安装、快速开始
- **指南**：查找数据文件、数据类型、自定义配置、编写配置、等待文件并执行命令、配置格式版本迁移
- **开发**：环境搭建、构建与版本、测试、内嵌配置、CI 与发版
- **参考**：命令行参考（自动生成）、数据配置格式、Checker 配置、模板变量与函数、内置数据类型清单

## 快速开始

从 [Releases](https://github.com/cemc-oper/nwpc-data-client/releases) 下载
对应平台的二进制，或从源码构建（`make`，产物在 `bin/`）。

查找 CMA-GFS 预报系统 2026 年 6 月 22 日 12 时次 24 小时的基础 GRIB2 数据文件：

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2026062212 \
    --forecast-time 24h
```

批量监视数据文件，文件生成且大小稳定后执行命令：

```bash
echo "000h 024h" |
    nwpc_data_client check local 2026062312 \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --location-level runtime \
        --execute-command 'cp {{.FilePath}} /path/to/archive/'
```

## 构建与开发

```bash
make            # 构建到 ./bin
make test       # 单元测试
make generate   # 修改 conf/ 后重新生成内嵌配置
```

版本号规则、交叉编译、集成测试、发版流程等见文档站"开发"分类。

## 许可证

Copyright &copy; 2019-2026 developers at cemc-oper.

`nwpc-data-client` 基于 [The MIT License](https://opensource.org/licenses/MIT) 授权。
