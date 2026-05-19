# CLI Contract

The initial CLI is a single `phanes` binary with subcommands for runtime, builder, cache, save, and launcher development.

## Runtime commands

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

Contract:

- `--bind` defaults to `127.0.0.1`.
- `--bind 0.0.0.0` fails.
- Non-loopback bind addresses fail.
- Missing or invalid cache fails startup with a diagnostic that tells the user to run builder first.
- Startup never runs builder and never downloads resources.

## Builder commands

```text
phanes builder detect \
  --roots <local-path>...

phanes builder build \
  --input <local-path> \
  --output ./cache \
  --mode local-cache

phanes builder build \
  --output ./cache \
  --mode embedded-minimal

phanes builder verify \
  --cache-dir ./cache \
  --deep
```

Contract:

- Inputs are local paths.
- Remote repository URLs are not accepted by the initial contract.
- `embedded-minimal` is for tests/smoke checks and must not include full copyrighted resources.
- `external-import` is explicit, optional, and not default.
- Output includes `manifest.pb`, `cache.sqlite`, blobs, logs, and `.cache-complete`. `manifest.json` is optional and generated for inspection only.

## Cache commands

```text
phanes cache inspect --cache-dir ./cache --json
phanes cache clean --cache-dir ./cache
```

Contract:

- `inspect` reports manifest, schema, source, artifact, and index metadata.
- `clean` may remove generated cache files only.
- `clean` must never touch save data.

## Save commands

```text
phanes save list-profiles --save-db ./save.sqlite
phanes save export --save-db ./save.sqlite --output ./backup.phanes-save
phanes save import --input ./backup.phanes-save --save-db ./save.sqlite
```

Contract:

- Save import/export is explicit.
- Overwrite requires explicit confirmation or a non-interactive force flag in future design.
- Save commands do not inspect or mutate resource cache unless a future contract explicitly says so.

## Launcher command

```text
phanes launcher dev --config ./phanes.config.json
```

Contract:

- Development launcher workflow uses the same runtime/builder/cache contracts as production launcher.
- Launcher starts runtime as a local child process and polls loopback HTTP.
