package contracts

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
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

	mustWrite(t, filepath.Join(dir, "manifest.pb"), "contract-fixture")
	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheSQLiteMissing) {
		t.Fatalf("missing sqlite error = %v, want ErrCacheSQLiteMissing", err)
	}

	mustWrite(t, filepath.Join(dir, "cache.sqlite"), "")
	if err := ValidateCacheRoot(dir); err != nil {
		t.Fatalf("complete fixture cache returned error: %v", err)
	}
}

func TestValidateCacheRootRejectsInProgressCache(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, CacheCompleteMarker), "contract-fixture")
	mustWrite(t, filepath.Join(dir, CacheInProgressMarker), "building")
	mustWrite(t, filepath.Join(dir, "manifest.pb"), "contract-fixture")
	mustWrite(t, filepath.Join(dir, "cache.sqlite"), "")

	if err := ValidateCacheRoot(dir); !errors.Is(err, ErrCacheInProgress) {
		t.Fatalf("in-progress cache error = %v, want ErrCacheInProgress", err)
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
