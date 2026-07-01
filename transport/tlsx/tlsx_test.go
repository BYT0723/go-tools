package tlsx

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/transport/connmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapTLS(t *testing.T) {
	svc := WrapTLS(&tls.Config{}, func(net.Conn) {})
	assert.Equal(t, "tls", svc.Name())
	assert.Equal(t, connmux.MatchTLS, svc.Match())
}

func TestTLSIntegration(t *testing.T) {
	cert := newSelfSignedCert(t)
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	mux := connmux.NewMux()

	var received []byte
	var done = make(chan struct{})
	mux.Route("tls", WrapTLS(tlsCfg, func(conn net.Conn) {
		defer conn.Close()
		defer close(done)
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		received = buf[:n]
		conn.Write([]byte("TLS_OK"))
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mux.Run(ctx)
	time.Sleep(100 * time.Millisecond)

	addr := mux.Addr()

	tlsConn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	require.NoError(t, err)
	defer tlsConn.Close()

	tlsConn.Write([]byte("HELLO_TLS"))
	buf := make([]byte, 1024)
	n, err := tlsConn.Read(buf)
	require.NoError(t, err)

	<-done
	assert.Equal(t, "TLS_OK", string(buf[:n]))
	assert.True(t, bytes.HasPrefix(received, []byte("HELLO_TLS")))

	cancel()
}

func newSelfSignedCert(t *testing.T) tls.Certificate {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}
