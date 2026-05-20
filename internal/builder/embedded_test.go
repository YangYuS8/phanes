package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/YangYuS8/phanes/internal/contracts"
)

func TestBuildEmbeddedMinimalIsVerifiable(t *testing.T) {
	output := filepath.Join(t.TempDir(), "cache")
	manifest, err := BuildEmbeddedMinimal(output)
	if err != nil {
		t.Fatalf("BuildEmbeddedMinimal: %v", err)
	}
	if manifest.GetCacheId() != EmbeddedMinimalCacheID {
		t.Fatalf("cache_id = %q", manifest.GetCacheId())
	}
	if _, err := os.Stat(filepath.Join(output, contracts.CacheInProgressMarker)); !os.IsNotExist(err) {
		t.Fatalf("in-progress marker should be absent after build, stat err=%v", err)
	}
	if _, err := contracts.VerifyCacheRoot(output, contracts.CacheVerifyOptions{DeepHashCheck: true}); err != nil {
		t.Fatalf("VerifyCacheRoot generated cache: %v", err)
	}
}

func TestBuildEmbeddedMinimalIsDeterministic(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first")
	second := filepath.Join(dir, "second")
	if _, err := BuildEmbeddedMinimal(first); err != nil {
		t.Fatal(err)
	}
	if _, err := BuildEmbeddedMinimal(second); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"manifest.pb", "manifest.json", "cache.sqlite"} {
		firstBytes, err := os.ReadFile(filepath.Join(first, name))
		if err != nil {
			t.Fatal(err)
		}
		secondBytes, err := os.ReadFile(filepath.Join(second, name))
		if err != nil {
			t.Fatal(err)
		}
		if string(firstBytes) != string(secondBytes) {
			t.Fatalf("%s is not deterministic", name)
		}
	}
}
