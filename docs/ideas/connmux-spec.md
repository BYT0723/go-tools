# Spec: connmux — Connection Multiplexer

## Objective

Build a Go library (`transport/connmux`) that accepts TCP connections on a single port, sniffs the initial bytes
to detect the protocol, and dispatches each connection to the correct protocol-specific handler. The target
architecture is **Mux as srvx.Service** owning the real listener and creating `VirtualListener` instances
that downstream services (e.g. SSH server, HTTP server) consume via standard `net.Listener.Accept()` —
these services don't know they're being multiplexed.

**User:** Go developers who need to serve multiple protocols on one port (e.g. HTTP + SSH on :443 to bypass firewalls that block :22).

**Success looks like:**
- Spin up one `connmux.Mux` on port 8080, register two backends (HTTP and SSH), and clients of both protocols connect successfully.
- `ssh user@host -p 8080` → SSH session works.
- `curl http://host:8080` → HTTP response works.
- `Mux.Destroy()` closes everything cleanly.
- All sub-services are standard `srvx.Service` instances; no custom interfaces needed.

## Tech Stack

| Component | Technology | Version / Source |
|-----------|-----------|------------------|
| Language | Go | 1.25 |
| Multi-service lifecycle | `github.com/BYT0723/go-tools/srvx` | (internal) |
| SSH | `github.com/BYT0723/go-tools/transport/sshx` | (internal, to be modified) |
| HTTP | `net/http` (stdlib) | Go 1.25 |
| Testing | `github.com/stretchr/testify` | v1.11.1 |

**No new external dependencies.** Protocol sniffing uses `io.Reader` byte-peeking only. VirtualListener is a pure channel-based `net.Listener` implementation.

## Commands

```
Build:   go build ./transport/connmux/...
Test:    go test -race ./transport/connmux/...
Vet:     go vet ./transport/connmux/...
```

## Project Structure

```
transport/connmux/
├── mux.go           → Mux struct, NewMux, srvx.Service (Name/Init/Run/Destroy),
│                      Start/Stop, AddService, serve (accept loop + sniff + dispatch)
├── listener.go      → VirtualListener: channel-based net.Listener for sub-services
├── matcher.go       → Matcher type (BytePrefix, TLS SNI), DefaultMatch (catch-all)
├── option.go        → Option type, WithListener (for testing), WithSniffSize
├── connmux_test.go  → Unit tests for matchers, integration tests with HTTP + sshx
└── doc.go           → Package documentation
```

**Also modified:**

```
transport/sshx/
└── option.go        → +WithListener(net.Listener) Option — allows external listener injection
```

**Rationale:** `listener.go` isolates the channel-based VirtualListener — it's the core primitive that makes the
architecture work. `matcher.go` isolates protocol detection logic so it's independently testable and extensible.
sshx only needs a ~10-line option addition to accept an external listener.

## Architecture

### Connection Flow

```
                    ┌─────────────────────┐
                    │   Real net.Listener  │
                    │   e.g. :443          │
                    └────────┬────────────┘
                             │ Accept()
                             ▼
                    ┌─────────────────────┐
                    │   Mux.serve()       │
                    │                     │
                    │   1. Accept conn    │
                    │   2. Peek N bytes   │
                    │   3. Match protocol │
                    │   4. Prepend bytes  │
                    │      + conn         │
                    │   5. Push to        │
                    │      VirtualListener│
                    └──┬──────┬──────┬────┘
                       │      │      │
              ┌────────┘      │      └────────┐
              ▼               ▼               ▼
      ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
      │VirtualLis   │ │VirtualLis   │ │VirtualLis   │
      │"SSH-2.0..." │ │"GET/POST..."│ │ (default)   │
      └─────┬───────┘ └─────┬───────┘ └─────┬───────┘
            │               │               │
      ┌─────▼──────┐  ┌─────▼──────┐  ┌─────▼──────┐
      │ sshx Server│  │ http Server│  │ fallback    │
      │ (Accept)   │  │ (Accept)   │  │ handler     │
      └────────────┘  └────────────┘  └────────────┘
```

None of the sub-services call `net.Listen`. They all receive connections through
`VirtualListener.Accept()` — they think they're serving on a regular port.

### VirtualListener

A `VirtualListener` is a channel-based `net.Listener`:

```go
type VirtualListener struct {
    ch     chan net.Conn
    closed chan struct{}
    addr   net.Addr
    mu     sync.Mutex
}

func (vl *VirtualListener) Accept() (net.Conn, error) {
    select {
    case conn, ok := <-vl.ch:
        if !ok {
            return nil, ErrListenerClosed
        }
        return conn, nil
    case <-vl.closed:
        return nil, ErrListenerClosed
    }
}

func (vl *VirtualListener) Close() error {
    vl.mu.Lock()
    defer vl.mu.Unlock()
    select {
    case <-vl.closed:
        return nil
    default:
        close(vl.closed)
    }
    return nil
}

func (vl *VirtualListener) Addr() net.Addr {
    return vl.addr
}
```

**Push-side:** Mux's `serve()` goroutine calls `push(conn)` which does a select
on `vl.ch <- conn` or `<-vl.closed`. When the VirtualListener is closed (its
parent service stopped), push returns without blocking.

### Byte Replay

Sniffing reads N bytes from the raw connection. Those bytes must be prepended
before the connection reaches the sub-service. The connection pushed into the
VirtualListener is wrapped:

```go
wrapped := &replayConn{
    Conn:   rawConn,
    reader: io.MultiReader(bytes.NewReader(sniffed), rawConn),
}

func (c *replayConn) Read(b []byte) (int, error) {
    return c.reader.Read(b)
}
```

This ensures the sub-service sees the complete byte stream, including the
sniffed protocol-identifying bytes (e.g. `SSH-2.0-...`).

### Matcher System

```go
type Matcher interface {
    Match(sniffed []byte) bool
}

type BytePrefixMatcher struct {
    Prefix []byte
}

func (m *BytePrefixMatcher) Match(sniffed []byte) bool {
    return bytes.HasPrefix(sniffed, m.Prefix)
}

// DefaultMatcher matches everything (catch-all).
type DefaultMatcher struct{}
func (DefaultMatcher) Match([]byte) bool { return true }
```

**Built-in matchers (pre-defined constants):**
- `MatchSSH` — prefix `SSH-` (matches `SSH-2.0-...`)
- `MatchHTTP1` — prefixes `GET `, `POST`, `HEAD`, `PUT `, `DELE`, `OPTI`, `PATC`, `CONN`, `TRAC`
- `MatchDefault` — always matches (catch-all, last in order)

**Custom matchers:** Users can provide their own `Matcher` via `WithMatcher`.

### Mux Lifecycle

```
Mux implements srvx.Service
  Init(ctx)  → validates config: sniffSize > 0 && <= 65536, at least one route registered
  Run(ctx)   → calls runOnce(ctx)
  Destroy(ctx) → calls stop(ctx, backgroundCtx)

// Mux derives an internal context from parent — parent triggers shutdown,
// internal context runs sub-services. This prevents sub-services from
// seeing a cancelled parent context during normal operation.
Internal flow (runOnce):
  1. deriveCtx, cancel := context.WithCancel(parentCtx)
     defer cancel()
  2. Creates real net.Listener on mux.addr
  3. For each route (in order): svc.SetListener(vl) → svc.Init(deriveCtx)
     If any Init fails: destroy previously-init'd services, return error
  4. Spawns goroutine per svc: svc.Run(deriveCtx) → wg.Done()
  5. Spawns serve(deriveCtx) accept loop
  6. <-parentCtx.Done()
  7. Close real listener (stops serve loop)
  8. Wait for serve goroutine to exit
  9. cancel() (signals sub-services to stop)
  10. wg.Wait() (all sub-service Run calls have returned)
  11. svc.Destroy(shutdownCtx) on each sub-service (reverse order)
  12. Close all VirtualListeners
```

### Serve Loop Details

```
serve(ctx):
  for {
    conn, err := listener.Accept()
    if err → return
    // Sniff with deadline to prevent slow-loris:
    conn.SetReadDeadline(time.Now().Add(5s))
    sniffed := make([]byte, m.sniffSize)
    n, _ := io.ReadAtLeast(conn, sniffed, 1)  // at least 1 byte
    sniffed = sniffed[:n]
    // Match against routes in registration order:
    for _, route := range m.routes {
      if route.matcher.Match(sniffed) { push(route.vl, wrap(conn, sniffed)); continue }
    }
    // No match: last route acts as fallback (usually MatchDefault)
    push(m.routes[len(m.routes)-1].vl, wrap(conn, sniffed))
  }

// Push is non-blocking — if VL is closed/draining, connection is dropped
push(vl, conn):
  select {
  case vl.ch <- conn:
  case <-vl.closed:
    conn.Close()  // service no longer accepting
  default:
    conn.Close()  // channel full, drop
  }
```

### ListenedService Interface

```go
type ListenedService interface {
    srvx.Service
    SetListener(net.Listener)
    Match() Matcher
}
```

Services that want to be routed by Mux implement this interface. Mux calls `Match()` to discover the protocol, `SetListener()` to inject the `VirtualListener`, then manages the service's full lifecycle via the embedded `srvx.Service` (Init → Run → Destroy).

### sshx Change: implement ListenedService

Add two methods to `sshx.Server`:

```go
func (s *Server) SetListener(l net.Listener) {
    s.mu.Lock()
    s.listener = l
    s.mu.Unlock()
}

func (s *Server) Match() connmux.Matcher {
    return connmux.MatchSSH
}
```

### Usage Example

```go
mux := connmux.NewMux(connmux.WithAddr(":443"))

sshSrv := sshx.NewServer(
    sshx.WithHostKey(key),
    sshx.WithPublicKeyAuth(authCallback),
)

httpSrv := connmux.HTTPService(handler)  // adapter: http.Handler → ListenedService

mux.Route("ssh", sshSrv)   // Mux calls sshSrv.Match() → MatchSSH
mux.Route("http", httpSrv) // Mux calls httpSrv.Match() → MatchHTTP1

// Mux runs everything — sub-services are started/stopped by Mux itself
srvx.Services(mux).Run(ctx)
```

## Code Style

Follow existing project conventions (matching sshx/frpx):

```go
type Option func(*Mux)

func WithAddr(addr string) Option { return func(m *Mux) { m.addr = addr } }

type ListenedService interface {
    srvx.Service
    SetListener(net.Listener)
    Match() Matcher
}

type Mux struct {
    mu       sync.Mutex
    addr     string
    listener net.Listener
    running  bool
    routes   []*serviceRoute
    wg       sync.WaitGroup
}

type serviceRoute struct {
    name    string
    matcher Matcher
    vl      *VirtualListener
    svc     ListenedService
}
```

- No comments unless critical
- Exported: `Mux`, `Option`, `ListenedService`, `VirtualListener`, `Matcher`, `MatchSSH`, `MatchHTTP1`, `MatchDefault`, `ErrListenerClosed`
- Unexported: `serviceRoute`, `replayConn`, `BytePrefixMatcher`, `DefaultMatcher`
- Functional options pattern
- Thread safety via `sync.Mutex`
- All goroutines tracked via `sync.WaitGroup`

## Testing Strategy

- **Framework:** `stretchr/testify/assert` + Go standard `testing`
- **Test location:** In-package (`package connmux`), file `connmux_test.go`
- **Test levels:**
  - Unit: Matcher matching, VirtualListener accept/close, option behavior
  - Integration:
    - HTTP + HTTP (two HTTP servers on one port, different routes matched by matcher)
    - HTTP + nothing (matcher that never matches, verify default fallback)
    - Byte replay correctness (verify sniffed bytes arrive at sub-service)
    - Lifecycle: Start/Stop, goroutine cleanup, re-start
- **Mock SSH:** Use `golang.org/x/crypto/ssh` test utilities for integration test with sshx
- **No real network endpoints** beyond loopback
- **Coverage target:** >80% on mux.go, listener.go, matcher.go

## Boundaries

### Always Do
- Run `go test -race ./transport/connmux/...` before declaring a step complete
- Follow existing naming conventions (`TestXxx`, `WithXxx`)
- Validate inputs in public functions
- Close all listeners, channels, and goroutines on Destroy

### Ask First
- Adding external dependencies
- Changing the sshx public API (WithListener is additive, not breaking)
- Modifying CI configuration

### Never Do
- Commit secrets
- Panic in library code (return errors)
- Rely on real network endpoints (except loopback) in tests
- Remove existing tests

## Success Criteria

- [ ] `connmux.Mux` on single port dispatches SSH connections to sshx.Server
- [ ] `connmux.Mux` on single port dispatches HTTP connections to http.Server
- [ ] Sniffed bytes are correctly replayed to sub-service (SSH handshake works)
- [ ] VirtualListener.Close() unblocks Accept() and stops push operations
- [ ] Mux.Stop() cleans up all goroutines, real listener, and VirtualListeners
- [ ] sshx.WithListener() allows external listener injection without breaking existing tests
- [ ] `go test -race ./transport/connmux/...` passes with zero races
- [ ] `go test -race ./transport/sshx/...` still passes (no regression)

## Open Questions

- **Sub-service lifecycle ordering:** Should Mux.Run() start sub-services, or leave that to the caller via `srvx.Services`? Recommendation: leave it to the caller. Mux owns its own listener + VirtualListners only; sub-service lifecycle is caller's responsibility.
- **ReplayConn concurrency:** Can the same replayConn be Read by two goroutines simultaneously? No — our VirtualListener guarantees each connection is handed to exactly one Accept() caller, so no concurrent read issue.
- **Matcher ordering:** What happens if two matchers both match the same sniffed bytes (e.g. custom matcher vs default)? Matchers are evaluated in registration order, first match wins. This is explicitly documented behavior.
- **TLS:** Should Mux handle TLS termination before sniffing? No — TCP payload is unencrypted at sniff time. If the user wants HTTPS + SSH over TLS, they need to handle TLS separately (e.g. via SNI routing in a reverse proxy). This is out of scope for connmux.
- **sniff size:** Should it be configurable? Yes, via `WithSniffSize(n int)`. Default: 256 bytes (enough for SSH banner + HTTP methods).
