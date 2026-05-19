# Release Contract

Phanes releases are produced by GitLab CI/CD from Git tags.

## Primary platform

The initial release target is Windows:

- `windows/amd64`
- `windows/arm64`

Other platforms may be added later, but Windows artifacts are the primary release output for the launcher/runtime workflow.

## Tag pipeline

Release jobs run only for tag pipelines:

```text
rules:
  - if: $CI_COMMIT_TAG
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

## GitLab Release

The `release:create` job uses the GitLab release keyword with the official `registry.gitlab.com/gitlab-org/cli:latest` image. Release asset links point to packaged job artifacts.

## Validation expectations

- `go test ./...`
- `buf lint`
- Docusaurus multilingual build
- Windows cross-compile of `./cmd/phanes`
- SHA256 checksums generated for every published archive

## Boundaries

Release artifacts must not bundle copyrighted game resources, GC-Resources, client patching tools, anti-cheat bypass tools, or remote service replacement assets.
