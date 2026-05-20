package storage

import (
	"context"
	"database/sql"
)

const CacheSchemaVersion = 1

func MigrateCache(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS schema_version (
  component TEXT PRIMARY KEY,
  version INTEGER NOT NULL,
  updated_at TEXT NOT NULL DEFAULT '1970-01-01T00:00:00Z'
)`,
		`CREATE TABLE IF NOT EXISTS cache_manifest (
  cache_id TEXT PRIMARY KEY,
  schema_version TEXT NOT NULL,
  game_version TEXT,
  created_at TEXT NOT NULL,
  manifest_blob BLOB NOT NULL
)`,
		`CREATE TABLE IF NOT EXISTS sources (
  source_id TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  local_path TEXT,
  digest_algorithm TEXT,
  digest_hex TEXT
)`,
		`CREATE TABLE IF NOT EXISTS artifacts (
  artifact_id TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  relative_path TEXT NOT NULL,
  size_bytes INTEGER NOT NULL,
  digest_algorithm TEXT NOT NULL,
  digest_hex TEXT NOT NULL
)`,
		`CREATE TABLE IF NOT EXISTS resource_index (
  resource_id TEXT PRIMARY KEY,
  resource_type TEXT NOT NULL,
  logical_path TEXT,
  artifact_id TEXT NOT NULL,
  offset INTEGER,
  length INTEGER,
  digest_hex TEXT,
  FOREIGN KEY (artifact_id) REFERENCES artifacts(artifact_id)
)`,
		`INSERT INTO schema_version(component, version, updated_at)
VALUES ('cache', 1, '1970-01-01T00:00:00Z')
ON CONFLICT(component) DO UPDATE SET version = excluded.version, updated_at = excluded.updated_at`,
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func VerifyCacheSchema(ctx context.Context, db *sql.DB) error {
	for _, table := range []string{"schema_version", "cache_manifest", "sources", "artifacts", "resource_index"} {
		exists, err := tableExists(ctx, db, table)
		if err != nil {
			return err
		}
		if !exists {
			return &MissingTableError{Table: table}
		}
	}
	return nil
}
