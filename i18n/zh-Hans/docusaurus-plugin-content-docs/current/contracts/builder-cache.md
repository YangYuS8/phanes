# Builder 与资源缓存契约

Builder 在 runtime 启动前准备本地资源缓存。Runtime 只消费已完成的缓存输出。

## 输入契约

允许的 source kind：

```text
embedded-minimal
local-install
local-archive
external-import
```

规则：

- 输入是本地路径或内置测试 fixture。
- 初始契约不包含远程 URL。
- GC-Resources 绝不是必需依赖。
- `external-import` 是显式选择的可选兼容行为。
- Builder 不得再分发完整版权资源集。

## 输出布局

```text
cache-root/
  manifest.pb
  manifest.json
  cache.sqlite
  blobs/
    sha256/
      ab/
        <digest>
  logs/
    build.log
```

`manifest.pb` 是权威数据。`manifest.json` 可选，仅用于人工检查。

## Atomicity

Builder 必须避免半成品缓存被当成完整缓存：

- 写入时存在 `.build-in-progress`
- 验证 manifest/artifacts/indexes
- 移除 `.build-in-progress`
- 写入 `.cache-complete`

Runtime 必须拒绝存在 `.build-in-progress` 或缺少 `.cache-complete` 的缓存。

## Runtime 接受条件

Runtime 只在以下条件满足时接受缓存：

- manifest 存在
- manifest schema version 受支持
- 必需 artifact 存在
- `cache.sqlite` 可打开
- 选择的 digest 检查通过
- `.cache-complete` 存在
- `.build-in-progress` 不存在

无效缓存行为：

- 启动失败或进入文档化 failed state
- 诊断提示运行 `phanes builder build` 或 `phanes builder verify`
- 不下载资源
- 不自动调用 builder

## 生命周期

- Cache 是生成的、可重建的。
- Cache clean/rebuild 不得触碰 save data。
- Cache 兼容性通过 manifest schema version 门控。
- Cache migration 可优先选择 rebuild。
