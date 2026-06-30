// Package connmux multiplexes multiple network protocols on a single TCP port.
// It sniffs the initial bytes of incoming connections to identify the protocol
// and dispatches each connection to the matching service via VirtualListener.
//
// Services implement the ListenedService interface (srvx.Service + SetListener + Match).
// Mux manages the full lifecycle of all routed services.
package connmux
