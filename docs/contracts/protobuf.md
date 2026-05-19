# Protobuf Contracts

Protobuf packages are small, versioned, and focused on cross-module messages. The initial package namespace is `phanes.<area>.v1`.

## `phanes.common.v1`

Shared primitives:

```proto
message Version {
  uint32 major = 1;
  uint32 minor = 2;
  uint32 patch = 3;
  string prerelease = 4;
}

message FileDigest {
  string algorithm = 1; // initially sha256
  string hex = 2;
}

message LocalPath {
  string path = 1;
}

message Problem {
  string code = 1;
  string message = 2;
  Severity severity = 3;
  repeated string remediation = 4;
}

enum Severity {
  SEVERITY_UNSPECIFIED = 0;
  SEVERITY_INFO = 1;
  SEVERITY_WARNING = 2;
  SEVERITY_ERROR = 3;
  SEVERITY_FATAL = 4;
}
```

## `phanes.cache.v1`

Generated local resource-cache manifest:

```proto
message CacheManifest {
  string cache_id = 1;
  phanes.common.v1.Version schema_version = 2;
  string game_region = 3;
  string game_version = 4;
  repeated CacheSource sources = 5;
  repeated CacheArtifact artifacts = 6;
  repeated CacheIndex indexes = 7;
  google.protobuf.Timestamp created_at = 8;
}

message CacheSource {
  string source_id = 1;
  SourceKind kind = 2;
  string display_name = 3;
  string local_path = 4;
  phanes.common.v1.FileDigest digest = 5;
}

enum SourceKind {
  SOURCE_KIND_UNSPECIFIED = 0;
  SOURCE_KIND_LOCAL_INSTALL = 1;
  SOURCE_KIND_LOCAL_ARCHIVE = 2;
  SOURCE_KIND_EMBEDDED_MINIMAL = 3;
  SOURCE_KIND_EXTERNAL_IMPORT = 4;
}

message CacheArtifact {
  string artifact_id = 1;
  ArtifactKind kind = 2;
  string relative_path = 3;
  uint64 size_bytes = 4;
  phanes.common.v1.FileDigest digest = 5;
}

enum ArtifactKind {
  ARTIFACT_KIND_UNSPECIFIED = 0;
  ARTIFACT_KIND_SQLITE_DB = 1;
  ARTIFACT_KIND_BLOB_PACK = 2;
  ARTIFACT_KIND_MANIFEST = 3;
  ARTIFACT_KIND_INDEX = 4;
}

message CacheIndex {
  string name = 1;
  string version = 2;
  uint64 record_count = 3;
}
```

`CacheManifest` does not contain a digest of itself. Cache completion and manifest integrity are represented by external files such as `.cache-complete` and, if needed, `manifest.pb.sha256`, avoiding self-referential hashing rules.

## `phanes.builder.v1`

Builder input, output, verification, and progress:

```proto
message BuildRequest {
  string output_dir = 1;
  repeated InputSource input_sources = 2;
  BuildMode mode = 3;
  bool allow_external_import = 4;
}

message InputSource {
  string source_id = 1;
  phanes.cache.v1.SourceKind kind = 2;
  string local_path = 3;
}

enum BuildMode {
  BUILD_MODE_UNSPECIFIED = 0;
  BUILD_MODE_EMBEDDED_MINIMAL = 1;
  BUILD_MODE_LOCAL_CACHE = 2;
  BUILD_MODE_VERIFY_ONLY = 3;
}

message BuildResult {
  bool success = 1;
  phanes.cache.v1.CacheManifest manifest = 2;
  repeated phanes.common.v1.Problem problems = 3;
}

message VerifyCacheRequest {
  string cache_dir = 1;
  bool deep_hash_check = 2;
}

message VerifyCacheResult {
  bool valid = 1;
  phanes.cache.v1.CacheManifest manifest = 2;
  repeated phanes.common.v1.Problem problems = 3;
}

message BuildProgress {
  string phase = 1;
  uint32 percent = 2;
  string current_item = 3;
  repeated phanes.common.v1.Problem problems = 4;
}
```

## `phanes.runtime.v1`

Runtime config, status, and lifecycle:

```proto
message RuntimeConfig {
  string bind_host = 1; // default: 127.0.0.1; non-loopback invalid
  uint32 http_port = 2;
  string cache_dir = 3;
  string save_db_path = 4;
  string log_dir = 5;
  RuntimeMode mode = 6;
}

enum RuntimeMode {
  RUNTIME_MODE_UNSPECIFIED = 0;
  RUNTIME_MODE_OFFLINE_LOCAL = 1;
  RUNTIME_MODE_DIAGNOSTIC = 2;
}

message RuntimeStatus {
  RuntimeState state = 1;
  string bind_host = 2;
  uint32 http_port = 3;
  string cache_id = 4;
  string active_profile_id = 5;
  repeated phanes.common.v1.Problem problems = 6;
  google.protobuf.Timestamp started_at = 7;
}

enum RuntimeState {
  RUNTIME_STATE_UNSPECIFIED = 0;
  RUNTIME_STATE_STARTING = 1;
  RUNTIME_STATE_READY = 2;
  RUNTIME_STATE_DEGRADED = 3;
  RUNTIME_STATE_STOPPING = 4;
  RUNTIME_STATE_STOPPED = 5;
  RUNTIME_STATE_FAILED = 6;
}

message StartSessionRequest { string profile_id = 1; }
message StartSessionResult {
  string session_id = 1;
  google.protobuf.Timestamp started_at = 2;
}
message StopRuntimeRequest { string reason = 1; }
message StopRuntimeResult { bool accepted = 1; }

message HealthzResponse {
  bool ok = 1;
  RuntimeState state = 2;
}

message DiagnosticsResponse {
  repeated phanes.common.v1.Problem problems = 1;
}

message StopSessionRequest {
  string session_id = 1;
  string reason = 2;
}

message StopSessionResult { bool accepted = 1; }
```

## `phanes.launcher.v1`

Launcher process supervision and log events:

```proto
message LauncherConfig {
  string runtime_binary_path = 1;
  string builder_binary_path = 2;
  string workspace_dir = 3;
  string cache_dir = 4;
  string save_db_path = 5;
  phanes.runtime.v1.RuntimeConfig runtime = 6;
}

message RuntimeProcessSpec {
  string executable = 1;
  repeated string args = 2;
  map<string, string> env = 3;
  string working_dir = 4;
}

message RuntimeProcessStatus {
  uint32 pid = 1;
  phanes.runtime.v1.RuntimeStatus runtime_status = 2;
  ProcessState process_state = 3;
}

enum ProcessState {
  PROCESS_STATE_UNSPECIFIED = 0;
  PROCESS_STATE_NOT_STARTED = 1;
  PROCESS_STATE_STARTING = 2;
  PROCESS_STATE_RUNNING = 3;
  PROCESS_STATE_EXITED = 4;
  PROCESS_STATE_FAILED = 5;
}

message LogEvent {
  google.protobuf.Timestamp time = 1;
  string source = 2;
  string level = 3;
  string message = 4;
  repeated phanes.common.v1.Problem problems = 5;
}
```

`RuntimeProcessSpec.env` is for explicit child-process environment overrides only. It must not bypass validated runtime config. Security- or network-relevant behavior must be represented in typed config fields instead of hidden environment variables. Launcher-set variables should be allowlisted before implementation.

## `phanes.save.v1`

Initial save/profile identity only:

```proto
message Profile {
  string profile_id = 1;
  string display_name = 2;
  google.protobuf.Timestamp created_at = 3;
  google.protobuf.Timestamp updated_at = 4;
}

message SaveMetadata {
  string save_id = 1;
  string profile_id = 2;
  uint32 schema_version = 3;
  google.protobuf.Timestamp last_played_at = 4;
}
```

## Rules

- Do not create one giant `phanes.proto`.
- Do not use `map<string, any>`-style escape hatches for core contracts.
- Add examples/fixtures and validation before implementing consumers.
- Version contract packages instead of changing semantics silently.
