// Package tlsx provides TLS utilities for connmux.
//
// # Server
//
// WrapTLS wraps a TLS config and connection handler into connmux.ListenedService.
// Connectons are decrypted via TLS handshake, then passed to the handler as a
// plaintext net.Conn. Use with any protocol — HTTP, gRPC, custom TCP.
//
//	mux.Route("tls", tlsx.WrapTLS(tlsCfg, func(conn net.Conn) {
//	    httpSrv.Serve(&tlsx.SingleConnListener{Conn: conn})
//	}))
package tlsx
