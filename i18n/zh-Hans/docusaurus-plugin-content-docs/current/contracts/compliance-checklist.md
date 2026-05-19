# 合规检查清单

用于设计 review、实现 review 和 handoff。

## Localhost-only runtime

- [ ] 默认 bind host 是 `127.0.0.1`。
- [ ] `0.0.0.0` 被拒绝。
- [ ] 非 loopback/公网 bind 地址被拒绝。
- [ ] Runtime status 报告实际 bind host 和 port。
- [ ] 未经维护者批准，不引入公网服务器功能。

## Runtime 不下载资源

- [ ] Runtime config 没有远程资源 URL 字段。
- [ ] 缓存缺失/无效时失败并给出诊断，而不是 fetch。
- [ ] Runtime 启动不调用 builder。
- [ ] Runtime package 不引入远程资源 client。

## 不要求 GC-Resources

- [ ] Builder 支持 `embedded-minimal` 和本地输入，无需 GC-Resources。
- [ ] `external-import` 可选且显式。
- [ ] 文档不要求用户下载 GC-Resources。
- [ ] 测试不依赖远程资源仓库。

## 不捆绑版权资源

- [ ] Fixture 是极小、合成或安全可再分发的。
- [ ] Builder 输出来自用户提供的本地输入。
- [ ] 仓库不提交完整版权游戏资源。

## 不做 unsafe bypass

- [ ] 无 anti-cheat bypass。
- [ ] 无 client protection bypass。
- [ ] 无官方认证绕过。
- [ ] 无未经授权的在线交互。
- [ ] 无商业服务替代。
- [ ] 无 deep client patching 指南。
- [ ] Launcher 不静默修改系统代理/设置。

## Clean-room/reference 姿态

- [ ] Grasscutter 和 Cultivation 只作概念参考。
- [ ] 未经许可证审查和维护者批准，不复制 GPL/AGPL/reference 项目源码。
- [ ] 不盲目克隆参考项目包结构和 handler 设计。

## Contract-first workflow

- [ ] 实现前先更新契约。
- [ ] 为契约变更添加示例或 fixture。
- [ ] 规划或添加验证/测试。
- [ ] 跨模块数据有类型或有文档。
- [ ] 核心边界不使用未文档化 JSON、无类型 map 或静默 schema drift。

## Storage separation

- [ ] Save DB 和 cache DB 分离。
- [ ] Cache clean/rebuild 不能删除 save data。
- [ ] Save migration 显式且考虑 backup。
- [ ] Runtime 尽可能只读打开 cache。
