package sshx

import "golang.org/x/crypto/ssh"

type Option func(*Server)

func WithAddr(addr string) Option {
	return func(s *Server) { s.addr = addr }
}

func WithHostKey(pemBytes []byte) Option {
	return func(s *Server) {
		key, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return
		}
		s.config.AddHostKey(key)
	}
}

func WithHandler(h Handler) Option {
	return func(s *Server) { s.handler = h }
}

func WithPasswordAuth(fn func(user, password string) bool) Option {
	return func(s *Server) {
		s.config.PasswordCallback = func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if fn(conn.User(), string(password)) {
				return nil, nil
			}
			return nil, ssh.ErrNoAuth
		}
	}
}

func WithPublicKeyAuth(fn func(conn ssh.ConnMetadata, key ssh.PublicKey) bool) Option {
	return func(s *Server) {
		s.config.PublicKeyCallback = func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if fn(conn, key) {
				return nil, nil
			}
			return nil, ssh.ErrNoAuth
		}
	}
}

func WithShellPath(path string) Option {
	return func(s *Server) { s.shellPath = path }
}

func WithUser(name, password string) Option {
	return WithPasswordAuth(func(user, pwd string) bool {
		return user == name && pwd == password
	})
}
