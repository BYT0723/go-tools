package sshx

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/BYT0723/go-tools/transport/connmux"
	"github.com/creack/pty"
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

		ownsListener bool
		shellPath    string
		wg           sync.WaitGroup
	}
)

func NewServer(opts ...Option) *Server {
	s := &Server{
		addr:      ":2222",
		config:    &xssh.ServerConfig{},
		shellPath: "/bin/bash",
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *Server) SetListener(l net.Listener) {
	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()
}

func (s *Server) Match() connmux.Matcher {
	return connmux.MatchSSH
}

func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server already running")
	}
	s.running = true

	if s.listener != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.serve(ctx)
		}()
		s.mu.Unlock()
		return nil
	}

	s.ownsListener = true
	s.mu.Unlock()

	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", s.addr)
	if err != nil {
		s.mu.Lock()
		s.running = false
		s.ownsListener = false
		s.mu.Unlock()
		return err
	}

	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.serve(ctx)
	}()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.mu.Lock()
	s.running = false
	if s.listener != nil && s.ownsListener {
		s.listener.Close()
	}
	s.mu.Unlock()
	s.wg.Wait()
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

func (s *Server) Name() string {
	return "sshx"
}

func (s *Server) Init(ctx context.Context) error {
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.Start(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	return s.Stop(context.Background())
}

func (s *Server) Destroy(ctx context.Context) error {
	return s.Stop(ctx)
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
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handleConn(conn)
		}()
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

	go s.handleGlobalRequests(sshConn, reqs)

	for newCh := range chans {
		switch newCh.ChannelType() {
		case "session":
			ch, reqs, err := newCh.Accept()
			if err != nil {
				continue
			}
			if s.handler != nil {
				go s.handler(sshConn, ch, reqs)
			} else {
			go s.handleSession(ch, reqs)
			}
		case "direct-tcpip":
			go s.handleDirectTcpip(newCh)
		default:
			newCh.Reject(xssh.UnknownChannelType, "unknown channel type")
		}
	}
}

func (s *Server) handleSession(ch xssh.Channel, reqs <-chan *xssh.Request) {
	defer ch.Close()

	var (
		cmd     *exec.Cmd
		ptyFile *os.File
		win     = &pty.Winsize{Rows: 24, Cols: 80}
		term    = "xterm"
		ptyReq  bool
		done    = make(chan uint32, 1)
	)

	for req := range reqs {
		switch req.Type {
		case "pty-req":
			ptyReq = true
			var p struct {
				Term   string
				Cols   uint32
				Rows   uint32
				Width  uint32
				Height uint32
				Modes  string
			}
			if err := xssh.Unmarshal(req.Payload, &p); err != nil {
				req.Reply(false, nil)
				continue
			}
			term = p.Term
			if p.Cols > 0 && p.Rows > 0 {
				win.Cols = uint16(p.Cols)
				win.Rows = uint16(p.Rows)
			}
			req.Reply(true, nil)

		case "shell":
			if !ptyReq {
				req.Reply(false, nil)
				ch.SendRequest("exit-status", false, xssh.Marshal(&struct{ Status uint32 }{1}))
				return
			}
			req.Reply(true, nil)
			cmd = exec.Command(s.shellPath)
			cmd.Env = append(os.Environ(), "TERM="+term)

			var err error
			ptyFile, err = pty.StartWithSize(cmd, win)
			if err != nil {
				ch.SendRequest("exit-status", false, xssh.Marshal(&struct{ Status uint32 }{1}))
				return
			}
			go ptyWait(cmd, ptyFile, ch, done)
			go drainSessionReqs(cmd, ptyFile, reqs)
			goto waitExit

		case "exec":
			var p struct{ Cmd string }
			if err := xssh.Unmarshal(req.Payload, &p); err != nil {
				req.Reply(false, nil)
				ch.SendRequest("exit-status", false, xssh.Marshal(&struct{ Status uint32 }{1}))
				return
			}
			req.Reply(true, nil)
			cmd = exec.Command(s.shellPath, "-c", p.Cmd)

			if ptyReq {
				cmd.Env = append(os.Environ(), "TERM="+term)
				var err error
				ptyFile, err = pty.StartWithSize(cmd, win)
				if err != nil {
					ch.SendRequest("exit-status", false, xssh.Marshal(&struct{ Status uint32 }{1}))
					return
				}
				go ptyWait(cmd, ptyFile, ch, done)
				go drainSessionReqs(cmd, ptyFile, reqs)
				goto waitExit
			}

			cmd.Env = os.Environ()
			cmd.Stdin = ch
			cmd.Stdout = ch
			cmd.Stderr = ch.Stderr()
			go func() {
				_ = cmd.Run()
				ec := uint32(0)
				if cmd.ProcessState != nil {
					ec = uint32(cmd.ProcessState.ExitCode())
				}
				done <- ec
			}()
			go drainSessionReqs(cmd, nil, reqs)
			goto waitExit

		default:
			req.Reply(false, nil)
		}
	}
	return

waitExit:
	ec := <-done
	ch.SendRequest("exit-status", false, xssh.Marshal(&struct{ Status uint32 }{ec}))
}

func ptyWait(cmd *exec.Cmd, ptyFile *os.File, ch xssh.Channel, done chan<- uint32) {
	defer func() {
		_ = cmd.Wait()
		ec := uint32(0)
		if cmd.ProcessState != nil {
			ec = uint32(cmd.ProcessState.ExitCode())
		}
		done <- ec
	}()
	go func() { _, _ = io.Copy(ptyFile, ch) }()
	_, _ = io.Copy(ch, ptyFile)
	_ = ptyFile.Close()
}

func drainSessionReqs(cmd *exec.Cmd, ptyFile *os.File, reqs <-chan *xssh.Request) {
	for req := range reqs {
		switch req.Type {
		case "window-change":
			if ptyFile == nil {
				continue
			}
			var p struct {
				Cols   uint32
				Rows   uint32
				Width  uint32
				Height uint32
			}
			if err := xssh.Unmarshal(req.Payload, &p); err != nil {
				continue
			}
			if p.Cols > 0 && p.Rows > 0 {
				_ = pty.Setsize(ptyFile, &pty.Winsize{
					Rows: uint16(p.Rows),
					Cols: uint16(p.Cols),
				})
			}
		case "signal":
			if cmd == nil || cmd.Process == nil {
				continue
			}
			var p struct{ Signal string }
			if err := xssh.Unmarshal(req.Payload, &p); err != nil {
				continue
			}
			var sig syscall.Signal
			switch p.Signal {
			case "HUP":
				sig = syscall.SIGHUP
			case "INT":
				sig = syscall.SIGINT
			case "KILL":
				sig = syscall.SIGKILL
			case "TERM":
				sig = syscall.SIGTERM
			case "QUIT":
				sig = syscall.SIGQUIT
			default:
				continue
			}
			_ = cmd.Process.Signal(sig)
		}
	}
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
