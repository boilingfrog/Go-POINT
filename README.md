# Go-POINT

[![CI](https://github.com/boilingfrog/Go-POINT/actions/workflows/ci.yml/badge.svg)](https://github.com/boilingfrog/Go-POINT/actions/workflows/ci.yml)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/boilingfrog/Go-POINT?style=flat)](https://github.com/boilingfrog/Go-POINT/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/boilingfrog/Go-POINT?style=flat)](https://github.com/boilingfrog/Go-POINT/forks)

Go-POINT 是一个面向中文 Go 开发者的工程知识库与可运行示例集合，持续整理 Go 语言机制、数据库、分布式系统、容器、网络和工程实践。

A Chinese-language Go engineering knowledge base with runnable examples for Go internals, databases, distributed systems, containers, and networking.

这个仓库既记录原理和源码阅读，也保留可以复现的最小示例，适合用于系统学习、问题排查和技术选型参考。内容以学习与研究为目的，不应未经评估直接用于生产环境。

## 适合谁

- 希望系统理解 Go 语言机制与标准库实现的开发者
- 需要查阅数据库、缓存、消息队列和分布式系统实践的后端工程师
- 希望通过小型可运行示例验证技术原理的学习者
- 愿意一起维护中文 Go 工程资料的贡献者

## 内容导航

| 方向 | 主要内容 |
| --- | --- |
| [Go 语言](golang/) | slice、map、channel、context、并发原语、GC、unsafe 等 |
| [MySQL](mysql/) | 原理、部署、复制、读写分离和故障排查 |
| [MongoDB](mongo/) | 使用、部署和运维实践 |
| [PostgreSQL](pgsql/) | PostgreSQL 学习笔记 |
| [Redis](redis/) | 数据结构、分布式锁和可运行示例 |
| [Elasticsearch](elasticsearch/) | 搜索与存储相关实践 |
| [Docker](docker/) / [Kubernetes](k8s/) | 容器、编排、Helm、Ingress 和探针 |
| [TCP](tcp/) / [gRPC](grpc/) | 网络协议和 RPC |
| [消息队列](mq/) | RabbitMQ 等消息队列实践 |
| [Linux](linux/) / [Git](git/) | 开发与运维基础工具 |

更多内容还包括架构、Bazel、Ansible、认证、事故复盘和 AI 工具实践。完整目录以仓库代码为准。

## 快速开始

环境要求：

- Go 1.18 或更高版本
- 部分集成示例需要 Docker、Redis、MySQL、MongoDB 或 Kubernetes

克隆并运行基础测试：

```bash
git clone https://github.com/boilingfrog/Go-POINT.git
cd Go-POINT
go test ./...
```

Redis 锁测试会在本机 Redis 不可用时自动跳过。若要完整运行：

```bash
docker run --rm -p 6379:6379 redis:7-alpine
go test ./redis/lock -v
```

## 项目原则

- **可验证**：示例尽量保持可编译、可运行，并通过 CI 检查。
- **解释背景**：不仅记录命令，也说明适用场景、原理和限制。
- **最小复现**：复杂问题优先沉淀为独立、可复现的示例。
- **持续维护**：修正过期内容，跟踪 Go 和相关基础设施的演进。
- **安全优先**：提交内容不得包含真实凭证、生产配置或敏感数据。

## 参与贡献

欢迎修正文档、补充示例、报告失效链接或提出新的学习主题。提交前请阅读：

- [贡献指南](CONTRIBUTING.md)
- [安全策略](SECURITY.md)
- [行为准则](CODE_OF_CONDUCT.md)
- [维护路线图](ROADMAP.md)

可以通过 [Issue](https://github.com/boilingfrog/Go-POINT/issues) 报告问题，或直接提交 Pull Request。如果内容对你有帮助，也欢迎 Star 仓库，让更多 Go 开发者发现这些资料。

## 维护者

项目由 [boilingfrog](https://github.com/boilingfrog) 发起并维护。技术文章同步发布于[博客园](https://www.cnblogs.com/ricklz/)。

## 许可证

本项目使用 [GNU General Public License v3.0](LICENSE)。引用或复用内容时，请遵守许可证并保留来源说明。
