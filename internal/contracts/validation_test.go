package contracts

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	cachev1 "github.com/YangYuS8/phanes/gen/go/phanes/cache/v1"
	commonv1 "github.com/YangYuS8/phanes/gen/go/phanes/common/v1"
	"github.com/YangYuS8/phanes/internal/storage"
	"google.golang.org/protobuf/proto"
)

func TestNormalizeBindHost(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "default", input: "", want: DefaultBindHost},
		{name: "literal loopback", input: "127.0.0.1", want: DefaultBindHost},
		{name: "all interfaces rejected", input: "0.0.0.0", wantErr: true},
		{name: "localhost alias rejected", input: "localhost", wantErr: true},
		{name: "ipv6 loopback deferred", input: "::1", wantErr: true},
		{name: "public address rejected", input: "192.0.2.10", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeBindHost(tt.input)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidBindHost) {
					t.Fatalf("expected ErrInvalidBindHost, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("NormalizeBindHost returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("NormalizeBindHost = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateCacheRootRequiresCompleteCache(t *testing.T) {
	dir := t.TempDir()

	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheIncomplete) {
		t.Fatalf("empty cache error = %v, want ErrCacheIncomplete", err)
	}

	mustWrite(t, filepath.Join(dir, CacheCompleteMarker), "contract-fixture")
	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheManifestMissing) {
		t.Fatalf("missing manifest error = %v, want ErrCacheManifestMissing", err)
	}

	mustWriteManifest(t, dir, "cache.sqlite", 0, nil)
	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheSQLiteMissing) {
		t.Fatalf("missing sqlite error = %v, want ErrCacheSQLiteMissing", err)
	}

	mustWriteCacheSQLite(t, filepath.Join(dir, "cache.sqlite"))
	mustWriteManifestForFile(t, dir, "cache.sqlite", nil)
	if err := ValidateCacheRoot(dir); err != nil {
		t.Fatalf("complete fixture cache returned error: %v", err)
	}
}

func TestValidateCacheRootRejectsInProgressCache(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, CacheCompleteMarker), "contract-fixture")
	mustWrite(t, filepath.Join(dir, CacheInProgressMarker), "building")
	mustWriteManifest(t, dir, "cache.sqlite", 0, nil)
	mustWriteCacheSQLite(t, filepath.Join(dir, "cache.sqlite"))

	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheInProgress) {
		t.Fatalf("in-progress cache error = %v, want ErrCacheInProgress", err)
	}
}

func TestValidateCacheRootRejectsInvalidManifest(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, CacheCompleteMarker), "contract-fixture")
	mustWrite(t, filepath.Join(dir, "manifest.pb"), "not protobuf")
	mustWriteCacheSQLite(t, filepath.Join(dir, "cache.sqlite"))

	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheManifestInvalid) {
		t.Fatalf("invalid manifest error = %v, want ErrCacheManifestInvalid", err)
	}
}

func TestValidateCacheRootRejectsInvalidSQLite(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, CacheCompleteMarker), "contract-fixture")
	mustWriteManifest(t, dir, "cache.sqlite", 0, nil)
	mustWrite(t, filepath.Join(dir, "cache.sqlite"), "not sqlite")

	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheSQLiteInvalid) {
		t.Fatalf("invalid sqlite error = %v, want ErrCacheSQLiteInvalid", err)
	}
}

func TestVerifyCacheRootRejectsArtifactPathEscape(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, CacheCompleteMarker), "contract-fixture")
	mustWriteCacheSQLite(t, filepath.Join(dir, "cache.sqlite"))
	mustWriteManifest(t, dir, "../save.sqlite", 0, nil)

	if _, err := VerifyCacheRoot(dir, CacheVerifyOptions{}); !errors.Is(err, ErrCacheArtifactInvalid) {
		t.Fatalf("path escape error = %v, want ErrCacheArtifactInvalid", err)
	}
}

func TestVerifyCacheRootRejectsDigestMismatch(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, CacheCompleteMarker), "contract-fixture")
	mustWriteCacheSQLite(t, filepath.Join(dir, "cache.sqlite"))
	mustWriteManifestForFile(t, dir, "cache.sqlite", &commonv1.FileDigest{Algorithm: "sha256", Hex: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"})

	if _, err := VerifyCacheRoot(dir, CacheVerifyOptions{DeepHashCheck: true}); !errors.Is(err, ErrCacheArtifactInvalid) {
		t.Fatalf("digest mismatch error = %v, want ErrCacheArtifactInvalid", err)
	}
}

func TestCacheCleanAllowedRejectsSaveInsideCache(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, "cache")
	saveDB := filepath.Join(cacheDir, "save.sqlite")

	if err := CacheCleanAllowed(cacheDir, saveDB); err == nil {
		t.Fatal("expected error when save database is inside cache directory")
	}
}

func TestEmbeddedMinimalFixtureSatisfiesCacheRootContract(t *testing.T) {
	fixture := filepath.Join("..", "..", "examples", "cache", "embedded-minimal")

	if err := ValidateCacheRoot(fixture); err != nil {
		t.Fatalf("embedded minimal fixture did not validate: %v", err)
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustWriteManifestForFile(t *testing.T, dir string, artifactPath string, digest *commonv1.FileDigest) {
	t.Helper()
	info, err := os.Stat(filepath.Join(dir, artifactPath))
	if err != nil {
		t.Fatal(err)
	}
	mustWriteManifest(t, dir, artifactPath, uint64(info.Size()), digest)
}

func mustWriteManifest(t *testing.T, dir string, artifactPath string, sizeBytes uint64, digest *commonv1.FileDigest) {
	t.Helper()
	if digest == nil {
		digest = &commonv1.FileDigest{Algorithm: "sha256", Hex: ""}
	}
	manifest := &cachev1.CacheManifest{
		CacheId: "test-cache",
		SchemaVersion: &commonv1.Version{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
		Artifacts: []*cachev1.CacheArtifact{
			{ArtifactId: "artifact", RelativePath: artifactPath, SizeBytes: sizeBytes, Digest: digest},
		},
	}
	data, err := proto.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.pb"), data, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
}

func minimalSQLiteHeader() []byte {
	data := make([]byte, 100)
	copy(data, []byte("SQLite format 3\x00"))
	return data
}

func mustWriteCacheSQLite(t *testing.T, path string) {
	t.Helper()
	db, err := storage.OpenReadWrite(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := storage.MigrateCache(context.Background(), db); err != nil {
		t.Fatal(err)
	}
}
