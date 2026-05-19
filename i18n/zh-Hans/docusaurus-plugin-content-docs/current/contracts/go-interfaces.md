# Go Interface 契约

Interface 应保持小，并放在靠近消费者的位置。这里是初始边界契约，不是最终包名要求。

## Runtime 访问 cache

```go
type CacheReader interface {
    Manifest(ctx context.Context) (*cachev1.CacheManifest, error)
    OpenResource(ctx context.Context, id string) (io.ReadCloser, ResourceInfo, error)
    Verify(ctx context.Context, opts VerifyOptions) (*builderv1.VerifyCacheResult, error)
}
```

Runtime 应尽可能以只读方式打开 cache。缺失资源应产生诊断，而不是触发下载或构建。

## Save storage

```go
type SaveStore interface {
    GetProfile(ctx context.Context, id string) (*savev1.Profile, error)
    ListProfiles(ctx context.Context) ([]*savev1.Profile, error)
    BeginSession(ctx context.Context, profileID string) (*runtimev1.StartSessionResult, error)
    EndSession(ctx context.Context, sessionID string) error
}
```

存档数据属于用户。破坏性 migration 必须显式，并应先支持 backup/export。

## Runtime service

```go
type RuntimeService interface {
    Status(ctx context.Context) (*runtimev1.RuntimeStatus, error)
    StartSession(ctx context.Context, req *runtimev1.StartSessionRequest) (*runtimev1.StartSessionResult, error)
    Shutdown(ctx context.Context, req *runtimev1.StopRuntimeRequest) (*runtimev1.StopRuntimeResult, error)
}
```

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

Builder 输入是本地路径。`external-import` 是显式可选兼容行为，不得成为 GC-Resources 依赖。

## Launcher supervision

```go
type RuntimeSupervisor interface {
    Start(ctx context.Context, spec launcherv1.RuntimeProcessSpec) error
    Stop(ctx context.Context) error
    Status(ctx context.Context) (*launcherv1.RuntimeProcessStatus, error)
    Logs(ctx context.Context) (<-chan *launcherv1.LogEvent, error)
}
```

Launcher 启动本地子进程并轮询 loopback HTTP；不得静默修改系统代理或 patch client。
