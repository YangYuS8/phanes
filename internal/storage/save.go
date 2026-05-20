package storage

import (
	"context"
	"database/sql"
)

const SaveSchemaVersion = 1

func MigrateSave(ctx context.Context, db *sql.DB) error {
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
		`CREATE TABLE IF NOT EXISTS profiles (
  profile_id TEXT PRIMARY KEY,
  display_name TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
)`,
		`CREATE TABLE IF NOT EXISTS save_metadata (
  save_id TEXT PRIMARY KEY,
  profile_id TEXT NOT NULL,
  schema_version INTEGER NOT NULL,
  last_played_at TEXT,
  FOREIGN KEY (profile_id) REFERENCES profiles(profile_id)
)`,
		`INSERT INTO schema_version(component, version, updated_at)
VALUES ('save', 1, '1970-01-01T00:00:00Z')
ON CONFLICT(component) DO UPDATE SET version = excluded.version, updated_at = excluded.updated_at`,
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func VerifySaveSchema(ctx context.Context, db *sql.DB) error {
	for _, table := range []string{"schema_version", "profiles", "save_metadata"} {
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
