package contracts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	cachev1 "github.com/YangYuS8/phanes/gen/go/phanes/cache/v1"
	"github.com/YangYuS8/phanes/internal/storage"
	"google.golang.org/protobuf/proto"
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
	ErrCacheManifestInvalid = errors.New("cache manifest is invalid")
	ErrCacheSQLiteMissing   = errors.New("cache sqlite index is missing")
	ErrCacheSQLiteInvalid   = errors.New("cache sqlite index is invalid")
	ErrCacheArtifactInvalid = errors.New("cache artifact is invalid")
)

type CacheVerifyOptions struct {
	DeepHashCheck bool
}

type CacheVerification struct {
	Manifest *cachev1.CacheManifest
}

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
	_, err := VerifyCacheRoot(cacheDir, CacheVerifyOptions{})
	return err
}

func VerifyCacheRoot(cacheDir string, opts CacheVerifyOptions) (*CacheVerification, error) {
	if exists(filepath.Join(cacheDir, CacheInProgressMarker)) {
		return nil, ErrCacheInProgress
	}
	if !exists(filepath.Join(cacheDir, CacheCompleteMarker)) {
		return nil, ErrCacheIncomplete
	}

	manifest, err := ReadCacheManifest(cacheDir)
	if err != nil {
		return nil, err
	}
	if !exists(filepath.Join(cacheDir, "cache.sqlite")) {
		return nil, ErrCacheSQLiteMissing
	}
	if err := validateSQLiteHeader(filepath.Join(cacheDir, "cache.sqlite")); err != nil {
		return nil, err
	}
	cacheDB, err := storage.OpenReadOnly(filepath.Join(cacheDir, "cache.sqlite"))
	if err != nil {
		return nil, err
	}
	defer cacheDB.Close()
	if err := cacheDB.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCacheSQLiteInvalid, err)
	}
	if err := storage.VerifyCacheSchema(context.Background(), cacheDB); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCacheSQLiteInvalid, err)
	}
	if err := validateCacheManifest(cacheDir, manifest, opts); err != nil {
		return nil, err
	}
	return &CacheVerification{Manifest: manifest}, nil
}

func ReadCacheManifest(cacheDir string) (*cachev1.CacheManifest, error) {
	path := filepath.Join(cacheDir, "manifest.pb")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrCacheManifestMissing
	}
	if err != nil {
		return nil, err
	}
	manifest := &cachev1.CacheManifest{}
	if err := proto.Unmarshal(data, manifest); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCacheManifestInvalid, err)
	}
	if manifest.GetCacheId() == "" || manifest.GetSchemaVersion() == nil {
		return nil, fmt.Errorf("%w: missing cache_id or schema_version", ErrCacheManifestInvalid)
	}
	return manifest, nil
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

func validateSQLiteHeader(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	header := make([]byte, 16)
	if _, err := io.ReadFull(file, header); err != nil {
		return fmt.Errorf("%w: %v", ErrCacheSQLiteInvalid, err)
	}
	if string(header) != "SQLite format 3\x00" {
		return ErrCacheSQLiteInvalid
	}
	return nil
}

func validateCacheManifest(cacheDir string, manifest *cachev1.CacheManifest, opts CacheVerifyOptions) error {
	version := manifest.GetSchemaVersion()
	if version.GetMajor() != 0 || version.GetMinor() != 1 {
		return fmt.Errorf("%w: unsupported schema version %d.%d.%d", ErrCacheManifestInvalid, version.GetMajor(), version.GetMinor(), version.GetPatch())
	}
	for _, artifact := range manifest.GetArtifacts() {
		artifactPath, err := safeCachePath(cacheDir, artifact.GetRelativePath())
		if err != nil {
			return err
		}
		info, err := os.Stat(artifactPath)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrCacheArtifactInvalid, artifact.GetRelativePath(), err)
		}
		if artifact.GetSizeBytes() != uint64(info.Size()) {
			return fmt.Errorf("%w: %s size mismatch", ErrCacheArtifactInvalid, artifact.GetRelativePath())
		}
		if opts.DeepHashCheck && artifact.GetDigest() != nil {
			if err := verifySHA256(artifactPath, artifact.GetDigest().GetAlgorithm(), artifact.GetDigest().GetHex()); err != nil {
				return fmt.Errorf("%w: %s: %v", ErrCacheArtifactInvalid, artifact.GetRelativePath(), err)
			}
		}
	}
	return nil
}

func safeCachePath(cacheDir string, relativePath string) (string, error) {
	if relativePath == "" || filepath.IsAbs(relativePath) || startsWithDotDot(filepath.Clean(relativePath)) {
		return "", fmt.Errorf("%w: unsafe relative path %q", ErrCacheArtifactInvalid, relativePath)
	}
	cacheAbs, err := filepath.Abs(cacheDir)
	if err != nil {
		return "", err
	}
	joinedAbs, err := filepath.Abs(filepath.Join(cacheDir, relativePath))
	if err != nil {
		return "", err
	}
	if joinedAbs != cacheAbs && !isWithin(joinedAbs, cacheAbs) {
		return "", fmt.Errorf("%w: path escapes cache root %q", ErrCacheArtifactInvalid, relativePath)
	}
	return joinedAbs, nil
}

func verifySHA256(path string, algorithm string, wantHex string) error {
	if !strings.EqualFold(algorithm, "sha256") {
		return fmt.Errorf("unsupported digest algorithm %q", algorithm)
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	got := hex.EncodeToString(hash.Sum(nil))
	if !strings.EqualFold(got, wantHex) {
		return fmt.Errorf("sha256 mismatch got %s want %s", got, wantHex)
	}
	return nil
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
