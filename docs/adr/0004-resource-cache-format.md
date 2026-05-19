# ADR 0004: Resource Cache Format

## Status

Proposed

## Context

Phanes must not require GC-Resources or any remote resource repository. The runtime must not download resources. Therefore the builder must produce a normalized local cache that the runtime can consume without network access.

## Decision

The resource cache is a generated, rebuildable, read-mostly directory:

```text
cache-root/
  manifest.pb
  manifest.json      # optional inspection view generated from manifest.pb
  cache.sqlite
  blobs/
    sha256/
      ab/
        <digest>
  logs/
    build.log
```

Required build modes:

- `embedded-minimal`: tiny test dataset for contract and smoke tests.
- `local-cache`: primary path generated from user-provided local inputs.

Allowed source kinds:

- `local-install`: user-selected local installation directory.
- `local-archive`: user-selected local archive.
- `embedded-minimal`: bundled safe test fixture.
- `external-import`: user-selected local directory or archive handled by a compatibility adapter; never remote, never required, and never an implicit GC-Resources dependency.

The cache manifest records schema version, cache ID, source list, artifact list, indexes, digests, creation time, and compatibility metadata. `manifest.pb` is authoritative. `manifest.json`, when present, is a generated human-readable mirror and is not an independent source of truth.

Builder writes cache output atomically. The initial contract requires a `.cache-complete` marker in the final cache root and requires runtime to reject caches containing `.build-in-progress` or missing `.cache-complete`.

Runtime accepts only completed and verified cache output. Missing or invalid cache is a startup failure with a clear diagnostic.

## Consequences

- Runtime remains local/offline and deterministic at startup.
- Cache can be deleted and rebuilt without touching save data.
- Cache compatibility can be gated by manifest schema version.
- Detailed resource indexing can be added incrementally behind the manifest contract.

## Open questions

- Exact minimal fixture content for `embedded-minimal` without bundling copyrighted resources.
- Exact canonical contents of the `.cache-complete` marker.

## Related

- [Builder/cache contract](../contracts/builder-cache.md)
- [SQLite schema contract](../contracts/sqlite-schema.md)
