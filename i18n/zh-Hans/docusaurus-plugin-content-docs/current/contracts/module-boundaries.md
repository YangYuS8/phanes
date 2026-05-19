# 模块边界

本文定义 Phanes 第一阶段模块边界。实现包结构可以变化，但跨模块职责不应静默漂移。

## Runtime

职责：

- 仅绑定 loopback，默认 `127.0.0.1`。
- 暴露健康、状态、诊断、session 生命周期和 shutdown 的本地 HTTP API。
- 读取已验证资源缓存。
- 读写本地存档数据。
- 输出结构化诊断。

非职责：

- 不下载资源。
- 不构建缓存或检测源。
- 不做公网服务器部署。
- 不做 client patching 或 bypass 行为。

## Builder

职责：

- 检测本地输入源。
- 将用户提供的本地输入标准化为版本化资源缓存。
- 验证缓存完整性。
- 生成 manifest、index、blob 和 build log。

非职责：

- 不管理 runtime 生命周期。
- 不要求 GC-Resources。
- 初始契约不下载远程仓库。
- 不捆绑完整版权资源集。

## Launcher

职责：

- 提供用户流程。
- 校验本地配置和缓存状态。
- 启停 runtime 子进程。
- 用户明确选择时运行 builder。
- 展示日志和诊断。

非职责：

- 不静默修改系统代理。
- 不 patch client。
- 不提供公网服务器控制。
- 不做隐藏的网络优先 daemon。

## Protocol

职责：

- 保存版本化 protobuf 定义。
- 提供 contract message 的生成 Go 类型。
- 定义兼容性和 schema version 规则。

## Storage

职责：

- 分离存档数据和资源缓存。
- 提供显式 migration。
- 保护存档数据不受缓存清理/重建影响。
- 为存档数据提供 backup/export/import 边界。

推荐分离：

```text
save.sqlite   用户所有、持久化、谨慎 migration
cache.sqlite  builder 生成、可重建、以读为主
```
