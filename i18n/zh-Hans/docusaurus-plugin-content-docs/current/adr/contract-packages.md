# ADR 0003：契约包

## 状态

Proposed

## 背景

Phanes 应避免未文档化 JSON、无类型 map、大型全局 struct 和跨模块 schema 漂移。Go 与 protobuf 应用于需要稳定类型契约的边界。

## 决策

定义小型、版本化的 protobuf 包：

```text
phanes.common.v1       通用基础类型、诊断、digest、版本
phanes.cache.v1        生成的资源缓存 manifest 与 artifact
phanes.builder.v1      build/verify 请求、结果和进度
phanes.runtime.v1      runtime 配置、状态、生命周期、session
phanes.launcher.v1     launcher 配置、进程状态、日志事件
phanes.save.v1         profile 与 save metadata 身份
```

模块边界使用生成的 Go 类型。JSON 可用于人工可编辑配置、UI 数据、日志、导入/导出 envelope 和进程 ready 文件，但 JSON 形状也必须文档化。

## 影响

- 契约变更显式且可 review。
- Runtime、builder、launcher 和 storage 可以通过版本门控演进。
- 后续可以加入 generated code 而不改变设计词汇。
- 早期契约保持小，不急于建模完整游戏状态。

## 未决问题

- 具体 protobuf 生成工具和 Go module path。
- 生成类型初期是否仅放在 `internal/`，或在有外部 API 需求后从 `pkg/` 导出。

## 相关

- [Protobuf 契约](../contracts/protobuf)
- [Go interface 契约](../contracts/go-interfaces)
