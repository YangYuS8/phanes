# Go Interface Contracts

Interfaces should remain small and live near consumers. These are initial boundary contracts, not mandatory final package names.

## Runtime-facing cache access

```go
type ResourceInfo struct {
    ID         string
    Type       string
    SizeBytes  uint64
    DigestHex  string
}

type VerifyOptions struct {
    DeepHashCheck bool
}

type CacheReader interface {
    Manifest(ctx context.Context) (*cachev1.CacheManifest, error)
    OpenResource(ctx context.Context, id string) (io.ReadCloser, ResourceInfo, error)
    Verify(ctx context.Context, opts VerifyOptions) (*builderv1.VerifyCacheResult, error)
}
```

Runtime opens cache read-only where possible. Missing resources are runtime diagnostics, not triggers to download or build resources.

## Save storage

```go
type SaveStore interface {
    GetProfile(ctx context.Context, id string) (*savev1.Profile, error)
    ListProfiles(ctx context.Context) ([]*savev1.Profile, error)
    BeginSession(ctx context.Context, profileID string) (*runtimev1.StartSessionResult, error)
    EndSession(ctx context.Context, sessionID string) error
}
```

Save data is user-owned. Destructive migrations must be explicit and should be preceded by backup/export support.

## Runtime service

```go
type RuntimeService interface {
    Status(ctx context.Context) (*runtimev1.RuntimeStatus, error)
    StartSession(ctx context.Context, req *runtimev1.StartSessionRequest) (*runtimev1.StartSessionResult, error)
    Shutdown(ctx context.Context, req *runtimev1.StopRuntimeRequest) (*runtimev1.StopRuntimeResult, error)
}
```

The runtime service validates loopback bind config before listening.

## Builder

```go
type SourceDetector interface {
    Detect(ctx context.Context, roots []string) ([]*builderv1.InputSource, error)
}

type CacheBuilder interface {
    Build(ctx context.Context, req *builderv1.BuildRequest, progress func(*builderv1.BuildProgress)) (*builderv1.BuildResult, error)
    Verify(ctx context.Context, req *builderv1.VerifyCacheRequest) (*builderv1.VerifyCacheResult, error)
}
```

Builder inputs are local paths. `external-import` is optional compatibility behavior and must not become a required GC-Resources path.

## Launcher supervision

```go
type RuntimeSupervisor interface {
    Start(ctx context.Context, spec launcherv1.RuntimeProcessSpec) error
    Stop(ctx context.Context) error
    Status(ctx context.Context) (*launcherv1.RuntimeProcessStatus, error)
    Logs(ctx context.Context) (<-chan *launcherv1.LogEvent, error)
}
```

Launcher starts local child processes and polls loopback HTTP. It must not silently mutate system proxy settings or patch clients.

## Interface rules

- Prefer `internal/` packages until an external API is proven necessary.
- Avoid large `Manager` interfaces that mix runtime, builder, launcher, and storage responsibilities.
- Avoid reflection-heavy registration and global mutable registries.
- Return typed protobuf messages or small explicit structs at boundaries.
