# SQLite Schema 契约

Phanes 分离存档数据和资源缓存：

```text
save.sqlite   用户所有、持久化、migration 需谨慎
cache.sqlite  Builder 生成、可重建、以读为主
```

缓存清理或重建命令不得删除或修改 `save.sqlite`。

## Save database

初始 schema 区域：

```sql
CREATE TABLE profiles (
  profile_id TEXT PRIMARY KEY,
  display_name TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE save_metadata (
  save_id TEXT PRIMARY KEY,
  profile_id TEXT NOT NULL,
  schema_version INTEGER NOT NULL,
  last_played_at TEXT,
  FOREIGN KEY (profile_id) REFERENCES profiles(profile_id)
);
```

暂不添加通用 KV 存档表。未来如需 `kv_state`，必须先为每个 namespace 定义 key 格式、value 编码、owner、migration 和兼容规则。

## Cache database

初始 schema 区域：

```sql
CREATE TABLE cache_manifest (
  cache_id TEXT PRIMARY KEY,
  schema_version TEXT NOT NULL,
  game_version TEXT,
  created_at TEXT NOT NULL,
  manifest_blob BLOB NOT NULL
);

CREATE TABLE sources (...);
CREATE TABLE artifacts (...);
CREATE TABLE resource_index (...);
```

Cache DB 要求：

- 仅由 Builder 生成。
- Runtime 尽可能只读打开。
- manifest schema version 必须受支持。
- 可删除、可重建，不影响 save DB。
- 不要求远程 source。

## Migration 规则

- Save 和 cache migration 分离。
- Save migration 保守，破坏性操作需要用户可见。
- Cache migration 可优先选择 rebuild。
- 在生产数据出现前先用 fixture 覆盖 migration。
