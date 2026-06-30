package connmux

import (
	"bytes"
)

type Matcher interface {
	Match(sniffed []byte) bool
}

type bytePrefixMatcher struct {
	prefix []byte
}

func (m *bytePrefixMatcher) Match(sniffed []byte) bool {
	return bytes.HasPrefix(sniffed, m.prefix)
}

type defaultMatcher struct{}

func (defaultMatcher) Match([]byte) bool { return true }

var (
	MatchSSH    Matcher = &bytePrefixMatcher{prefix: []byte("SSH-")}
	MatchHTTP1  Matcher = http1Matcher{}
	MatchDefault Matcher = defaultMatcher{}
)

type http1Matcher struct{}

var httpMethods = [][]byte{
	[]byte("GET "),
	[]byte("POST"),
	[]byte("HEAD"),
	[]byte("PUT "),
	[]byte("DELE"),
	[]byte("OPTI"),
	[]byte("PATC"),
	[]byte("CONN"),
	[]byte("TRAC"),
}

func (http1Matcher) Match(sniffed []byte) bool {
	for _, m := range httpMethods {
		if bytes.HasPrefix(sniffed, m) {
			return true
		}
	}
	return false
}
