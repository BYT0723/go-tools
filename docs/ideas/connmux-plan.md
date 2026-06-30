# Implementation Plan: connmux — Connection Multiplexer

## Overview

Build `transport/connmux`: a single-port multi-protocol multiplexer that sniffs TCP connections and dispatches them to protocol-specific `net.Listener` instances. Implements `srvx.Service` for unified lifecycle management with the rest of the repo.

**New package:** `transport/connmux/` (5 files)
**Modified package:** `transport/sshx/` (+2 methods: SetListener, Match)

## Architecture Decisions

- **Mux implements srvx.Service, but doesn't manage sub-service lifecycles internally.** Mux owns its real listener + VirtualListeners. Sub-services (sshx, http.Server wrapper) are also srvx.Service instances, but their lifecycle is wired by the caller via `srvx.Services` — not by Mux itself.
- **`Route()` API, auto-match.** `mux.Route(name, svc)` registers a `ListenedService`. Mux calls `svc.Match()` to get the matcher, creates a `VirtualListener`, and passes it via `svc.SetListener(vl)`. The service declares its protocol, not the caller.
- **Channel-based VirtualListener.** Uses `chan net.Conn` + `chan struct{}` for shutdown signaling. No polling, no timers.
- **Byte replay via `replayConn`.** Wraps `io.MultiReader(buf, rawConn)` to prepend sniffed bytes. Sub-service sees the complete byte stream.
- **First-match-wins ordering.** Matchers are evaluated in `Listen()` registration order. Custom matchers have priority over built-ins if registered first.

## Dependency Graph

```
VirtualListener (listener.go)     Matcher (matcher.go)
        │                                │
        └────────────┬───────────────────┘
                     │
              Mux core (mux.go + option.go + doc.go)
                     │
        ┌────────────┴────────────┐
        │                         │
  Integration tests          sshx.WithListener
  (connmux_test.go)          (sshx/option.go + server.go)
```

## Task List

### Phase 1: Foundation

---

#### Task 1: Matcher interface and built-in matchers

**Description:** Define the `Matcher` interface, implement `BytePrefixMatcher`, `DefaultMatcher`, and export built-in matcher instances (`MatchSSH`, `MatchHTTP1`). Pure logic, no network I/O.

**Acceptance criteria:**
- [ ] `Matcher` interface with `Match(sniffed []byte) bool` method
- [ ] `BytePrefixMatcher` matches bytes with given prefix, is case-sensitive
- [ ] `DefaultMatcher` always returns true
- [ ] `MatchSSH` matches `SSH-` prefix
- [ ] `MatchHTTP1` matches all 9 HTTP/1.x method prefixes: `GET `, `POST`, `HEAD`, `PUT `, `DELE`, `OPTI`, `PATC`, `CONN`, `TRAC`
- [ ] `MatchHTTP1` correctly rejects non-HTTP (e.g. `GETTER`, binary data)

**Verification:**
- [ ] `go test -race -run TestMatcher ./transport/connmux/...` passes

**Dependencies:** None

**Files likely touched:**
- `transport/connmux/matcher.go` (new)
- `transport/connmux/doc.go` (new)

**Estimated scope:** XS (1 new file)

---

#### Task 2: VirtualListener — channel-based net.Listener

**Description:** Implement `VirtualListener` implementing `net.Listener` (Accept, Close, Addr). Uses a buffered `chan net.Conn` for incoming connections and a `chan struct{}` for shutdown signaling. `Push(conn)` non-blocking on closed listener.

**Acceptance criteria:**
- [ ] `Accept()` returns connections from the channel, or `ErrListenerClosed` when closed
- [ ] `Close()` is idempotent; subsequent `Accept()` calls return `ErrListenerClosed`
- [ ] `Addr()` returns the advertised address
- [ ] `Push(conn)` succeeds when listener is open
- [ ] `Push(conn)` drops silently (select with default) when listener is closed
- [ ] Concurrent `Accept()` + `Close()` does not race

**Verification:**
- [ ] `go test -race -run TestVirtualListener ./transport/connmux/...` passes

**Dependencies:** None (Task 1 not required, but both run in same phase)

**Files likely touched:**
- `transport/connmux/listener.go` (new)

**Estimated scope:** S (1 new file + tests in connmux_test.go)

---

### Checkpoint: Foundation
- [ ] `go test -race ./transport/connmux/...` passes for matchers and VirtualListeners
- [ ] `go vet ./transport/connmux/...` clean

---

### Phase 2: Core Mux

---

#### Task 3: Mux struct, Route API, lifecycle, and accept loop

**Description:** Implement `Mux` struct with `srvx.Service` interface, `Route(name, svc)` API, accept loop with sniff+dispatch, and full lifecycle management. Mux derives its own context for sub-services and manages Init→Run→Destroy ordering correctly.

**Lifecycle ordering (critical):**
- Run: deriveCtx → real listener → SetListener → Init (rollback on failure) → Run goroutines → serve goroutine
- Stop: close real listener → wait serve exit → cancel deriveCtx → wait sub-service goroutines → Destroy → close VLs
- Non-blocking push: `select { vl.ch <- conn; <-vl.closed; default: conn.Close() }`

**Acceptance criteria:**
- [ ] `ListenedService` interface defined: embeds `srvx.Service`, adds `SetListener(net.Listener)` and `Match() Matcher`
- [ ] `Route(name string, svc ListenedService)` registered before `Start()`; panics if called after
- [ ] `Route()` calls `svc.Match()` to auto-derive the matcher, creates `VirtualListener`, calls `svc.SetListener(vl)`
- [ ] `Init(ctx)` validates: sniffSize > 0 and < 65536, at least one route registered
- [ ] `Run(ctx)` derives internal context; parent ctx only triggers shutdown
- [ ] Partial Init failure destroys previously-init'd services and returns error
- [ ] `serve()` accepts, sets `SetReadDeadline(5s)`, sniffs, matches in order, dispatches
- [ ] Non-blocking VL push: drops connection if VL closed or channel full
- [ ] No-match fallback: last route receives unmatched connections
- [ ] `Stop()` closing order: real listener → serve exit → cancel sub ctx → wg.Wait → Destroy (reverse order) → close VLs
- [ ] `Name()` returns `"connmux"`
- [ ] `WithAddr(addr)` and `WithSniffSize(n)` options work

**Verification:**
- [ ] `go test -race -run TestMux ./transport/connmux/...` passes
- [ ] Manual: start/stop cycle, verify goroutines exit

**Dependencies:** Task 1 (matcher), Task 2 (VirtualListener)

**Files likely touched:**
- `transport/connmux/mux.go` (new)
- `transport/connmux/option.go` (new)
- `transport/connmux/doc.go` (append)

**Estimated scope:** M (3 new files)

---

### Checkpoint: Core Mux
- [ ] `go test -race ./transport/connmux/...` passes for all unit tests
- [ ] Mux lifecycle (Init→Run→Destroy) works with no goroutine leaks
- [ ] `go vet ./transport/connmux/...` clean

---

### Phase 3: sshx Integration

---

#### Task 4: sshx: implement ListenedService (SetListener + Match)

**Description:** Add `SetListener(net.Listener)` method and `Match() connmux.Matcher` method to `sshx.Server`. Modify `Start()` to skip listener creation when external listener is set. All existing sshx tests must pass unchanged.

**Acceptance criteria:**
- [ ] `sshx.Server` satisfies `connmux.ListenedService` via `SetListener(l net.Listener)` + `Match() Matcher` returning `connmux.MatchSSH`
- [ ] sshx `Start()` uses `SetListener`-injected listener instead of creating one; creates its own when not set
- [ ] sshx `Stop()` closes the listener regardless of source
- [ ] All existing sshx tests pass: `go test -race ./transport/sshx/...`

**Verification:**
- [ ] `go test -race ./transport/sshx/...` full pass, zero regressions
- [ ] Manual: create sshx with `WithListener(vl)`, verify SSH client can connect via VL

**Dependencies:** Task 3 (Mux provides VirtualListener for testing)

**Files likely touched:**
- `transport/sshx/option.go`
- `transport/sshx/server.go` (Start())

**Estimated scope:** S (2 files, ~15 lines changed)

---

### Checkpoint: sshx Integrated
- [ ] All existing sshx tests pass
- [ ] sshx can accept connections from a VirtualListener

---

### Phase 4: Integration Tests & Polish

---

#### Task 5: Integration tests

**Description:** Write integration tests covering the full stack: Mux + HTTP server, Mux + multi-HTTP (two servers), byte replay correctness, lifecycle (Start→Stop→re-Start), unmatched fallback.

**Acceptance criteria:**
- [ ] `TestMuxHTTP` — HTTP client connects to Mux, reaches correct HTTP handler
- [ ] `TestMuxMultiHTTP` — two HTTP matchers (different prefixes), connections routed correctly
- [ ] `TestMuxByteReplay` — verify sniffed bytes reach the sub-service intact (send known payload, read it back)
- [ ] `TestMuxLifecycle` — Start→Stop→re-Start, verify cleanup and restart
- [ ] `TestMuxUnmatched` — connection that matches nothing, verify default handler or connection close
- [ ] `TestMuxConcurrent` — multiple concurrent connections, verify no races

**Verification:**
- [ ] `go test -race ./transport/connmux/...` passes with all tests
- [ ] Coverage >80% on mux.go, listener.go, matcher.go

**Dependencies:** Task 3, Task 4

**Files likely touched:**
- `transport/connmux/connmux_test.go` (new)

**Estimated scope:** M (1 new file, ~5 test functions)

---

#### Task 6: CI and documentation

**Description:** Update CI workflow to include `transport/connmux` in build/test matrix, update AGENTS.md and README.md.

**Acceptance criteria:**
- [ ] `go build ./transport/connmux/...` succeeds
- [ ] `go test -race ./transport/connmux/...` included in CI (or confirmed already covered by `./...`)
- [ ] AGENTS.md module list includes connmux
- [ ] README.md module table includes connmux

**Verification:**
- [ ] `go build ./...` succeeds
- [ ] `go vet ./...` clean
- [ ] README and AGENTS.md updated

**Dependencies:** Task 5

**Files likely touched:**
- `AGENTS.md`
- `README.md`

**Estimated scope:** XS (2 files, minor edits)

---

### Final Checkpoint: Complete
- [ ] All 6 tasks done
- [ ] `go test -race ./transport/connmux/...` passes (zero races)
- [ ] `go test -race ./transport/sshx/...` passes (no regression)
- [ ] `go build ./...` + `go vet ./...` clean
- [ ] All spec success criteria met

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| `replayConn` data race under concurrent Reads | Med | Each connection is handed to exactly one Accept() caller; `io.MultiReader` is not safe for concurrent reads. Document as part of the contract. |
| Sub-service goroutine blocks on its own listener when Mux is shutting down | Med | Mux cancels deriveCtx → sub-service Run() returns → wg.Wait → then VL is closed after services exit. VL push is non-blocking during serve loop. |
| `SetReadDeadline` causes spurious timeout under real network latency | Low | 5s deadline is generous for initial bytes. Configurable via `WithSniffDeadline`. |
| `io.ReadAtLeast` with min=1 on tiny payloads (e.g., client sends only 1 byte) matches incorrectly | Low | Matchers check `len(sniffed)` before comparing. MatchSSH needs 4 bytes; 1 byte won't match and falls through to default. |

## Open Questions

- **What happens when no matcher matches?** Spec says fall back to last-registered listener (which should be `DefaultMatcher`). If the user never registered a default, the connection is closed with a log warning. Confirm this is acceptable.
- **Should `Listen()` return an error instead of panicking when called after Start?** Panic is simpler and catches programmer error early. Consistent with Go's `http.Server.ListenAndServe` panicking on misconfiguration. But we could return error for a gentler API. Spec says panic for now.
