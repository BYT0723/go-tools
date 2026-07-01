package connmux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchSSH(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		match bool
	}{
		{"SSH banner", []byte("SSH-2.0-OpenSSH_9.0"), true},
		{"SSH short", []byte("SSH-"), true},
		{"not SSH", []byte("GET / HTTP/1.1"), false},
		{"empty", []byte{}, false},
		{"partial S", []byte("S"), false},
		{"partial SS", []byte("SS"), false},
		{"partial SSH", []byte("SSH"), false},
		{"SSH lowercase", []byte("ssh-2.0"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.match, MatchSSH.Match(tt.input))
		})
	}
}

func TestMatchHTTP1(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		match bool
	}{
		{"GET", []byte("GET / HTTP/1.1"), true},
		{"POST", []byte("POST /api HTTP/1.1"), true},
		{"HEAD", []byte("HEAD / HTTP/1.1"), true},
		{"PUT", []byte("PUT / HTTP/1.1"), true},
		{"DELETE", []byte("DELETE / HTTP/1.1"), true},
		{"OPTIONS", []byte("OPTIONS / HTTP/1.1"), true},
		{"PATCH", []byte("PATCH / HTTP/1.1"), true},
		{"CONNECT", []byte("CONNECT proxy:443 HTTP/1.1"), true},
		{"TRACE", []byte("TRACE / HTTP/1.1"), true},
		{"SSH banner", []byte("SSH-2.0-OpenSSH"), false},
		{"binary data", []byte{0x00, 0x01, 0x02, 0x03}, false},
		{"empty", []byte{}, false},
		{"GETTER (no space)", []byte("GETTER"), false},
		{"get lowercase", []byte("get /"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.match, MatchHTTP1.Match(tt.input))
		})
	}
}

func TestMatchDefault(t *testing.T) {
	assert.True(t, MatchDefault.Match(nil))
	assert.True(t, MatchDefault.Match([]byte{}))
	assert.True(t, MatchDefault.Match([]byte{0x00}))
	assert.True(t, MatchDefault.Match([]byte("anything")))
}

func TestMatchHTTP2(t *testing.T) {
	assert.True(t, MatchHTTP2.Match([]byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")))
	assert.True(t, MatchHTTP2.Match([]byte("PRI * HTTP/2.0")))
	assert.False(t, MatchHTTP2.Match([]byte("PRI")))
	assert.False(t, MatchHTTP2.Match([]byte("GET / HTTP/2.0")))
	assert.False(t, MatchHTTP2.Match([]byte{}))
}

func TestMatchTLS(t *testing.T) {
	assert.True(t, MatchTLS.Match([]byte{0x16, 0x03, 0x01})) // TLS 1.0
	assert.True(t, MatchTLS.Match([]byte{0x16, 0x03, 0x03})) // TLS 1.2/1.3
	assert.True(t, MatchTLS.Match([]byte{0x16, 0x03}))       // minimal
	assert.False(t, MatchTLS.Match([]byte{0x16}))
	assert.False(t, MatchTLS.Match([]byte{0x17, 0x03})) // application data record
	assert.False(t, MatchTLS.Match([]byte{}))
}
