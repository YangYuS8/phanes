# CLI 契约

初始 CLI 是单个 `phanes` 二进制，包含 runtime、builder、cache、save 和 launcher 开发命令。

## Runtime 命令

```text
phanes runtime start \
  --config ./phanes.config.json \
  --bind 127.0.0.1 \
  --port 0 \
  --cache-dir ./cache \
  --save-db ./save.sqlite

phanes runtime status \
  --url http://127.0.0.1:<port>

phanes runtime stop \
  --url http://127.0.0.1:<port>
```

契约：

- `--bind` 默认 `127.0.0.1`。
- `--bind 0.0.0.0` 失败。
- 非 loopback 绑定地址失败。
- 缓存缺失或无效时启动失败，并提示先运行 builder。
- 启动不运行 builder，也不下载资源。

## Builder 命令

```text
phanes builder detect --roots <local-path>...
phanes builder build --input <local-path> --output ./cache --mode local-cache
phanes builder build --output ./cache --mode embedded-minimal
phanes builder verify --cache-dir ./cache --deep
```

契约：

- 输入是本地路径。
- 初始契约不接受远程仓库 URL。
- `embedded-minimal` 用于测试/冒烟检查，不包含完整版权资源。
- `external-import` 显式、可选，且不是默认模式。
- 输出包含 `manifest.pb`、`cache.sqlite`、blobs、logs 和 `.cache-complete`；`manifest.json` 是可选检查视图。

## Cache 命令

```text
phanes cache inspect --cache-dir ./cache --json
phanes cache clean --cache-dir ./cache
```

`clean` 只能移除生成的 cache 文件，绝不能触碰 save data。

## Save 命令

```text
phanes save list-profiles --save-db ./save.sqlite
phanes save export --save-db ./save.sqlite --output ./backup.phanes-save
phanes save import --input ./backup.phanes-save --save-db ./save.sqlite
```

Save import/export 必须显式，覆盖应要求确认或未来的非交互 force flag。
