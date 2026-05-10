package ssh

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"sync"

	xssh "golang.org/x/crypto/ssh"
)

var ErrServerClosed = fmt.Errorf("ssh server closed")

type (
	Handler func(sess *xssh.ServerConn, ch xssh.Channel, reqs <-chan *xssh.Request)

	Server struct {
		mu       sync.Mutex
		addr     string
		config   *xssh.ServerConfig
		listener net.Listener
		running  bool
		handler  Handler
	}
)

func NewServer(opts ...Option) *Server {
	s := &Server{
		addr:   ":2222",
		config: &xssh.ServerConfig{},
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server already running")
	}
	s.running = true
	s.mu.Unlock()

	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", s.addr)
	if err != nil {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		return err
	}

	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()

	go s.serve(ctx)
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}
	s.mu.Unlock()
	return nil
}

func (s *Server) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.addr
}

func (s *Server) serve(ctx context.Context) {
	for {
		conn, err := s.accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			s.mu.Lock()
			stop := !s.running
			s.mu.Unlock()
			if stop {
				return
			}
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) accept() (net.Conn, error) {
	s.mu.Lock()
	l := s.listener
	s.mu.Unlock()
	if l == nil {
		return nil, ErrServerClosed
	}
	return l.Accept()
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	sshConn, chans, reqs, err := xssh.NewServerConn(conn, s.config)
	if err != nil {
		return
	}
	defer sshConn.Close()

	go xssh.DiscardRequests(reqs)

	for newCh := range chans {
		if newCh.ChannelType() != "session" {
			newCh.Reject(xssh.UnknownChannelType, "only session channels supported")
			continue
		}

		ch, reqs, err := newCh.Accept()
		if err != nil {
			continue
		}

		if s.handler != nil {
			go s.handler(sshConn, ch, reqs)
		} else {
			go defaultHandler(ch, reqs)
		}
	}
}

func defaultHandler(ch xssh.Channel, reqs <-chan *xssh.Request) {
	defer ch.Close()

	go func() {
		for req := range reqs {
			switch req.Type {
			case "shell":
				req.Reply(true, nil)
			case "exec":
				payload := struct{ Cmd string }{}
				xssh.Unmarshal(req.Payload, &payload)
				ch.Write([]byte(payload.Cmd + ": command not found\n"))
				req.Reply(true, nil)
				ch.SendRequest("exit-status", false, xssh.Marshal(&struct{ Status uint32 }{1}))
				return
			default:
				req.Reply(false, nil)
			}
		}
	}()

	io.Copy(io.Discard, ch)
}

func GenerateHostKey() ([]byte, error) {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	b := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(k),
	}
	return pem.EncodeToMemory(b), nil
}
