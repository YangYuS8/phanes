# ADR 0003: Contract Packages

## Status

Proposed

## Context

Phanes should avoid undocumented JSON shapes, untyped maps, large global structs, and accidental schema drift across module boundaries. Go and protobuf should be used where typed contracts improve stability.

## Decision

Define small versioned protobuf packages for cross-module messages:

```text
phanes.common.v1       shared primitives, diagnostics, digests, versions
phanes.cache.v1        generated resource-cache manifest and artifacts
phanes.builder.v1      build/verify requests, results, and progress
phanes.runtime.v1      runtime config, status, lifecycle, sessions
phanes.launcher.v1     launcher config, process status, log events
phanes.save.v1         profile and save metadata identity
```

Use generated Go types at module boundaries. JSON is allowed for human-editable config, UI-facing data, logs, import/export envelopes, and process-ready files, but those JSON shapes must still be documented.

## Consequences

- Contract changes are explicit and reviewable.
- Runtime, builder, launcher, and storage can evolve with version gates.
- The project can add generated code later without changing the design vocabulary.
- Early contracts stay small instead of trying to model complete game state.

## Open questions

- Exact protobuf generation tooling and Go module path.
- Whether generated types should live only under `internal/` initially or be re-exported from `pkg/` once external API needs exist.

## Related

- [Protobuf contracts](../contracts/protobuf.md)
- [Go interface contracts](../contracts/go-interfaces.md)
