// Package frpx provides an FRP (Fast Reverse Proxy) V2 client for exposing
// local services behind NAT/firewall through a standard frps server.
//
// # Quick Start
//
// Create a client, register proxies, and run:
//
//	client := frpx.NewClient(
//	    frpx.WithServerAddr("frps.example.com:7000"),
//	    frpx.WithToken("your-secret-token"),
//	    frpx.WithProxy(frpx.ProxyConfig{
//	        Name:      "web",
//	        Type:      "tcp",
//	        LocalAddr: "127.0.0.1:8080",
//	        RemotePort: 6000,
//	    }),
//	)
//	client.Run(ctx)
//
// External clients can now reach your local :8080 via frps.example.com:6000.
//
// # Proxy Types
//
// Three proxy types are supported:
//
//   - "tcp": raw TCP forwarding. Requires RemotePort.
//   - "http": HTTP reverse proxy. Requires CustomDomains (Host-based routing).
//   - "https": HTTPS reverse proxy. Requires CustomDomains.
//
//	client := frpx.NewClient(
//	    frpx.WithServerAddr("frps:7000"),
//	    frpx.WithToken("token"),
//	    frpx.WithProxy(frpx.ProxyConfig{
//	        Name:      "api",
//	        Type:      "http",
//	        LocalAddr: "127.0.0.1:3000",
//	        CustomDomains: []string{"api.example.com"},
//	    }),
//	)
//
// # Authentication
//
// Authentication uses a token derived privilege key:
//
//	privilege_key = hex(md5(token + strconv.FormatInt(timestamp, 10)))
//
// The token is configured via WithToken() and must match the frps server.
//
// # Reconnection
//
// The client automatically reconnects on connection loss with exponential
// backoff. All registered proxies are re-registered after each reconnection.
//
//   - WithReconnectBackoff: initial backoff duration (default 2s)
//   - WithHeartbeatInterval: ping interval (default 30s)
//
// # Wire Protocol
//
// This client implements FRP V2 protocol with AEAD-GCM encryption:
//
//  1. Magic bytes exchange
//  2. ClientHello / ServerHello handshake
//  3. HKDF key derivation (HMAC-SHA256)
//  4. AEAD AES-256-GCM encrypted stream
//
// See docs/decisions/adr-001-frp-v2.md for the protocol decision record.
//
// # Service Lifecycle
//
// Client implements srvx.Service (Init/Run/Destroy):
//
//	svcs := &srvx.Services{}
//	svcs.Register(client)
//	svcs.Run(ctx)
package frpx
