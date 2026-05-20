# Release Contract

Phanes releases are produced by GitHub Actions from Git tags.

## Primary platform

The initial release target is Windows:

- `windows/amd64`
- `windows/arm64`

Other platforms may be added later, but Windows artifacts are the primary release output for the launcher/runtime workflow.

## Tag pipeline

Release jobs run only for tag pushes:

```text
on:
  push:
    tags:
      - 'v*'
```

## Artifacts

Release packaging produces:

```text
release/phanes-windows-amd64.zip
release/phanes-windows-arm64.zip
release/SHA256SUMS.txt
```

The archives include:

- `phanes-<goos>-<goarch>.exe`
- `LICENSE`
- `NOTICE`
- `README.md`
- `README.zh-CN.md`

## GitHub Release

The `create-release` job uses `gh release create` with the workflow `GITHUB_TOKEN`. Release assets include the Windows archives and `SHA256SUMS.txt`.

## Validation expectations

- `go test ./...`
- `buf lint`
- Docusaurus multilingual build
- Windows cross-compile of `./cmd/phanes`
- SHA256 checksums generated for every published archive

## Boundaries

Release artifacts must not bundle copyrighted game resources, GC-Resources, client patching tools, anti-cheat bypass tools, or remote service replacement assets.
