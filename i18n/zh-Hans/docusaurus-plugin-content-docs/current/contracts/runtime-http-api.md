# Runtime HTTP API 契约

Runtime 暴露本地 HTTP API，用于 launcher/CLI 诊断和生命周期控制。监听地址必须只绑定 loopback，初始只接受字面量 `127.0.0.1`。

## Endpoints

```text
GET  /v1/healthz
GET  /v1/status
GET  /v1/diagnostics
POST /v1/session/start
POST /v1/session/stop
POST /v1/runtime/shutdown
```

响应尽可能是 protobuf message 的 JSON 编码。

## `GET /v1/healthz`

返回 `phanes.runtime.v1.HealthzResponse`。

## `GET /v1/status`

返回 `phanes.runtime.v1.RuntimeStatus`，包含状态、实际 bind host、端口、cache ID、活动 profile、问题列表和启动时间。

## `GET /v1/diagnostics`

返回 `phanes.runtime.v1.DiagnosticsResponse`。

## `POST /v1/session/start`

请求：`phanes.runtime.v1.StartSessionRequest`。

响应：`phanes.runtime.v1.StartSessionResult`。

该接口启动本地/offline session，不是官方认证流程，也不得联系官方服务。

## `POST /v1/session/stop`

请求：`phanes.runtime.v1.StopSessionRequest`。

响应：`phanes.runtime.v1.StopSessionResult`。

## `POST /v1/runtime/shutdown`

请求：`phanes.runtime.v1.StopRuntimeRequest`。

响应：`phanes.runtime.v1.StopRuntimeResult`。

## 绑定规则

- 缺省 bind config 解析为 `127.0.0.1`。
- 显式字面量 `127.0.0.1` 被接受。
- `localhost`、`::1` 和其他别名需要先更新验证决策。
- `0.0.0.0` 和公网地址被拒绝。

## 范围外

- 公网 admin API
- reverse proxy 支持
- 公网 TLS/auth 配置
- 远程服务发现
- 多用户服务器控制
