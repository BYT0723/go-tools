# Transport

Network transport utilities for Go services.

## Packages

| Package                               | Description                                              |
| ------------------------------------- | -------------------------------------------------------- |
| [sshx](#sshx)                         | Embeddable SSH server (PTY shell, exec, port forwarding) |
| [frpx](#frpx)                         | FRP V2 client (expose local services behind NAT)         |
| [connmux](#connmux)                   | Single-port multi-protocol connection multiplexer        |
| [httpx](httpx/)                       | HTTP client wrapper (encoder, decoder, compressor)       |
| [httpx/middleware](httpx/middleware/) | Gin and Echo middleware (logging, trace context)         |

---

## sshx

Embeddable SSH server with PTY shell, remote command execution, and TCP port forwarding.

### Quick start

```go
key, _ := sshx.GenerateHostKey()

srv := sshx.NewServer(
    sshx.WithAddr(":2222"),
    sshx.WithHostKey(key),
    sshx.WithUser("admin", "password"),
)

srv.Start(ctx)
```

### Auth

```go
// Password
sshx.WithUser("admin", "password")

// Custom password validation
sshx.WithPasswordAuth(func(user, password string) bool {
    return user == "admin" && checkHash(password, storedHash)
})

// Public key
sshx.WithPublicKeyAuth(func(conn ssh.ConnMetadata, key ssh.PublicKey) bool {
    return isAuthorizedKey(key)
})
```

### Capabilities

```bash
ssh user@host -p 2222                    # interactive PTY shell
ssh user@host -p 2222 ls -la             # remote exec
ssh -L 8080:localhost:80 user@host -p 2222  # local port forward
ssh -R 9999:localhost:22 user@host -p 2222  # remote port forward
```

### Lifecycle

```go
srv := sshx.NewServer(...)
svcs := &srvx.Services{}
svcs.Register(srv)
svcs.Run(ctx)
```

---

## frpx

FRP (Fast Reverse Proxy) V2 client. Exposes local services to the internet through a standard frps server.

### Quick start

```go
client := frpx.NewClient(
    frpx.WithServerAddr("frps.example.com:7000"),
    frpx.WithToken("your-secret-token"),
    frpx.WithProxy(frpx.ProxyConfig{
        Name:       "web",
        Type:       "tcp",
        LocalAddr:  "127.0.0.1:8080",
        RemotePort: 6000,
    }),
)
client.Run(ctx)
// External clients can now reach :8080 via frps:6000
```

### Proxy types

```go
// TCP
frpx.WithProxy(frpx.ProxyConfig{
    Name: "ssh", Type: "tcp", LocalAddr: "127.0.0.1:22", RemotePort: 2222,
})

// HTTP (host-based routing)
frpx.WithProxy(frpx.ProxyConfig{
    Name: "api", Type: "http", LocalAddr: "127.0.0.1:3000",
    CustomDomains: []string{"api.example.com"},
})

// HTTPS
frpx.WithProxy(frpx.ProxyConfig{
    Name: "secure", Type: "https", LocalAddr: "127.0.0.1:8443",
    CustomDomains: []string{"secure.example.com"},
})
```

### Reconnection

The client auto-reconnects on connection loss:

```go
frpx.WithReconnectBackoff(5 * time.Second)    // initial backoff (default 2s)
frpx.WithHeartbeatInterval(30 * time.Second)  // ping interval (default 30s)
```

### Wire protocol (V2)

```
Magic(7B) → ClientHello(JSON) → ServerHello(JSON) → AEAD-GCM stream
                        ↑ HKDF(HMAC-SHA256, token, transcript)
```

### Lifecycle

```go
client := frpx.NewClient(...)
svcs := &srvx.Services{}
svcs.Register(client)
svcs.Run(ctx)
```

---

## connmux

Single-port multi-protocol connection multiplexer. Sniffs initial bytes to detect the protocol, then dispatches each connection to the correct service.

### Quick start

```go
hostKey, _ := sshx.GenerateHostKey()

mux := connmux.NewMux(connmux.WithAddr(":443"))

// SSH
sshSrv := sshx.NewServer(
    sshx.WithHostKey(hostKey),
    sshx.WithUser("admin", "pass"),
)

// HTTP (stdlib, Echo, or Gin wrapped in http.Server)
httpSrv := &http.Server{Handler: handler}

mux.Route("ssh", sshSrv)                                 // detected by "SSH-" prefix
mux.Route("http", httpx.WrapHTTPServer(httpSrv))         // detected by "GET"/"POST"/... prefix

svcs := &srvx.Services{}
svcs.Register(mux)
svcs.Run(ctx)
```

### How it works

```
                   :443 TCP listener
                           │
                     ┌─────▼─────┐
                     │    Mux    │
                     │   sniff   │
                     │   match   │
                     └─┬──────┬──┘
                       │      │
              "SSH-2.0"│      │"GET /..."
                    ┌──▼───┐ ┌▼──────┐
                    │ sshx │ │ http  │
                    │server│ │server │
                    └──────┘ └───────┘
```

### Built-in matchers

| Matcher        | Prefix                                                                 | Example               |
| -------------- | ---------------------------------------------------------------------- | --------------------- |
| `MatchSSH`     | `SSH-`                                                                 | `SSH-2.0-OpenSSH_9.0` |
| `MatchHTTP1`   | `GET `, `POST`, `HEAD`, `PUT `, `DELE`, `OPTI`, `PATC`, `CONN`, `TRAC` | `GET / HTTP/1.1`      |
| `MatchDefault` | _anything_                                                             | catch-all             |

### Custom matcher

```go
type myMatcher struct{}

func (myMatcher) Match(sniffed []byte) bool {
    return bytes.HasPrefix(sniffed, []byte("MYPROTO"))
}
```

### HTTP service adapter

Use `httpx.WrapHTTPServer()` to wrap any `*http.Server`:

```go
mux.Route("http", httpx.WrapHTTPServer(&http.Server{Handler: handler}))
```

For Echo or Gin, wrap the framework in an `*http.Server`:

```go
// Echo
e := echo.New()
mux.Route("http", httpx.WrapHTTPServer(&http.Server{Handler: e}))

// Gin
g := gin.New()
mux.Route("http", httpx.WrapHTTPServer(&http.Server{Handler: g}))
```

### Lifecycle

Mux.Run() manages the full lifecycle of all routed services:

1. Creates the real TCP listener
2. Calls `SetListener()` on each service (injects VirtualListener)
3. Calls `Init()` on each service (rolls back on failure)
4. Starts `Run()` goroutines for successfully inited services
5. Starts accept/sniff/dispatch serve loop
6. On context cancellation: stops listener, stops serve, cancels sub-services, waits, calls `Destroy()`

### Integration with sshx

sshx.Server implements `connmux.ListenedService` directly — no adapter needed:

```go
sshSrv := sshx.NewServer(...)
mux.Route("ssh", sshSrv) // sshx.Match() returns connmux.MatchSSH
```
