# AGENTS.md

This repository is developed with **opencode** and **oh-my-opencode-slim**.

oh-my-opencode-slim already provides the agent roles and task-dispatch style. Do **not** create a second role system in this file. Use the existing opencode skills, commands, and MCP servers to plan and execute work.

This file only explains what Phanes is trying to become, the hard boundaries, and the shared collaboration tools available in this environment.

---

## 1. Project Name

**Phanes**

Phanes is not a Grasscutter fork and should not be treated as a direct rewrite of Grasscutter.

The name intentionally separates this project from the old Grasscutter ecosystem. Grasscutter, Cultivation, and similar projects are useful references, but Phanes should be designed as a new project with its own architecture.

---

## 2. What We Want to Build

Phanes is a **local-first offline runtime, resource-cache builder, and launcher framework** for sandbox research.

The intended user experience:

```text
User opens Phanes Launcher
  ↓
Launcher checks local configuration and local resource cache
  ↓
Launcher starts a localhost-only runtime
  ↓
The game connects only to local services
  ↓
Runtime uses local save data and local resource cache
  ↓
When the game exits, Phanes shuts down cleanly
```

The goal is to make the experience feel lightweight and local/offline.

Internally, Phanes may still use a local runtime service. The goal is not to turn the original client into a fully standalone executable by deep patching.

---

## 3. Desired Result

A successful Phanes implementation should eventually provide:

```text
local runtime
local launcher
local account/save data
local resource cache
resource-cache builder
typed protocol contracts
clear diagnostics
one-command or one-click local startup
```

It should avoid requiring:

```text
MongoDB
GC-Resources
remote resource repositories
manual management of multiple heavy services
public server deployment
```

---

## 4. Hard Boundaries

These rules are mandatory.

### 4.1 Localhost by Default

The runtime must default to:

```text
127.0.0.1
```

Do not bind to `0.0.0.0` by default.

Do not add public server hosting features unless the human maintainer explicitly changes the project scope.

### 4.2 Do Not Depend on GC-Resources

GC-Resources is not a required dependency.

Phanes should not require users to download GC-Resources or any remote resource repository.

Preferred resource model:

```text
embedded-minimal  # tiny test dataset
local-cache       # generated local cache; primary target
external-import   # optional compatibility/import mode
```

### 4.3 Runtime Must Not Download Resources

Runtime startup should use local files only.

If resources need to be prepared, use a builder step before runtime execution.

### 4.4 Do Not Bundle Copyrighted Game Resources

Phanes should not redistribute full copyrighted game resources.

Local cache generation should be designed around user-provided local inputs and normalized local cache output.

### 4.5 No Unsafe Bypass Work

Do not implement or document:

```text
anti-cheat bypass
client protection bypass
official service authentication bypass
unauthorized online interaction
commercial service replacement
deep client patching instructions
```

If a task appears to require one of these, stop and ask the human maintainer.

---

## 5. Important Design Direction

Use a **contract-first** workflow.

Before writing implementation code for cross-module behavior, define the contract first.

Contracts may include:

```text
protobuf schemas
Go interfaces
SQLite schemas
CLI command contracts
config file formats
runtime status API
builder input/output format
launcher/runtime communication format
```

Preferred workflow:

```text
1. Define or update the contract
2. Add examples or fixtures
3. Add tests/validation for the contract
4. Implement against the contract
5. Change contracts only through explicit discussion
```

Do not let modules communicate through undocumented JSON shapes, untyped maps, implicit structs, or accidental assumptions.

Use Go and protobuf strengths wherever useful.

Prefer:

```text
protobuf messages
generated Go code
typed Go interfaces
versioned schemas
explicit adapters
small stable contracts
```

Avoid for core boundaries:

```text
map[string]any
large global structs
reflection-heavy registration
silent schema drift
ad-hoc undocumented JSON
```

JSON is fine for config, logs, UI-facing data, and import/export where appropriate.

---

## 6. Major Components

Phanes will likely need these conceptual components. opencode may choose the exact repository structure.

```text
Runtime:
  localhost-only local service
  dispatch-like HTTP endpoints
  local session lifecycle
  save-data access
  resource-cache access
  diagnostics

Builder:
  local source detection
  resource extraction/normalization
  resource-cache generation
  resource-cache verification

Launcher:
  desktop control panel
  runtime process control
  resource build workflow
  local config UI
  logs and diagnostics

Protocol:
  protobuf definitions
  generated Go types
  packet/session abstractions
  version adapters

Storage:
  SQLite save data
  SQLite or equivalent resource cache
  migrations
  backup/export/import boundaries
```

Keep save data and resource cache conceptually separate:

```text
save data:
  user-owned, persistent, should not be casually destroyed

resource cache:
  generated, rebuildable, read-mostly
```

---

## 7. Suggested Tech Preferences

These are preferences, not a rigid plan.

```text
Runtime/CLI:
  Go

Protocol:
  protobuf + generated Go code

Storage:
  SQLite

Launcher:
  Tauri 2

Frontend:
  Svelte, React, or another opencode-chosen stack

Docs:
  Markdown

Collaboration:
  GitHub issues/PRs through gh
```

opencode may propose better choices when justified.

---

## 8. Available Tools

The local environment has useful tools available.

### 8.1 gh

`gh` is installed.

Use GitHub as the shared coordination surface when useful.

Recommended use:

```text
Issues:
  task handoff
  contract proposal
  blocker
  design question
  review request

Pull Requests:
  implementation changes
  documentation changes
  schema/proto changes

Issue/PR comments:
  progress updates
  review feedback
  handoff notes
```

Before scripting gh commands, check local help because flags may differ by version:

```bash
gh --help
gh issue --help
gh pr --help
```

Good patterns:

```bash
gh issue list
gh issue create
gh issue comment <id>
gh pr list
gh pr create
gh pr comment <id>
```

Use concise GitHub comments for handoffs:

```markdown
## Handoff

Current state:
- ...

Blocked on:
- ...

Needed from next agent:
- ...

Validation:
- ...
```

Do not rely on hidden local reasoning as project memory. If another agent needs to know it, write it into an issue, MR, ADR, or repo doc.

### 8.2 rtk

`rtk` may be available or installable.

Reference: https://github.com/rtk-ai/rtk

rtk is useful for reducing noisy command output before it enters the agent context.

Use it when inspecting large outputs, such as:

```bash
rtk git status
rtk git diff
rtk git log -n 10
rtk ls .
rtk read <file>
rtk grep <pattern> .
rtk go test ./...
rtk err <command>
rtk test <command>
```

Before relying on a command, check:

```bash
rtk --help
rtk --version
```

If rtk is not available or behaves unexpectedly, fall back to normal shell commands.

### 8.3 opencode Skills and MCP Servers

Use available opencode skills and MCP servers.

The project intentionally leaves planning and implementation strategy to opencode. Use skills/MCP to:

```text
inspect repository structure
search references
draft contracts
generate protobufs
coordinate implementation
create GitHub issues/PRs
run tests
review changes
```

Do not ignore existing opencode coordination conventions from oh-my-opencode-slim.

---

## 9. Reference Links

Use these references to understand the background. They are references, not dependencies.

### Old server reference

Grasscutter:

```text
https://github.com/Grasscutters/Grasscutter
```

Useful for understanding historical concepts:

```text
dispatch
local account flow
session
packet handler
resource loader
game systems
```

Do not copy code without explicit license review.

### Old launcher/proxy reference

Cultivation:

```text
https://github.com/Grasscutters/Cultivation
```

Useful for understanding launcher responsibilities:

```text
local launcher
runtime/server startup
proxy-like routing ideas
process monitoring
user configuration
```

Do not blindly copy behavior. Phanes should avoid unsafe client patching and should not silently mutate system settings.

### Resource reference only

GC-Resources:

```text
https://gitlab.com/YuukiPS/GC-Resources
```

This is reference-only.

Phanes must not require GC-Resources.

### Tooling reference

rtk:

```text
https://github.com/rtk-ai/rtk
```

Use for compact command output when useful.

### Suggested technology references

```text
Go:
  https://go.dev/

Protocol Buffers:
  https://protobuf.dev/

Tauri:
  https://tauri.app/

SQLite:
  https://www.sqlite.org/

GitHub CLI:
  https://cli.github.com/
```

---

## 10. Documentation Expectations

Keep docs focused on:

```text
what Phanes should achieve
what boundaries must not be crossed
what contracts exist
how modules communicate
what assumptions were made
how to validate behavior
```

Avoid:

```text
overly rigid roadmaps
premature implementation plans
duplicating oh-my-opencode-slim agent roles
large task lists inside AGENTS.md
```

When an important cross-module decision is made, prefer an ADR:

```text
docs/adr/0001-contract-first.md
docs/adr/0002-resource-cache-format.md
docs/adr/0003-runtime-status-contract.md
```

Simple ADR shape:

```markdown
# ADR 000N: Title

## Status

Proposed | Accepted | Superseded

## Context

## Decision

## Consequences

## Related
```

---

## 11. Validation Expectations

Changes should generally be validated with the relevant checks.

Examples:

```bash
go test ./...
go test ./internal/...
go vet ./...
```

If rtk is available:

```bash
rtk go test ./...
rtk err go test ./...
rtk git diff
```

For launcher work, use the package manager selected by opencode and document the command.

Before committing large changes, inspect the diff:

```bash
git diff
# or
rtk git diff
```

---

## 12. Current Accepted Decisions

These decisions are currently accepted:

```text
Project name: Phanes
Project type: local-first offline runtime + builder + launcher
Relationship to Grasscutter: reference only, not a fork
Runtime language preference: Go
Protocol preference: protobuf
Storage preference: SQLite
Launcher preference: Tauri 2
Resource strategy: local cache, no GC-Resources dependency
Default runtime binding: 127.0.0.1
Coordination: opencode + oh-my-opencode-slim + gh
Output reduction tool: rtk when useful
Planning/design: delegated to opencode agents
AGENTS.md purpose: context and boundaries only
```

Agents may propose changes, but should not silently change these assumptions.

---

## 13. When to Ask the Human Maintainer

Ask before doing any of the following:

```text
changing project scope
adding public server support
adding client patching
adding anti-cheat/protection bypass logic
copying code from GPL/AGPL projects
bundling or redistributing game resources
making GC-Resources required
changing the license
changing from local-first to network-first
```

---

## 14. Short Version

Build **Phanes** as a new local-first project.

Use:

```text
Go
protobuf
SQLite
Tauri
local resource cache
gh for GitHub coordination
rtk for compact command output
opencode/oh-my-opencode-slim for planning and task execution
```

Do not build:

```text
a Grasscutter fork
a public private-server framework
a GC-Resources-dependent project
a client bypass/patching project
a copyrighted resource bundle
```

Define contracts first, then write code.
