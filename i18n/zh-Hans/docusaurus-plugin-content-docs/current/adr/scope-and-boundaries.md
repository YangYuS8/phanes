# ADR 0001：范围与边界

## 状态

Accepted

## 背景

Phanes 的目标是成为本地优先的离线 runtime、资源缓存构建器与启动器框架，用于沙盒研究。项目应保持轻量、本地和离线体验，同时允许通过本地 runtime 服务简化架构。

Grasscutter、Cultivation 等历史项目可以作为概念参考，但不能作为 Phanes 的实现基础。

## 决策

Phanes 必须遵守以下边界：

- Runtime 默认仅监听字面量 `127.0.0.1`。
- `0.0.0.0` 和公网服务器部署不在范围内。
- Runtime 启动只读取本地存档和本地资源缓存。
- Runtime 启动不得下载或构建资源。
- GC-Resources 仅作参考，不得成为必需依赖。
- 不得捆绑或再分发完整版权游戏资源。
- 本地缓存由用户提供的本地输入生成。
- 不做 client patching、anti-cheat bypass、protection bypass、官方认证绕过、未经授权的在线交互或商业服务替代。
- 外部项目只能用于概念参考；复制源码需要明确许可证审查和维护者批准。

## 影响

- 缓存缺失或无效时，Runtime 可以快速失败并提示用户先运行 builder。
- Builder 与 Runtime 职责分离。
- Launcher 只负责本地进程监督和诊断展示，不修改系统设置或 patch client。
- 网络优先、公网部署和 bypass 类功能都需要维护者明确变更范围。

## 合规检查

- 初始配置只接受字面量 `127.0.0.1`。`0.0.0.0`、公网地址、`localhost`、`::1` 和其他别名都需要先更新验证决策。
- Runtime 契约不得包含远程资源 URL 字段。
- Builder 输入默认接受本地路径，不接受远程仓库。
- 清理缓存不能触碰存档数据。
- Review 时检查是否复制了参考项目代码。

## 相关

- [模块边界](../contracts/module-boundaries)
- [合规检查清单](../contracts/compliance-checklist)
