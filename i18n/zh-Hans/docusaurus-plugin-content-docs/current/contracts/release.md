# 发布契约

Phanes release 由 GitHub Actions 从 Git tag 生成。

## 主要平台

初始发布目标以 Windows 为主：

- `windows/amd64`
- `windows/arm64`

未来可以增加其他平台，但 Windows artifact 是 launcher/runtime 工作流的主要发布输出。

## Tag pipeline

Release job 只在 tag push 运行：

```text
on:
  push:
    tags:
      - 'v*'
```

## Artifacts

发布打包输出：

```text
release/phanes-windows-amd64.zip
release/phanes-windows-arm64.zip
release/SHA256SUMS.txt
```

压缩包包含：

- `phanes-<goos>-<goarch>.exe`
- `LICENSE`
- `NOTICE`
- `README.md`
- `README.zh-CN.md`

## GitHub Release

`create-release` job 使用 `gh release create` 和 workflow `GITHUB_TOKEN`。Release assets 包含 Windows archives 和 `SHA256SUMS.txt`。

## 验证期望

- `go test ./...`
- `buf lint`
- Docusaurus 双语构建
- `./cmd/phanes` Windows cross-compile
- 为所有发布 archive 生成 SHA256 checksum

## 边界

发布产物不得捆绑版权游戏资源、GC-Resources、client patching 工具、anti-cheat bypass 工具或远程服务替代资产。
