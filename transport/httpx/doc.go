// Package httpx provides HTTP transport utilities for Go services.
//
// # Client
//
// httpx wraps net/http.Client with encoder/decoder/compressor support.
// See client.go and option.go for the HTTP client API.
//
// # Server
//
// WrapHTTPServer wraps an *http.Server into connmux.ListenedService,
// allowing it to be multiplexed with other protocols on a single port:
//
//	mux := connmux.NewMux(connmux.WithAddr(":443"))
//	mux.Route("http", httpx.WrapHTTPServer(&http.Server{Handler: handler}))
//
// For Echo or Gin, wrap the framework as the handler:
//
//	httpx.WrapHTTPServer(&http.Server{Handler: echoInstance})
//	httpx.WrapHTTPServer(&http.Server{Handler: ginEngine})
//
// # Middleware
//
// Sub-packages provide middleware for popular frameworks:
//   - middleware/ginx: trace/span/request ID injection + API logging for Gin
//   - middleware/echox: trace/span/request ID injection + API logging for Echo
package httpx
