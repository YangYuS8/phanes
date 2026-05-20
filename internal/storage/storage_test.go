package storage

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestMigrateSaveCreatesExpectedTables(t *testing.T) {
	db, err := OpenReadWrite(filepath.Join(t.TempDir(), "save.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := MigrateSave(ctx, db); err != nil {
		t.Fatalf("MigrateSave: %v", err)
	}
	if err := VerifySaveSchema(ctx, db); err != nil {
		t.Fatalf("VerifySaveSchema: %v", err)
	}
}

func TestMigrateCacheCreatesExpectedTables(t *testing.T) {
	db, err := OpenReadWrite(filepath.Join(t.TempDir(), "cache.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := MigrateCache(ctx, db); err != nil {
		t.Fatalf("MigrateCache: %v", err)
	}
	if err := VerifyCacheSchema(ctx, db); err != nil {
		t.Fatalf("VerifyCacheSchema: %v", err)
	}
}

func TestVerifyCacheSchemaReportsMissingTables(t *testing.T) {
	db, err := OpenReadWrite(filepath.Join(t.TempDir(), "cache.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = VerifyCacheSchema(context.Background(), db)
	var missing *MissingTableError
	if !errors.As(err, &missing) {
		t.Fatalf("VerifyCacheSchema error = %v, want MissingTableError", err)
	}
}

func TestOpenReadOnlyDoesNotCreateMissingDB(t *testing.T) {
	db, err := OpenReadOnly(filepath.Join(t.TempDir(), "missing.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err == nil {
		t.Fatal("expected read-only open of missing database to fail")
	}
}
