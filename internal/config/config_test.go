package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/YangYuS8/phanes/internal/contracts"
)

func TestLoadExampleConfig(t *testing.T) {
	cfg, err := Load(filepath.Join("..", "..", "examples", "config", "phanes.config.json"))
	if err != nil {
		t.Fatalf("Load example config: %v", err)
	}
	if cfg.Runtime.BindHost != contracts.DefaultBindHost {
		t.Fatalf("bind host = %q, want %q", cfg.Runtime.BindHost, contracts.DefaultBindHost)
	}
}

func TestValidateDefaultsBindHost(t *testing.T) {
	cfg := Config{Cache: CacheConfig{Dir: "cache"}, Save: SaveConfig{DBPath: "save/save.sqlite"}}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if cfg.Runtime.BindHost != contracts.DefaultBindHost {
		t.Fatalf("default bind host = %q", cfg.Runtime.BindHost)
	}
}

func TestValidateRejectsUnsafeBindHost(t *testing.T) {
	cfg := Config{Runtime: RuntimeConfig{BindHost: "0.0.0.0"}, Cache: CacheConfig{Dir: "cache"}, Save: SaveConfig{DBPath: "save/save.sqlite"}}
	if err := cfg.Validate(); !errors.Is(err, contracts.ErrInvalidBindHost) {
		t.Fatalf("Validate error = %v, want ErrInvalidBindHost", err)
	}
}

func TestValidateRejectsSaveInsideCache(t *testing.T) {
	cfg := Config{Cache: CacheConfig{Dir: "cache"}, Save: SaveConfig{DBPath: "cache/save.sqlite"}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected save-inside-cache validation error")
	}
}

func TestLoadRejectsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}
