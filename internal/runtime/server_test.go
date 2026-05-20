package runtime

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/YangYuS8/phanes/internal/builder"
	"github.com/YangYuS8/phanes/internal/config"
	"github.com/YangYuS8/phanes/internal/contracts"
)

func TestRuntimeHTTPStatusLifecycle(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, "cache")
	if _, err := builder.BuildEmbeddedMinimal(cacheDir); err != nil {
		t.Fatal(err)
	}

	cfg := config.Config{
		Runtime: config.RuntimeConfig{BindHost: contracts.DefaultBindHost, HTTPPort: 0},
		Cache:   config.CacheConfig{Dir: cacheDir},
		Save:    config.SaveConfig{DBPath: filepath.Join(dir, "save.sqlite")},
	}

	ready := make(chan *Server, 1)
	server, err := NewServer(cfg, Options{OnReady: func(server *Server) { ready <- server }})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errc := make(chan error, 1)
	go func() { errc <- server.Run(ctx) }()

	var running *Server
	select {
	case running = <-ready:
	case <-time.After(5 * time.Second):
		t.Fatal("runtime did not become ready")
	}

	var health map[string]any
	getJSON(t, running.URL()+"/v1/healthz", &health)
	if health["ok"] != true {
		t.Fatalf("health ok = %v", health["ok"])
	}

	var status map[string]any
	getJSON(t, running.URL()+"/v1/status", &status)
	if status["bind_host"] != contracts.DefaultBindHost {
		t.Fatalf("bind_host = %v", status["bind_host"])
	}
	if status["cache_id"] != builder.EmbeddedMinimalCacheID {
		t.Fatalf("cache_id = %v", status["cache_id"])
	}

	resp, err := http.Post(running.URL()+"/v1/runtime/shutdown", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	select {
	case err := <-errc:
		if err != nil {
			t.Fatalf("runtime returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runtime did not shutdown")
	}
}

func TestNewServerRejectsInvalidBind(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, "cache")
	if _, err := builder.BuildEmbeddedMinimal(cacheDir); err != nil {
		t.Fatal(err)
	}
	_, err := NewServer(config.Config{
		Runtime: config.RuntimeConfig{BindHost: "0.0.0.0"},
		Cache:   config.CacheConfig{Dir: cacheDir},
		Save:    config.SaveConfig{DBPath: filepath.Join(dir, "save.sqlite")},
	}, Options{})
	if err == nil {
		t.Fatal("expected invalid bind error")
	}
}

func getJSON(t *testing.T, url string, target any) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status %d", url, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatal(err)
	}
}
