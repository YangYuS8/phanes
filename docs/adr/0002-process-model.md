# ADR 0002: Process Model

## Status

Proposed

## Context

Phanes needs a launcher, runtime, and builder. Combining them too early would hide lifecycle responsibilities and make it easier for runtime startup to accidentally prepare resources or perform launcher-specific work.

## Decision

Use separate logical processes for the first implementation phase:

```text
Launcher process
  owns user-facing workflow
  validates local config/cache state
  starts/stops Runtime as a child process
  starts Builder explicitly when the user requests cache preparation
  displays logs and diagnostics

Runtime process
  binds only to loopback
  exposes local HTTP/status/session APIs
  reads verified resource cache
  reads/writes save data
  never downloads or builds resources

Builder process
  reads user-provided local inputs
  writes normalized local cache output
  verifies cache completeness/integrity
  has no runtime dependency
```

The initial launcher/runtime communication model is child-process supervision plus loopback HTTP polling. Avoid custom IPC, gRPC, plugin systems, and background daemon managers until a real need appears.

## Consequences

- CLI validation can exercise runtime and builder independently.
- Crashes and logs are easier to isolate.
- Runtime has a narrow contract: consume cache + save data, expose local APIs.
- Launcher remains a supervisor and UX surface rather than a hidden system mutator.

## Open questions

- Whether builder should be invoked as a child process by the launcher or as a Go library behind a CLI wrapper in early development.
- Whether the launcher frontend will be Svelte, React, or another Tauri-compatible stack.

## Related

- [Runtime HTTP API](../contracts/runtime-http-api.md)
- [CLI contract](../contracts/cli.md)
