# Runtime HTTP API Contract

The runtime exposes a local HTTP API for launcher/CLI diagnostics and lifecycle control. The listener must bind only to loopback, defaulting to `127.0.0.1`.

## Endpoints

```text
GET  /v1/healthz
GET  /v1/status
GET  /v1/diagnostics
POST /v1/session/start
POST /v1/session/stop
POST /v1/runtime/shutdown
```

Responses are JSON encodings of protobuf messages.

## `GET /v1/healthz`

Purpose: cheap process health check.

Returns `phanes.runtime.v1.HealthzResponse`.

## `GET /v1/status`

Returns `phanes.runtime.v1.RuntimeStatus`.

Required fields include:

- state
- actual bind host
- actual HTTP port
- cache ID when cache is accepted
- active profile ID when a session is active
- problems
- started timestamp

## `GET /v1/diagnostics`

Returns `phanes.runtime.v1.DiagnosticsResponse`.

## `POST /v1/session/start`

Request: `phanes.runtime.v1.StartSessionRequest`.

Response: `phanes.runtime.v1.StartSessionResult`.

This starts a local/offline session for a local profile. It is not an official authentication flow and must not contact official services.

## `POST /v1/session/stop`

Request: `phanes.runtime.v1.StopSessionRequest`.

Response: `phanes.runtime.v1.StopSessionResult`.

## `POST /v1/runtime/shutdown`

Request: `phanes.runtime.v1.StopRuntimeRequest`.

Response: `phanes.runtime.v1.StopRuntimeResult`.

Launcher should use this endpoint for graceful shutdown before killing the child process.

## Binding and validation rules

- Missing bind config resolves to literal `127.0.0.1`.
- Explicit literal `127.0.0.1` is accepted.
- `localhost`, `::1`, and other loopback aliases require a documented validation update before use.
- `0.0.0.0` is rejected.
- Public/non-loopback addresses are rejected.
- `/v1/status` must expose the actual bind host and port.

## Out of scope

- public admin API
- reverse proxy support
- public TLS/auth configuration
- remote service discovery
- multi-user server controls
