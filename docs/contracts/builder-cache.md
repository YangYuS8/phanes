# Builder and Resource Cache Contract

The builder prepares local resource-cache output before runtime startup. Runtime consumes completed cache output only.

## Input contract

Allowed source kinds:

```text
embedded-minimal
local-install
local-archive
external-import
```

Rules:

- Inputs are local paths or embedded test fixtures.
- Remote URLs are not part of the initial contract.
- GC-Resources is never required.
- `external-import` is optional compatibility behavior and must be explicitly selected.
- Builder must not redistribute full copyrighted resource sets.

## Output layout

```text
cache-root/
  manifest.pb
  manifest.json
  cache.sqlite
  blobs/
    sha256/
      ab/
        <digest>
  logs/
    build.log
```

`manifest.pb` is authoritative. `manifest.json` is for human inspection and must be generated consistently from the same manifest data if present.
`manifest.json` is optional and is never an independent source of truth.

## Manifest requirements

The cache manifest records:

- schema version
- cache ID
- source list and source kinds
- artifact list and relative paths
- digests
- indexes and record counts
- build timestamp
- compatibility metadata such as game version/region when known

## Atomicity

Builder must avoid exposing half-built cache output as complete.

Initial required strategy:

- write output with `.build-in-progress` present
- verify required manifest/artifacts/indexes
- remove `.build-in-progress`
- write `.cache-complete`

Builders may also use temporary directories and atomic rename, but the final cache root still needs `.cache-complete` for runtime acceptance.

Runtime must reject output marked in-progress or missing completion evidence.

## Runtime acceptance criteria

Runtime accepts a cache only when:

- manifest exists
- manifest schema version is supported
- required artifacts exist
- `cache.sqlite` can be opened
- selected digest checks pass
- `.cache-complete` exists
- `.build-in-progress` does not exist

Runtime behavior on invalid cache:

- fail startup or enter a documented failed state
- emit a diagnostic telling the user to run `phanes builder build` or `phanes builder verify`
- do not download resources
- do not invoke builder automatically

## Cache lifecycle

- Cache is generated and rebuildable.
- Cache clean/rebuild must not touch save data.
- Cache compatibility should be gated by manifest schema version.
- Cache migrations may prefer rebuild over in-place mutation.
