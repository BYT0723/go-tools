// Package connmux multiplexes multiple network protocols on a single TCP port.
// It sniffs the initial bytes of incoming connections to identify the protocol
// and dispatches each connection to the matching service via VirtualListener.
//
// # Quick Start
//
// The core workflow: create a Mux, create protocol services, route them, and run:
//
//	mux := connmux.NewMux(connmux.WithAddr(":443"))
//
//	sshSrv := sshx.NewServer(sshx.WithHostKey(key), sshx.WithUser("admin", "pass"))
//
//	mux.Route("ssh", sshSrv)                            // SSH connections matched by "SSH-" prefix
//	mux.Route("http", httpx.WrapHTTPServer(&http.Server{...})) // HTTP via GET/POST/... prefix
//
//	srvx.Services{}.Register(mux).Run(ctx)
//
// For Echo or Gin, wrap the framework in an http.Server:
//
//	connmux.WrapHTTPServer(&http.Server{Handler: echoInstance})
//	connmux.WrapHTTPServer(&http.Server{Handler: ginEngine})
//
// # Protocol Detection
//
// Mux reads the first N bytes of each connection (default 256) and evaluates
// each route's Matcher in registration order. The first match wins. The last
// registered route acts as a catch-all (typically using MatchDefault).
//
// Built-in matchers:
//   - MatchSSH: matches "SSH-" prefix (SSH-2.0 banners)
//   - MatchHTTP1: matches HTTP/1.x method prefixes (GET, POST, HEAD, PUT, DELETE,
//     OPTIONS, PATCH, CONNECT, TRACE)
//   - MatchHTTP2: matches "PRI * HTTP/2.0" prefix (HTTP/2 connection preface, gRPC)
//   - MatchTLS: matches {0x16, 0x03} prefix (TLS handshake record)
//   - MatchDefault: matches everything (catch-all)
//
// Custom matchers implement the Matcher interface:
//
//	type myMatcher struct{}
//	func (myMatcher) Match(sniffed []byte) bool { return bytes.HasPrefix(sniffed, []byte("MYPROTO")) }
//
// # ListenedService Interface
//
// Services routed through Mux must implement ListenedService:
//
//	type ListenedService interface {
//	    Name() string
//	    Init(ctx context.Context) error
//	    Run(ctx context.Context) error
//	    Destroy(ctx context.Context) error
//	    SetListener(net.Listener)
//	    Match() Matcher
//	}
//
// Match() declares the protocol this service handles.
// SetListener() receives the VirtualListener that Mux will push connections to.
// Init/Run/Destroy follow the srvx.Service lifecycle, managed automatically by Mux.Run().
//
// # HTTP Service Adapter
//
// Use httpx.WrapHTTPServer to wrap an *http.Server into ListenedService:
//
//	mux.Route("http", httpx.WrapHTTPServer(&http.Server{Handler: handler}))
//
// Echo and Gin are supported by passing the framework as the handler:
//
//	httpx.WrapHTTPServer(&http.Server{Handler: echoInstance})
//	httpx.WrapHTTPServer(&http.Server{Handler: ginEngine})
//
// # SSH Service Integration
//
// sshx.Server implements ListenedService directly via SetListener() and Match() methods:
//
//	sshSrv := sshx.NewServer(sshx.WithHostKey(key), sshx.WithPasswordAuth(authFn))
//	mux.Route("ssh", sshSrv) // sshx.Match() returns connmux.MatchSSH
//
// # Lifecycle
//
// Mux.Run() manages the full lifecycle of all routed services:
//
//  1. Creates the real TCP listener
//  2. Calls SetListener() on each service (injects VirtualListener)
//  3. Calls Init() on each service (rolls back on failure)
//  4. Starts Run() goroutines for all successfully inited services
//  5. Starts the accept/sniff/dispatch serve loop
//  6. On context cancellation: stops the real listener, stops serve, cancels
//     sub-service contexts, waits for Run() to return, calls Destroy()
package connmux
