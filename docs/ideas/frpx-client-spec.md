# Spec: frpx FRP Client

## Objective

Implement a Go library (`transport/frpx`) that acts as an FRP (Fast Reverse Proxy) client, compatible with the
standard frps server (github.com/fatedier/frp). The client exposes local services behind NAT/firewall to the public
internet through frps.

**User:** Go developers embedding reverse proxy client capability in their services.

**Success looks like:**
- `client.RegisterProxy(tcpProxy)` → remote port available on frps, traffic forwarded to local service
- `client.RegisterProxy(httpProxy)` → HTTP requests to frps with matching Host header routed to local service
- Connection loss → automatic reconnection with backoff
- `client.Run(ctx)` blocks until `ctx.Done()`, properly cleans up

## Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.25 |
| Wire protocol | FRP V1 binary framing + JSON body | compatible with frp v0.x |
| Crypto | XOR-based V1 crypto (SHA256) | standard library |
| Testing | `github.com/stretchr/testify` + mock frps | v1.11.1 |
| Lifecycle | `github.com/BYT0723/go-tools/srvx` | (internal) |

**No new external dependencies.** All crypto (SHA256, MD5) from `crypto/` standard library. JSON from `encoding/json`.
Binary framing from `encoding/binary`.

## Commands

```
Build:     go build ./transport/frpx/...
Test:      go test -race ./transport/frpx/...
Vet:       go vet ./transport/frpx/...
CI:        go test -race -count=1 -timeout 120s $(go list ./...)
```

## Project Structure

```
transport/frpx/
├── client.go         → Client struct, NewClient, srvx.Service (Name/Init/Run/Destroy),
│                        connect(), login(), heartbeatWorker(), mainLoop()
├── option.go         → Option type, WithServerAddr, WithToken, WithUser, WithHeartbeatInterval,
│                        WithProxy
├── message.go        → Message type definitions and V1 codec:
│                        Frame encode/decode (4B length + 1B type + JSON body)
│                        Login, LoginResp, NewProxy, NewProxyResp,
│                        NewWorkConn, ReqWorkConn, StartWorkConn, Ping, Pong
├── proxy.go          → ProxyConfig (TCP/HTTP/HTTPS), RegisterProxy(), handleTCPProxy(),
│                        handleHTTPProxy(), handleWorkConn
├── workconn.go       → Work connection lifecycle: ReqWorkConn handler,
│                        dialWorkConn(), forwardWorkConn()
├── crypto.go         → V1 XOR crypto reader/writer (salt exchange + per-direction key derivation)
├── client_test.go    → Unit tests for messages/options, integration tests with mock frps
├── mockfrps.go       → Minimal FRP server for testing (handles V1 handshake + basic protocol)
└── doc.go            → Package documentation
```

**Rationale:** Separation follows sshx pattern — `option.go` for configuration, `forward.go`/`workconn.go` for
distinct concerns, single `client_test.go` for all tests. `message.go` isolates the wire protocol layer.

## Protocol Details

### V1 Wire Framing

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                  Total Length (uint32, Big Endian)             |  = 1 (type) + len(body)
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|     Type (byte)      |
+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
|        JSON Body (variable length)                            |
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Read:** Read 4 bytes → total length, Read 1 byte → type, Read (totalLen - 1) bytes → JSON body.

**Write:** Marshal JSON body, write 4-byte length (1 + len(body)), write 1-byte type, write body.

### V1 XOR Crypto

On connect, both sides wrap the connection with XOR crypto:

1. Client creates XOR writer: generates 16-byte random salt, on first write prepends salt then XORs all
   data with `SHA256(token + salt)[:16]` in a stream cipher fashion
2. Server creates XOR reader: reads 16-byte salt from connection, derives key with `SHA256(token + salt)[:16]`,
   XORs all data with the key
3. Server also creates XOR writer with its own salt (client reads that salt, reverses direction)

**Result:** Client→Server data XOR'd with `SHA256(token + client_salt)`, Server→Client data XOR'd with
`SHA256(token + server_salt)`.

### Message Types (V1 type codes)

| Type | Code | Struct | Direction | Purpose |
|------|------|--------|-----------|---------|
| `'o'` (0x6F) | TypeLogin | `Login` | C→S | Authenticate and establish session |
| `'1'` (0x31) | TypeLoginResp | `LoginResp` | S→C | Login result (run_id, error) |
| `'p'` (0x70) | TypeNewProxy | `NewProxy` | C→S | Register a new proxy |
| `'2'` (0x32) | TypeNewProxyResp | `NewProxyResp` | S→C | Proxy registration result |
| `'c'` (0x63) | TypeCloseProxy | `CloseProxy` | C→S | Unregister a proxy |
| `'r'` (0x72) | TypeReqWorkConn | `ReqWorkConn{}` | S→C | Server requests a new work connection |
| `'w'` (0x77) | TypeNewWorkConn | `NewWorkConn` | C→S | Client opens work connection |
| `'s'` (0x73) | TypeStartWorkConn | `StartWorkConn` | S→C | Server binds work conn to proxy |
| `'h'` (0x68) | TypePing | `Ping` | C→S | Heartbeat |
| `'4'` (0x34) | TypePong | `Pong` | S→C | Heartbeat response |

### Authentication (Token)

```
privilege_key = hex(md5(token + strconv.FormatInt(timestamp, 10)))
timestamp     = time.Now().Unix()
```

Login message includes: version, hostname, os, arch, user, privilege_key, timestamp, pool_count.

### Proxy Registration

After login, send `NewProxy` messages on the control connection for each proxy. Key fields:

| Field | TCP | HTTP | HTTPS |
|-------|-----|------|-------|
| `proxy_type` | `"tcp"` | `"http"` | `"https"` |
| `remote_port` | Required | — | — |
| `custom_domains` | — | Required | Required |
| `subdomain` | — | Optional | Optional |
| `host_header_rewrite` | — | Optional | — |
| `http_user` / `http_pwd` | — | Optional | — |

### Work Connection Flow (TCP)

```
External Client              frps                      frpx (us)
     |                         |                          |
     |--- TCP to :6000 ------>|                          |
     |                         |--- ReqWorkConn ('r') -->|  (control conn)
     |                         |                          |--- open new TCP ----->|
     |                         |               |<--- NewWorkConn ('w') -------|
     |                         |               |--- StartWorkConn ('s') ----->|
     |                         |                          |  (work conn: proxyName)
     |                         |                          |--- dial local svc ->|
     |<======================== bidirectional pipe ===========================>|
```

### HTTP Work Connection Flow

```
External Client              frps                      frpx (us)
     |                         |                          |
     |--- HTTP GET / --------->|                          |
     |   Host: app.example.com |                          |
     |                         |--- ReqWorkConn ('r') -->|
     |                         |   ... same flow as TCP ...
     |                         |                          |--- proxy HTTP ---> local :8080
     |<====================== HTTP response ================================>|
```

## Code Style

Follow existing project conventions (matching sshx):

```go
// Exported types
type Client struct {
    mu       sync.Mutex
    addr     string
    token    string
    user     string
    running  bool
    proxies  map[string]*proxyConfig
    wg       sync.WaitGroup
}

// Functional options
type Option func(*Client)

func WithServerAddr(addr string) Option {
    return func(c *Client) { c.addr = addr }
}

// srvx.Service interface
func (c *Client) Name() string                      { return "frpx" }
func (c *Client) Init(ctx context.Context) error     { return nil }
func (c *Client) Run(ctx context.Context) error      { ... }
func (c *Client) Destroy(ctx context.Context) error  { ... }
```

- No comments unless critical
- Exported types: `Client`, `Option`, `ProxyConfig`
- Unexported: handler functions, message types, crypto helpers
- Thread safety via `sync.Mutex` on shared state
- All goroutines tracked via `sync.WaitGroup`

## Testing Strategy

- **Framework:** `stretchr/testify/assert` + Go standard `testing`
- **Test location:** In-package (`package frpx`), file `client_test.go`
- **Mock server:** `mockfrps.go` — a minimal FRP server that:
  1. Accepts TCP connections
  2. Wraps with XOR crypto
  3. Handles Login → LoginResp
  4. Handles NewProxy → NewProxyResp
  5. Handles ReqWorkConn / StartWorkConn flow
  6. Responds to Ping with Pong
- **Integration tests:** Client connects to mock frps, registers proxy, verifies traffic flows
- **Test levels:**
  - Unit: Message encode/decode, option behavior, auth key computation
  - Integration: Full client→mock frps lifecycle (login, proxy, ping, reconnect)

**No real network endpoints.** All tests use mock frps on loopback.

## Boundaries

### Always Do
- Run `go test -race ./transport/frpx/...` before declaring a step complete
- Follow existing naming conventions (`TestXxx`, `WithXxx`, `TypeXxx`)
- Validate inputs in public functions (addr, token, proxy config)
- Close all connections and goroutines on Destroy

### Ask First
- Adding external dependencies
- Changing the public API after initial release
- Modifying CI configuration

### Never Do
- Commit secrets or tokens
- Panic in library code (return errors)
- Rely on real network or real frps in tests
- Remove existing tests

## Success Criteria

- [ ] Client connects to mock frps, logs in, registers TCP proxy
- [ ] External TCP connection to frps proxy port is forwarded to local service
- [ ] Client connects to mock frps, registers HTTP proxy with custom domain
- [ ] HTTP request with matching Host header reaches local HTTP service
- [ ] Heartbeat ping/pong works; connection timeout triggers reconnection
- [ ] Client stops gracefully (no goroutine leaks, all connections closed)
- [ ] `go test -race ./transport/frpx/...` passes with zero races

## Open Questions

- **Proxy removal**: Should `Client` support `UnregisterProxy(name string)`? Yes — `CloseProxy` message.
- **Reconnection**: Keep existing proxy registrations on reconnect? Yes — re-register all proxies after login.
- **Work connection pool**: Should we pre-create work connections? Defer — one-at-a-time for MVP.
- **HTTP host header rewrite**: Client-side or server-side? Client-side (proxy config sets rewrite).
