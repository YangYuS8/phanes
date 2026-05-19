# ADR 0002：进程模型

## 状态

Proposed

## 背景

Phanes 需要 launcher、runtime 和 builder。过早把它们合并会模糊生命周期职责，也容易让 runtime 启动阶段承担资源准备或 launcher 工作。

## 决策

第一阶段采用分离的逻辑进程：

```text
Launcher 进程
  负责用户流程
  校验本地配置和缓存状态
  启停 Runtime 子进程
  用户明确请求时运行 Builder
  展示日志和诊断

Runtime 进程
  仅绑定 loopback
  暴露本地 HTTP/status/session API
  读取已验证资源缓存
  读写存档数据
  不下载、不构建资源

Builder 进程
  读取用户提供的本地输入
  写入标准化本地缓存
  验证缓存完整性
  不依赖 runtime
```

初始 launcher/runtime 通信采用子进程监督和 loopback HTTP 轮询。暂不引入自定义 IPC、gRPC、插件系统或后台 daemon 管理器。

## 影响

- CLI 可以独立验证 runtime 和 builder。
- 崩溃和日志更容易隔离。
- Runtime 契约保持狭窄：消费 cache + save data，暴露本地 API。
- Launcher 是监督和 UX 表层，不是隐藏的系统设置修改器。

## 未决问题

- Builder 早期由 launcher 作为子进程调用，还是先作为 Go library 由 CLI 包装。
- Launcher 前端使用 Svelte、React 或其他 Tauri 兼容栈。

## 相关

- [Runtime HTTP API](../contracts/runtime-http-api)
- [CLI 契约](../contracts/cli)
