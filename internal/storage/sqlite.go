package storage

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const driverName = "sqlite"

func OpenReadWrite(path string) (*sql.DB, error) {
	dsn, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return sql.Open(driverName, dsn)
}

func OpenReadOnly(path string) (*sql.DB, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return sql.Open(driverName, fmt.Sprintf("file:%s?mode=ro", filepath.ToSlash(abs)))
}

func tableExists(ctx context.Context, db *sql.DB, table string) (bool, error) {
	var name string
	err := db.QueryRowContext(ctx, `SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return name == table, nil
}
