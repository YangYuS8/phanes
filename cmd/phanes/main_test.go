package main

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/YangYuS8/phanes/internal/builder"
	"github.com/YangYuS8/phanes/internal/config"
	"github.com/YangYuS8/phanes/internal/contracts"
	runtimehttp "github.com/YangYuS8/phanes/internal/runtime"
)

func TestRuntimeStartRejectsUnsafeBind(t *testing.T) {
	err := run([]string{"runtime", "start", "--bind", "0.0.0.0", "--cache-dir", "cache", "--save-db", "save/save.sqlite"})
	if err == nil || !strings.Contains(err.Error(), "127.0.0.1") {
		t.Fatalf("expected unsafe bind error, got %v", err)
	}
}

func TestRuntimeStartValidatesExampleConfig(t *testing.T) {
	err := run([]string{"runtime", "start", "--config", "../../examples/config/phanes.config.json"})
	if err == nil || !strings.Contains(err.Error(), "verify cache before runtime start") {
		t.Fatalf("expected cache verification error, got %v", err)
	}
}

func TestBuilderVerifyFixture(t *testing.T) {
	err := run([]string{"builder", "verify", "--cache-dir", "../../examples/cache/embedded-minimal", "--deep"})
	if err != nil {
		t.Fatalf("builder verify fixture: %v", err)
	}
}

func TestBuilderBuildEmbeddedMinimal(t *testing.T) {
	output := filepath.Join(t.TempDir(), "cache")
	if err := run([]string{"builder", "build", "--mode", "embedded-minimal", "--output", output}); err != nil {
		t.Fatalf("builder build embedded-minimal: %v", err)
	}
	if err := run([]string{"builder", "verify", "--cache-dir", output, "--deep"}); err != nil {
		t.Fatalf("builder verify generated embedded-minimal: %v", err)
	}
}

func TestCacheCleanRejectsSaveInsideCache(t *testing.T) {
	err := run([]string{"cache", "clean", "--cache-dir", "cache", "--save-db", "cache/save.sqlite"})
	if err == nil {
		t.Fatal("expected cache/save separation error")
	}
}

func TestVersionCommand(t *testing.T) {
	if err := run([]string{"version"}); err != nil {
		t.Fatalf("version command: %v", err)
	}
}

func TestRuntimeStatusAndStopCommands(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, "cache")
	if _, err := builder.BuildEmbeddedMinimal(cacheDir); err != nil {
		t.Fatal(err)
	}

	ready := make(chan *runtimehttp.Server, 1)
	server, err := runtimehttp.NewServer(config.Config{
		Runtime: config.RuntimeConfig{BindHost: contracts.DefaultBindHost, HTTPPort: 0},
		Cache:   config.CacheConfig{Dir: cacheDir},
		Save:    config.SaveConfig{DBPath: filepath.Join(dir, "save.sqlite")},
	}, runtimehttp.Options{OnReady: func(server *runtimehttp.Server) { ready <- server }})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errc := make(chan error, 1)
	go func() { errc <- server.Run(ctx) }()

	var running *runtimehttp.Server
	select {
	case running = <-ready:
	case <-time.After(5 * time.Second):
		t.Fatal("runtime did not become ready")
	}

	if err := run([]string{"runtime", "status", "--url", running.URL()}); err != nil {
		t.Fatalf("runtime status: %v", err)
	}
	if err := run([]string{"runtime", "stop", "--url", running.URL()}); err != nil {
		t.Fatalf("runtime stop: %v", err)
	}

	select {
	case err := <-errc:
		if err != nil {
			t.Fatalf("runtime returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runtime did not stop")
	}
}
