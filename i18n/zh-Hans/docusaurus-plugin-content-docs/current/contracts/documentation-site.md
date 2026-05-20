# 文档站契约

Phanes 文档以仓库 Markdown 管理，并通过 Docusaurus 发布到 GitHub Pages。

## Source of truth

- 英文源文档位于 `docs/`。
- 中文本地化文档位于 `i18n/zh-Hans/docusaurus-plugin-content-docs/current/`。
- Docusaurus 使用 `docs/intro.md` 作为英文入口页。
- ADR 位于 `docs/adr/`。
- 契约文档位于 `docs/contracts/`。
- `sidebars.js` 显式定义顺序，确保重要契约易于查阅。

## 本地命令

```bash
npm install
npm run start
npm run start:zh
npm run build
npm run build:en
npm run build:zh
```

`npm run build` 必须在 GitHub Pages 部署前通过。

## GitHub Pages

GitHub Actions 使用 `actions/deploy-pages` 构建并发布 Docusaurus 站点。

仓库设置必须启用 GitHub Pages，并将 Source 设为 **GitHub Actions**。否则 `actions/deploy-pages` 创建 deployment 时会以 404 失败。

站点配置为项目 Pages 路径：

```text
baseUrl: /phanes/
```

如 GitHub Pages URL 或项目路径变化，需要更新 `docusaurus.config.js` 并运行 `npm run build` 验证链接。

## 内容边界

文档站必须保持项目硬边界：

- 不提供 client patching 指南
- 不提供 anti-cheat/protection bypass 指南
- 不提供官方服务认证绕过指南
- 不要求 GC-Resources
- 不捆绑版权资源数据
- runtime 默认保持 localhost-only
