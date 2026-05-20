# Phanes

中文说明见：[README.zh-CN.md](README.zh-CN.md).

Phanes is a local-first offline runtime, resource-cache builder, and launcher framework for sandbox research.

The intended flow is:

```text
Launcher checks local config and resource cache
  -> starts a localhost-only runtime on 127.0.0.1
  -> runtime reads local save data and local resource cache
  -> launcher shuts runtime down when the game exits
```

Phanes is **not** a Grasscutter fork. Grasscutter, Cultivation, and GC-Resources are reference material only; Phanes must keep its own architecture and contracts.

## Non-negotiable boundaries

- Runtime defaults to `127.0.0.1` and must not bind to `0.0.0.0`.
- Runtime startup uses local files only and must not download resources.
- GC-Resources must not be required.
- Full copyrighted game resources must not be bundled or redistributed.
- Do not implement or document client patching, anti-cheat bypass, protection bypass, official service authentication bypass, unauthorized online interaction, or commercial service replacement.
- Save data and resource cache are separate: save data is user-owned and persistent; cache data is generated and rebuildable.

## Contract-first design

Before implementation, cross-module behavior is defined through versioned contracts:

- ADRs in [`docs/adr/`](docs/adr/)
- module boundaries in [`docs/contracts/module-boundaries.md`](docs/contracts/module-boundaries.md)
- protobuf package plan in [`docs/contracts/protobuf.md`](docs/contracts/protobuf.md)
- Go interface contracts in [`docs/contracts/go-interfaces.md`](docs/contracts/go-interfaces.md)
- SQLite schema contracts in [`docs/contracts/sqlite-schema.md`](docs/contracts/sqlite-schema.md)
- CLI contract in [`docs/contracts/cli.md`](docs/contracts/cli.md)
- runtime HTTP/status API in [`docs/contracts/runtime-http-api.md`](docs/contracts/runtime-http-api.md)
- builder/cache contract in [`docs/contracts/builder-cache.md`](docs/contracts/builder-cache.md)
- compliance checklist in [`docs/contracts/compliance-checklist.md`](docs/contracts/compliance-checklist.md)

Preferred stack:

- Go for runtime, builder, and CLI
- Protocol Buffers for typed cross-module messages
- SQLite for save data and resource cache indexes
- Tauri 2 for the launcher

## Current status

The project is in contract-design phase. Implementation should begin only after the relevant ADRs, examples/fixtures, and validation checks are in place.

## Documentation site

Project documentation is intended to be managed with Docusaurus and published to GitHub Pages through GitHub Actions. The checked-in contract docs under `docs/` are the source for that site.

Repository setup requirement: enable **Settings → Pages → Source → GitHub Actions** before the `pages` workflow can deploy.

Local docs commands:

```bash
npm install
npm run start
npm run build
npm run build:en
npm run build:zh
```

The Docusaurus site provides English and Simplified Chinese documentation. English source docs live under `docs/`; Chinese localized docs live under `i18n/zh-Hans/docusaurus-plugin-content-docs/current/`.

## Releases

Tagged GitHub Actions release workflows build Windows-focused release artifacts for `windows/amd64` and `windows/arm64`, package them as `.zip` files, generate `SHA256SUMS.txt`, and create a GitHub Release.

## License

Phanes is licensed under the Apache License, Version 2.0. See [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE).
