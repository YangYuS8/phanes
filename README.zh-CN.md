# Phanes

English version: [README.md](README.md).

Phanes 是一个本地优先的离线 runtime、资源缓存构建器与启动器框架，用于沙盒研究。

目标流程：

```text
Launcher 检查本地配置和资源缓存
  -> 启动仅监听 127.0.0.1 的本地 runtime
  -> runtime 读取本地存档和本地资源缓存
  -> 游戏退出后 launcher 干净关闭 runtime
```

Phanes **不是** Grasscutter fork。Grasscutter、Cultivation 和 GC-Resources 只作为参考材料；Phanes 必须保持自己的架构与契约。

## 不可变边界

- Runtime 默认使用 `127.0.0.1`，不得绑定到 `0.0.0.0`。
- Runtime 启动只使用本地文件，不得下载资源。
- 不得要求 GC-Resources。
- 不得捆绑或再分发完整的版权游戏资源。
- 不得实现或记录 client patching、anti-cheat bypass、protection bypass、官方服务认证绕过、未经授权的在线交互或商业服务替代。
- 存档数据与资源缓存必须分离：存档数据属于用户且应持久保存；资源缓存是生成的、可重建的。

## Contract-first 设计

实现前，跨模块行为先通过版本化契约定义：

- ADR：[`docs/adr/`](docs/adr/)
- 模块边界：[`docs/contracts/module-boundaries.md`](docs/contracts/module-boundaries.md)
- Protobuf 契约：[`docs/contracts/protobuf.md`](docs/contracts/protobuf.md)
- Go interface 契约：[`docs/contracts/go-interfaces.md`](docs/contracts/go-interfaces.md)
- SQLite schema 契约：[`docs/contracts/sqlite-schema.md`](docs/contracts/sqlite-schema.md)
- CLI 契约：[`docs/contracts/cli.md`](docs/contracts/cli.md)
- Runtime HTTP/status API：[`docs/contracts/runtime-http-api.md`](docs/contracts/runtime-http-api.md)
- Builder/cache 契约：[`docs/contracts/builder-cache.md`](docs/contracts/builder-cache.md)
- 合规检查清单：[`docs/contracts/compliance-checklist.md`](docs/contracts/compliance-checklist.md)

偏好技术栈：

- Runtime、builder 和 CLI 使用 Go
- 跨模块消息使用 Protocol Buffers
- 存档数据和资源缓存索引使用 SQLite
- Launcher 使用 Tauri 2

## 当前状态

项目处于契约设计阶段。应先补齐相关 ADR、示例/fixtures 和验证检查，再开始实现。

## 文档站

项目文档使用 Docusaurus 管理，并通过 GitHub Actions 发布到 GitHub Pages。`docs/` 下的英文契约文档是文档站源内容；中文本地化文档位于 `i18n/zh-Hans/docusaurus-plugin-content-docs/current/`。

仓库设置要求：需要先启用 **Settings → Pages → Source → GitHub Actions**，`pages` workflow 才能部署。

本地文档命令：

```bash
npm install
npm run start
npm run build
npm run build:en
npm run build:zh
```

## 发布

带 tag 的 GitHub Actions release workflow 会优先构建 Windows 发布产物，包括 `windows/amd64` 和 `windows/arm64`，打包为 `.zip`，生成 `SHA256SUMS.txt`，并创建 GitHub Release。

## 许可证

Phanes 使用 Apache License, Version 2.0。详见 [`LICENSE`](LICENSE) 和 [`NOTICE`](NOTICE)。
