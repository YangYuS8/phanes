package contracts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultBindHost       = "127.0.0.1"
	CacheCompleteMarker   = ".cache-complete"
	CacheInProgressMarker = ".build-in-progress"
)

var (
	ErrInvalidBindHost      = errors.New("bind host must be literal 127.0.0.1")
	ErrCacheIncomplete      = errors.New("cache is incomplete")
	ErrCacheInProgress      = errors.New("cache build is in progress")
	ErrCacheManifestMissing = errors.New("cache manifest is missing")
	ErrCacheSQLiteMissing   = errors.New("cache sqlite index is missing")
)

func NormalizeBindHost(bindHost string) (string, error) {
	if bindHost == "" {
		return DefaultBindHost, nil
	}
	if bindHost != DefaultBindHost {
		return "", fmt.Errorf("%w: %q", ErrInvalidBindHost, bindHost)
	}
	return bindHost, nil
}

func ValidateCacheRoot(cacheDir string) error {
	if exists(filepath.Join(cacheDir, CacheInProgressMarker)) {
		return ErrCacheInProgress
	}
	if !exists(filepath.Join(cacheDir, CacheCompleteMarker)) {
		return ErrCacheIncomplete
	}
	if !exists(filepath.Join(cacheDir, "manifest.pb")) {
		return ErrCacheManifestMissing
	}
	if !exists(filepath.Join(cacheDir, "cache.sqlite")) {
		return ErrCacheSQLiteMissing
	}
	return nil
}

func CacheCleanAllowed(cacheDir string, saveDBPath string) error {
	cacheAbs, err := filepath.Abs(cacheDir)
	if err != nil {
		return err
	}
	saveAbs, err := filepath.Abs(saveDBPath)
	if err != nil {
		return err
	}
	if cacheAbs == saveAbs {
		return errors.New("cache path must not equal save database path")
	}
	if isWithin(saveAbs, cacheAbs) {
		return errors.New("save database must not be inside cache directory")
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isWithin(path string, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && rel != "" && !startsWithDotDot(rel)
}

func startsWithDotDot(path string) bool {
	return path == ".." || len(path) > 3 && path[:3] == "../"
}
