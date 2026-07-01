# ADR-002: Single-Port Multi-Protocol Connection Multiplexer

## Status
Accepted

## Date
2026-07-01

## Context

We need to serve multiple network protocols (SSH, HTTP/1.1, gRPC/HTTP/2, TLS-terminated traffic) on a single
TCP port. This is useful for:

- Bypassing firewalls that only allow port 80/443 (e.g. SSH + HTTP on :443)
- Reducing the number of open ports in containerized deployments
- Simplifying infrastructure by running one listener per host

The challenge: TCP connections to a single port from different protocols (SSH vs HTTP) carry no inherent
protocol identifier. We must detect the protocol from the initial bytes of each connection before the
application-layer handler can process it.

## Decision

Build `transport/connmux` — a connection multiplexer that **sniffs the first N bytes** of each incoming
connection and **dispatches to a protocol-specific service** via a VirtualListener.

### Architecture

```
:443 TCP listener
       │
  ┌────▼────┐
  │   Mux   │  sniff N bytes → match in route order → dispatch
  └─┬──┬──┬─┘
    │  │  │
 SSH- GET PRI*
  │  │  │
  ▼  ▼  ▼
sshx http http2
```

### Core Abstractions

**ListenedService interface** — services declare their protocol and accept a listener:

```go
type ListenedService interface {
    srvx.Service
    SetListener(net.Listener)  // injected VirtualListener
    Match() Matcher             // self-declared protocol matcher
}
```

**Route() API** — one-line registration, no manual matcher parameter:

```go
mux.Route("ssh", sshSrv)     // sshSrv.Match() → MatchSSH
mux.Route("http", httpSrv)   // httpSrv.Match() → MatchHTTP1
```

The service declares what it IS. The caller doesn't need to know protocol magic bytes.

**VirtualListener** — a channel-based `net.Listener`. Sub-services call `Accept()` just like a real TCP
listener. They don't know they're being multiplexed.

**Lifecycle** — Mux implements `srvx.Service`. `Mux.Run()` manages the full lifecycle of all routed
services: `SetListener → Init → Run → (ctx cancel) → Destroy`. Caller only needs `srvx.Services{}.Register(mux).Run(ctx)`.

### Built-in Protocol Matchers

| Matcher | Prefix | Bytes | Protocol |
|---------|--------|-------|----------|
| `MatchSSH` | `SSH-` | 4 | SSH-2.0 |
| `MatchHTTP1` | `GET `, `POST`, `HEAD`, `PUT `, `DELE`, `OPTI`, `PATC`, `CONN`, `TRAC` | 4-7 | HTTP/1.x |
| `MatchHTTP2` | `PRI * HTTP/2.0` | 16 | HTTP/2 (gRPC, web) |
| `MatchTLS` | `{0x16, 0x03}` | 2 | TLS handshake |
| `MatchDefault` | *anything* | — | Catch-all |

Sniff size defaults to 256 bytes. All prefix lengths fit comfortably within this.

### Byte Replay

Sniffed bytes are consumed from the connection but must reach the sub-service intact (SSH needs `SSH-2.0-...`,
HTTP needs `GET /...`). Each dispatched connection is wrapped in a `replayConn` using `io.MultiReader`
to prepend the sniffed bytes before the raw connection stream.

### Slow-Loris Protection

`conn.SetReadDeadline(5s)` is set before sniffing. Connections that send fewer bytes than needed within
the deadline are dropped.

## Alternatives Considered

### Single service per port (no multiplexing)
- Pros: Simpler, no protocol detection needed
- Cons: Multiple ports, firewall issues, orchestration complexity
- Rejected: The goal is single-port operation

### cmux (soheilhy/cmux) as external dependency
- Pros: Battle-tested Go connection multiplexer with built-in HTTP1/HTTP2/TLS matchers
- Cons: External dependency, different API style (listener-based registration), no srvx.Service lifecycle
- Rejected: We can implement a lighter version (~250 lines) that fits our srvx.Service + functional options patterns

### Three-parameter Route API: `Route(name, svc, matcher)`
- Pros: Explicit — caller chooses the matcher
- Cons: Redundant — the service already knows its protocol (sshx knows it speaks SSH);
  `Route("ssh", sshSrv, MatchSSH)` repeats information
- Rejected: Switched to self-declaring `Match()` on the service. Caller only writes `Route("ssh", sshSrv)`

### HTTP/1.1 and HTTP/2 as a single matcher
- Pros: `*http.Server` handles both protocols transparently; one route serves both
- Cons: Can't route gRPC and HTTP/1.1 to different services or handlers
- Rejected: Following cmux precedent — HTTP/1.1 and HTTP/2 are distinct at the wire level
  (`GET` vs `PRI * HTTP/2.0`). Users who want unified handling can route both to the same service.

### TLS termination inside connmux
- Pros: Would allow routing HTTPS vs gRPC-over-TLS by SNI/ALPN
- Cons: Requires TLS certificate management and key material inside connmux; bloats the scope
- Rejected: connmux only sniffs TLS (forwarding the encrypted stream). TLS termination belongs
  at a higher layer (reverse proxy, service mesh).

## Consequences

- **Zero external dependencies** — All matching is `bytes.HasPrefix`; VirtualListener is channels;
  lifecycle is srvx (same repo)
- **Match ordering matters** — First match wins. Put `MatchTLS` before `MatchHTTP1` if you want TLS
  traffic to hit the TLS handler instead of matching `GET`/`POST` bytes inside the encrypted stream
  (they won't — `0x16` ≠ `G`, but the principle holds)
- **Default sniff size 256 bytes** — Sufficient for all built-in matchers. Users with custom protocols
  needing larger prefixes should use `WithSniffSize(n)`
- **VirtualListener buffer 128** — Burst connections beyond buffer capacity are dropped (non-blocking push)
- **HTTP/2 matcher enables gRPC** — gRPC uses HTTP/2 transport; `MatchHTTP2` routes gRPC connections
  to a service that can handle HTTP/2 (e.g. `*http.Server` with `grpc.Server` as handler)

## Key Implementation Decisions
- Channel-based VirtualListener allows sub-services to use standard `net.Listener.Accept()`
- Mux manages sub-service lifecycle (Init → Run → Destroy) — callers don't need to wire srvx.Services
  for individual sub-services inside Mux
- Shutdown order: close real listener → wait for serve goroutine → close all VirtualListeners →
  cancel sub-service contexts → wait for services → destroy services
- `httpService.Close()` (not `Shutdown()`) for Destroy — quick termination; users wanting
  graceful shutdown should implement it in their own adapter
- `sniffDeadline` configured via `WithSniffDeadline(d)`, defaults to 5s

## See Also
- [FRP V2 ADR](adr-001-frp-v2.md)
- [connmux spec](../ideas/connmux-spec.md)
- [connmux implementation plan](../ideas/connmux-plan.md)
- [cmux (reference implementation)](https://github.com/soheilhy/cmux)
