// Package sshx provides an embeddable SSH server for Go services.
//
// # Quick Start
//
// Create a server with host key and authentication, then start it:
//
//	key, _ := sshx.GenerateHostKey()
//	srv := sshx.NewServer(
//	    sshx.WithAddr(":2222"),
//	    sshx.WithHostKey(key),
//	    sshx.WithUser("admin", "password"),
//	)
//	srv.Start(ctx)
//
// # Authentication
//
// Two auth methods are supported:
//
//   - Password: WithUser(name, password) for simple cases, or
//     WithPasswordAuth(func(user, password string) bool) for custom validation
//   - Public key: WithPublicKeyAuth(func(conn ssh.ConnMetadata, key ssh.PublicKey) bool)
//
// # Capabilities
//
// The server supports:
//
//   - PTY shell: `ssh user@host` gets an interactive bash shell
//   - Remote exec: `ssh user@host <command>` runs commands
//   - TCP forwarding (-L): `ssh -L 8080:localhost:80 user@host`
//   - TCP reverse forwarding (-R): `ssh -R 9999:localhost:22 user@host`
//
// # Service Lifecycle
//
// Server implements srvx.Service (Init/Run/Destroy):
//
//	svcs := &srvx.Services{}
//	svcs.Register(srv)
//	svcs.Run(ctx)
//
// # Integration with connmux
//
// Server implements connmux.ListenedService. It can be routed through a
// connection multiplexer to share a port with HTTP or other protocols:
//
//	mux := connmux.NewMux(connmux.WithAddr(":443"))
//	mux.Route("ssh", sshSrv)
//	mux.Route("http", httpSrv)
//	mux.Run(ctx)
//
// See transport/connmux for details.
package sshx
