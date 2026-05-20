package builder

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	cachev1 "github.com/YangYuS8/phanes/gen/go/phanes/cache/v1"
	commonv1 "github.com/YangYuS8/phanes/gen/go/phanes/common/v1"
	"github.com/YangYuS8/phanes/internal/contracts"
	"github.com/YangYuS8/phanes/internal/storage"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const EmbeddedMinimalCacheID = "embedded-minimal-v1"

func BuildEmbeddedMinimal(outputDir string) (*cachev1.CacheManifest, error) {
	if outputDir == "" {
		return nil, fmt.Errorf("output directory is required")
	}
	if err := os.RemoveAll(outputDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "logs"), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, contracts.CacheInProgressMarker), []byte("building\n"), 0o644); err != nil {
		return nil, err
	}

	cacheSQLite := filepath.Join(outputDir, "cache.sqlite")
	if err := createMinimalCacheSQLite(cacheSQLite); err != nil {
		return nil, err
	}
	sqliteData, err := os.ReadFile(cacheSQLite)
	if err != nil {
		return nil, err
	}

	manifest := embeddedMinimalManifest(uint64(len(sqliteData)), sha256Hex(sqliteData))
	manifestBytes, err := proto.Marshal(manifest)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "manifest.pb"), manifestBytes, 0o644); err != nil {
		return nil, err
	}

	jsonBytes, err := protojson.MarshalOptions{Multiline: true, Indent: "  ", UseEnumNumbers: false}.Marshal(manifest)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "manifest.json"), append(jsonBytes, '\n'), 0o644); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "logs", "build.log"), []byte("embedded-minimal cache fixture generated\n"), 0o644); err != nil {
		return nil, err
	}
	if err := os.Remove(filepath.Join(outputDir, contracts.CacheInProgressMarker)); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, contracts.CacheCompleteMarker), []byte("contract-fixture\n"), 0o644); err != nil {
		return nil, err
	}
	return manifest, nil
}

func embeddedMinimalManifest(cacheSQLiteSize uint64, cacheSQLiteSHA256 string) *cachev1.CacheManifest {
	return &cachev1.CacheManifest{
		CacheId: EmbeddedMinimalCacheID,
		SchemaVersion: &commonv1.Version{
			Major:      0,
			Minor:      1,
			Patch:      0,
			Prerelease: "contract-fixture",
		},
		GameRegion:  "fixture",
		GameVersion: "fixture",
		Sources: []*cachev1.CacheSource{
			{
				SourceId:    "embedded-minimal",
				Kind:        cachev1.SourceKind_SOURCE_KIND_EMBEDDED_MINIMAL,
				DisplayName: "Embedded minimal contract fixture",
				Digest:      &commonv1.FileDigest{Algorithm: "sha256", Hex: sha256Hex([]byte("embedded-minimal"))},
			},
		},
		Artifacts: []*cachev1.CacheArtifact{
			{
				ArtifactId:   "cache-sqlite",
				Kind:         cachev1.ArtifactKind_ARTIFACT_KIND_SQLITE_DB,
				RelativePath: "cache.sqlite",
				SizeBytes:    cacheSQLiteSize,
				Digest:       &commonv1.FileDigest{Algorithm: "sha256", Hex: cacheSQLiteSHA256},
			},
		},
		Indexes: []*cachev1.CacheIndex{
			{Name: "resource_index", Version: "0.1.0", RecordCount: 0},
		},
		CreatedAt: timestamppb.New(time.Unix(0, 0).UTC()),
	}
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func createMinimalCacheSQLite(path string) error {
	db, err := storage.OpenReadWrite(path)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := storage.MigrateCache(context.Background(), db); err != nil {
		return err
	}
	return vacuumSQLite(db)
}

func vacuumSQLite(db *sql.DB) error {
	_, err := db.Exec("VACUUM")
	return err
}
