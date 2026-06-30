# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a collection of Go utility modules (`github.com/BYT0723/go-tools`). Each module is a standalone package providing common functionality for Go applications: configuration loading, logging, monitoring, data structures, transport utilities, and more.

The project is structured as a library, not an application. The `cmd/go-tool` directory contains a small test binary demonstrating usage of weak pointers.

## Common Commands

### Building and Testing
- `go build ./...` - Build all packages
- `go test ./...` - Run all tests
- `go test -race ./...` - Run tests with race detector
- `go test -count=1 -timeout 120s $(go list ./... | grep -v -E 'monitor/(ping|web|prometheus|multi|snmp)')` - CI test command
- `go test -coverprofile=cover.cov ./...` - Generate coverage report
- `go test ./ds` - Test a specific package

### Dependency Management
- `go mod tidy` - Clean up dependencies
- `go mod download` - Download modules
- `go list ./...` - List all packages

### Lint & Vet
- `go vet ./...` - Static analysis (no external linter; golangci-lint was removed due to false positives)

## Architecture

### Module Structure
Each top-level directory is a separate Go module (package) with its own responsibility:

- **cfg**: Configuration loading using viper, supports local files and remote etcd
- **channelx**: Channel utilities
- **contextx**: Context utilities (trace ID, request ID, logger injection, etc.)
- **ds**: Data structures (counters, mutexes, pools, caches, queues, sets, maps, hubs, stacks)
- **funny**: Experimental/graph utilities
- **i18n**: Internationalization wrappers
- **logx**: Logging facade with implementations for zap and zerolog
- **monitor**: Monitoring framework (core types, ping, Prometheus, SNMP, web, multi)
- **transport**: Transport utilities (HTTP/httpx, SSH)
- **osx**: OS utilities
- **mathx**: Math utilities
- **spider**: Web scraping utilities
- **srvx**: Service lifecycle management (Init/Run/Destroy)
- **unsafex**: Unsafe operations with benchmarks

### Key Patterns

1. **Configuration**: Centralized in `cfg` package. Use `cfg.Init()` with options. Tags `cfg:"field"`.
2. **Logging**: Use `logx` facade; implementations in `logx/zaplogger` and `logx/zerologger`.
3. **Monitoring**: Implement `monitor.Monitor` interface; `monitor/component` provides reusable lifecycle/alert logic.
4. **Testing**: Use **`testing.T.Run` + `github.com/stretchr/testify/assert`** (goconvey was fully removed). Tests are in-package (`package ds` not `package ds_test`).
5. **Networking tests**: All tests MUST use `httptest.Server` or mock data; never rely on real network (ICMP/HTTP/prometheus endpoints).
6. **Data Structures**: Thread-safe and non-thread-safe variants (e.g., `counter` vs `mutexCounter`).

### Dependencies
Major external dependencies:
- `github.com/stretchr/testify` - testing assertions
- `github.com/spf13/viper` - configuration
- `go.uber.org/zap` and `github.com/rs/zerolog` - logging
- `github.com/prometheus-community/pro-bing` - ping monitoring
- `github.com/gin-gonic/gin` and `github.com/labstack/echo/v4` - HTTP frameworks

## Known Issues & Fixes

### Race conditions fixed (2026-05-10)
- **`ds/cache.go:Release()`** — now holds mutex when writing `ctx`/`cf` fields; cleanup goroutine captures `ctx.Done()` before loop to avoid nil deref after Release.
- **`ds/fastHub.go:Publish()`** — snapshots subscribers into a slice under lock, releases lock, then iterates the slice. Prevents race with concurrent `Unsubscribe()`.
- **`ds/cache_test.go`** — uses `c.Get()` instead of reading `c.entries` directly; expire time reduced to 50ms (was 1s) for fast tests.

### Source bugs (not fixed, noted for awareness)
- **`ds/arrayStack.Pop()` and `ds/linkStack.Pop()`** — panic on empty stack (no empty check before accessing internal slice/node)
- **`ds/linkStack.Pop()`** — does not decrement `size` field
- **`transport/httpx/encoder/json.go:RequestHeader()`** — infinite recursion bug (calls `d.RequestHeader()` instead of `d.compressor.RequestHeader()`)
- **`monitor/web/monitor.go:NewMonitor()`** — does not set `m.method` field
- **`transport/httpx/decoder/json.go:Decode()`** — decodes twice, second Decode always fails with EOF
- **`monitor/snmp/`** — incomplete/untracked code, excluded from build

## CI

GitHub Actions workflow at `.github/workflows/test.yml`:
- Trigger: push/PR to main/master
- Go 1.25, runs `go vet ./...`, `go build ./...`, `go test -race -timeout 120s` (excludes some network-dependent monitor subpackages)

## Commit Style
Conventional commits: `feat:`, `fix:`, `chore:`, `perf:`, `refactor:`, `doc:`.
