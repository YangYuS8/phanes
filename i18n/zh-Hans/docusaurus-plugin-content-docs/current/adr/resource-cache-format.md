# ADR 0004：资源缓存格式

## 状态

Proposed

## 背景

Phanes 不得要求 GC-Resources 或任何远程资源仓库。Runtime 也不得下载资源。因此 Builder 必须生成一个标准化的本地缓存，供 runtime 在无网络情况下消费。

## 决策

资源缓存是生成的、可重建的、以读为主的目录：

```text
cache-root/
  manifest.pb
  manifest.json      # 从 manifest.pb 生成的可选检查视图
  cache.sqlite
  blobs/
    sha256/
      ab/
        <digest>
  logs/
    build.log
```

必需 build mode：

- `embedded-minimal`：用于契约和 smoke test 的极小测试数据集。
- `local-cache`：从用户提供的本地输入生成的主路径。

允许 source kind：

- `local-install`：用户选择的本地安装目录。
- `local-archive`：用户选择的本地归档。
- `embedded-minimal`：安全的内置测试 fixture。
- `external-import`：用户选择的本地目录或归档，由兼容 adapter 处理；不得远程下载，不得成为 GC-Resources 依赖。

`manifest.pb` 是权威数据。`manifest.json` 如存在，只是生成的人类可读镜像。

Builder 必须避免半成品缓存被 runtime 接受。初始契约要求最终缓存根目录包含 `.cache-complete`，且 runtime 拒绝存在 `.build-in-progress` 或缺少 `.cache-complete` 的缓存。

## 影响

- Runtime 启动保持本地、离线和确定性。
- 缓存可删除、可重建，不影响存档数据。
- 通过 manifest schema version 做兼容性门控。
- 详细资源索引可以后续在 manifest 契约后增量加入。

## 未决问题

- `embedded-minimal` 的最小 fixture 内容，必须不包含版权资源。
- `.cache-complete` marker 的规范内容。

## 相关

- [Builder/cache 契约](../contracts/builder-cache)
- [SQLite schema 契约](../contracts/sqlite-schema)
