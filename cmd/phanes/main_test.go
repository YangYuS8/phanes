package main

import (
	"strings"
	"testing"
)

func TestRuntimeStartRejectsUnsafeBind(t *testing.T) {
	err := run([]string{"runtime", "start", "--bind", "0.0.0.0", "--cache-dir", "cache", "--save-db", "save/save.sqlite"})
	if err == nil || !strings.Contains(err.Error(), "127.0.0.1") {
		t.Fatalf("expected unsafe bind error, got %v", err)
	}
}

func TestRuntimeStartValidatesExampleConfig(t *testing.T) {
	err := run([]string{"runtime", "start", "--config", "../../examples/config/phanes.config.json"})
	if err == nil || !strings.Contains(err.Error(), "implementation pending") {
		t.Fatalf("expected validated not implemented error, got %v", err)
	}
}

func TestBuilderVerifyFixture(t *testing.T) {
	err := run([]string{"builder", "verify", "--cache-dir", "../../examples/cache/embedded-minimal"})
	if err != nil {
		t.Fatalf("builder verify fixture: %v", err)
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
