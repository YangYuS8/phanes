# Module Boundaries

This document defines the first Phanes module boundaries. It is intentionally contract-level: implementation packages may change, but cross-module responsibilities should not drift silently.

## Runtime

Responsibilities:

- Bind to loopback only, defaulting to `127.0.0.1`.
- Expose local HTTP APIs for health, status, diagnostics, session lifecycle, and shutdown.
- Read verified resource-cache output.
- Read/write local save data.
- Emit structured diagnostics.

Non-responsibilities:

- No resource downloads.
- No cache building or source detection.
- No public server hosting.
- No client patching or bypass behavior.

Inputs:

- `RuntimeConfig`
- cache directory containing a completed manifest and cache indexes
- save SQLite path

Outputs:

- loopback HTTP API
- diagnostics/logs
- session state updates in save storage

## Builder

Responsibilities:

- Detect local input sources.
- Normalize user-provided local inputs into a versioned resource cache.
- Verify cache completeness and integrity.
- Produce manifest, indexes, blobs, and build logs.

Non-responsibilities:

- No runtime lifecycle management.
- No required GC-Resources dependency.
- No remote repository downloads in the initial contract.
- No bundled copyrighted full resource sets.

Inputs:

- local paths
- explicit build mode
- output directory

Outputs:

- completed cache directory
- `CacheManifest`
- `BuildResult` / `VerifyCacheResult`

## Launcher

Responsibilities:

- Provide user-facing workflow.
- Validate local config and cache state.
- Start/stop runtime as a child process.
- Run builder explicitly when the user chooses to prepare/rebuild cache.
- Show logs and diagnostics.

Non-responsibilities:

- No silent system proxy mutation.
- No client patching.
- No public server controls.
- No hidden network-first daemon behavior.

Communication:

- child process control
- runtime ready/status file
- loopback HTTP polling
- stdout/stderr log capture

## Protocol

Responsibilities:

- Hold versioned protobuf definitions.
- Provide generated Go types for contract messages.
- Define compatibility and schema versioning rules.

Non-responsibilities:

- No monolithic global schema for all future data.
- No undocumented maps for core boundaries.

## Storage

Responsibilities:

- Keep save data and resource cache separate.
- Provide explicit migrations.
- Protect save data from cache cleanup/rebuild flows.
- Support backup/export/import boundaries for save data.

Recommended split:

```text
save.sqlite   user-owned, persistent, careful migrations
cache.sqlite  generated, rebuildable, read-mostly
```
