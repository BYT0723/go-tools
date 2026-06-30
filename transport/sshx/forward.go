package sshx

import (
	"fmt"
	"io"
	"net"
	"sync"

	xssh "golang.org/x/crypto/ssh"
)

func (s *Server) handleGlobalRequests(sshConn *xssh.ServerConn, reqs <-chan *xssh.Request) {
	var (
		mu   sync.Mutex
		fwds = make(map[string]net.Listener)
	)

	for req := range reqs {
		switch req.Type {
		case "tcpip-forward":
			var payload struct {
				Address string
				Port    uint32
			}
			if err := xssh.Unmarshal(req.Payload, &payload); err != nil {
				if req.WantReply {
					req.Reply(false, nil)
				}
				continue
			}

			laddr := net.JoinHostPort(payload.Address, fmt.Sprintf("%d", payload.Port))
			listener, err := net.Listen("tcp", laddr)
			if err != nil {
				if req.WantReply {
					req.Reply(false, nil)
				}
				continue
			}

			actualPort := uint32(listener.Addr().(*net.TCPAddr).Port)
			key := net.JoinHostPort(payload.Address, fmt.Sprintf("%d", actualPort))

			mu.Lock()
			fwds[key] = listener
			mu.Unlock()

			if req.WantReply {
				req.Reply(true, xssh.Marshal(&struct{ Port uint32 }{actualPort}))
			}

			go s.acceptForwarded(sshConn, listener, payload.Address, actualPort)

		case "cancel-tcpip-forward":
			var payload struct {
				Address string
				Port    uint32
			}
			if err := xssh.Unmarshal(req.Payload, &payload); err != nil {
				if req.WantReply {
					req.Reply(false, nil)
				}
				continue
			}

			key := net.JoinHostPort(payload.Address, fmt.Sprintf("%d", payload.Port))
			mu.Lock()
			ln, ok := fwds[key]
			if ok {
				delete(fwds, key)
				ln.Close()
			}
			mu.Unlock()

			if req.WantReply {
				req.Reply(ok, nil)
			}

		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}

	mu.Lock()
	for _, ln := range fwds {
		ln.Close()
	}
	mu.Unlock()
}

func (s *Server) acceptForwarded(sshConn *xssh.ServerConn, listener net.Listener, addr string, port uint32) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		originAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
		if !ok {
			conn.Close()
			continue
		}

		ch, reqs, err := sshConn.OpenChannel("forwarded-tcpip", xssh.Marshal(&struct {
			ConnectedAddr string
			ConnectedPort uint32
			OriginAddr    string
			OriginPort    uint32
		}{
			ConnectedAddr: addr,
			ConnectedPort: port,
			OriginAddr:    originAddr.IP.String(),
			OriginPort:    uint32(originAddr.Port),
		}))
		if err != nil {
			conn.Close()
			continue
		}

		go xssh.DiscardRequests(reqs)
		go func() {
			defer ch.Close()
			defer conn.Close()
			pipeConn(ch, conn)
		}()
	}
}

func (s *Server) handleDirectTcpip(newCh xssh.NewChannel) {
	var payload struct {
		DestAddr   string
		DestPort   uint32
		OriginAddr string
		OriginPort uint32
	}
	if err := xssh.Unmarshal(newCh.ExtraData(), &payload); err != nil {
		newCh.Reject(xssh.ConnectionFailed, "invalid payload")
		return
	}

	ch, reqs, err := newCh.Accept()
	if err != nil {
		return
	}
	defer ch.Close()
	go xssh.DiscardRequests(reqs)

	addr := net.JoinHostPort(payload.DestAddr, fmt.Sprintf("%d", payload.DestPort))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	pipeConn(ch, conn)
}

func pipeConn(ch xssh.Channel, conn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(ch, conn)
	}()
	go func() {
		defer wg.Done()
		_, _ = io.Copy(conn, ch)
		ch.CloseWrite()
	}()
	wg.Wait()
}
