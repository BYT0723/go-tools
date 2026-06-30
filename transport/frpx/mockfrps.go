package frpx

import (
	"encoding/json"
	"io"
	"net"
)

type mockFrps struct {
	listener net.Listener
	runID    string
}

func newMockFrps() *mockFrps {
	return &mockFrps{runID: "mock-run-id"}
}

func (m *mockFrps) listen() (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	m.listener = l
	go m.serve()
	return l.Addr().String(), nil
}

func (m *mockFrps) close() {
	if m.listener != nil {
		m.listener.Close()
	}
}

func (m *mockFrps) serve() {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			return
		}
		go m.handleConn(conn)
	}
}

func (m *mockFrps) handleConn(conn net.Conn) {
	defer conn.Close()

	rl := io.LimitReader(conn, 7)
	magic := make([]byte, 7)
	io.ReadFull(rl, magic)

	readV2Frame(conn)

	sh := serverHello{
		Capabilities: serverCaps{
			Message: messageSel{Codec: "json"},
			Crypto:  cryptoSel{},
		},
	}
	shJSON, _ := json.Marshal(sh)
	writeV2Frame(conn, FrameTypeServerHello, shJSON)

	typ, body, err := readV2Msg(conn)
	if err != nil || typ != msgTypeLogin {
		return
	}
	var login Login
	json.Unmarshal(body, &login)
	_ = login

	writeV2Msg(conn, msgTypeLoginResp, LoginResp{
		RunID: m.runID,
	})

	m.processMessages(conn)
}

func (m *mockFrps) processMessages(rw io.ReadWriter) {
	for {
		typ, body, err := readV2Msg(rw)
		if err != nil {
			return
		}

		switch typ {
		case msgTypeNewProxy:
			var proxy NewProxy
			json.Unmarshal(body, &proxy)
			writeV2Msg(rw, msgTypeNewProxyResp, NewProxyResp{
				ProxyName: proxy.ProxyName,
			})

		case msgTypePing:
			writeV2Msg(rw, msgTypePong, Pong{})
		}
	}
}
