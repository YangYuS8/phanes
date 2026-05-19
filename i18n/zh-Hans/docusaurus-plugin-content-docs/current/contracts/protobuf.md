# Protobuf 契约

Protobuf 包应小型、版本化，并聚焦跨模块消息。初始命名空间为 `phanes.<area>.v1`。

## 包列表

- `phanes.common.v1`：共享基础类型、版本、digest、诊断。
- `phanes.cache.v1`：生成的本地资源缓存 manifest、source、artifact、index。
- `phanes.builder.v1`：builder 输入、输出、验证和进度。
- `phanes.runtime.v1`：runtime 配置、状态、生命周期、session 和诊断响应。
- `phanes.launcher.v1`：launcher 配置、runtime 进程监督和日志事件。
- `phanes.save.v1`：profile 与 save metadata 身份。

## 当前 proto 文件

```text
proto/phanes/common/v1/common.proto
proto/phanes/cache/v1/cache.proto
proto/phanes/builder/v1/builder.proto
proto/phanes/runtime/v1/runtime.proto
proto/phanes/launcher/v1/launcher.proto
proto/phanes/save/v1/save.proto
```

## 规则

- 不创建单一巨型 `phanes.proto`。
- 核心契约不得使用 `map<string, any>` 风格逃生口。
- 实现消费者前，应先添加示例/fixtures 和验证。
- 通过版本化 package 演进契约，不能静默改变语义。
- `CacheManifest` 不包含自身 digest；manifest 完整性由 `.cache-complete`、`manifest.pb.sha256` 等外部文件表达，以避免自引用 hash。
