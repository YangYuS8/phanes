---
slug: /
---

# Phanes 文档

Phanes 是一个本地优先的离线 runtime、资源缓存构建器与启动器框架，用于沙盒研究。

本站是仓库内 contract-first 设计文档的中文浏览入口。英文源文档位于 `docs/`，中文本地化文档位于 `i18n/zh-Hans/docusaurus-plugin-content-docs/current/`。

## 从这里开始

- [范围与边界](adr/scope-and-boundaries)
- [v1 路线图](roadmap/v1)
- [模块边界](contracts/module-boundaries)
- [Protobuf 契约](contracts/protobuf)
- [合规检查清单](contracts/compliance-checklist)

## 硬边界

- Runtime 默认使用字面量 `127.0.0.1`。
- Runtime 不下载资源。
- GC-Resources 仅作参考，不是依赖。
- 不做 client patching、anti-cheat bypass、protection bypass 或官方服务认证绕过。
- 存档数据和资源缓存必须保持分离。
