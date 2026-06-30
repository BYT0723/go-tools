package httpx

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/transport/connmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapHTTPServer(t *testing.T) {
	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	})}
	svc := WrapHTTPServer(srv)
	assert.Equal(t, "http", svc.Name())
	assert.Equal(t, connmux.MatchHTTP1, svc.Match())
}

func TestHTTPServiceIntegration(t *testing.T) {
	mux := connmux.NewMux()

	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("httpx-http-ok"))
	})}

	mux.Route("http", WrapHTTPServer(srv))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mux.Run(ctx)
	time.Sleep(100 * time.Millisecond)

	addr := mux.Addr()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n"))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	require.NoError(t, err)
	resp := string(buf[:n])
	assert.Contains(t, resp, "httpx-http-ok")
	conn.Close()

	cancel()
}
