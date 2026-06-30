package frpx

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const (
	FrameTypeClientHello uint16 = 1
	FrameTypeServerHello uint16 = 2
	FrameTypeMessage     uint16 = 16
)

const (
	msgTypeLogin          uint16 = 1
	msgTypeLoginResp      uint16 = 2
	msgTypeNewProxy       uint16 = 3
	msgTypeNewProxyResp   uint16 = 4
	msgTypeNewWorkConn    uint16 = 6
	msgTypeReqWorkConn    uint16 = 7
	msgTypeStartWorkConn  uint16 = 8
	msgTypePing           uint16 = 11
	msgTypePong           uint16 = 12
)

type (
	Login struct {
		Version      string `json:"version,omitempty"`
		Hostname     string `json:"hostname,omitempty"`
		Os           string `json:"os,omitempty"`
		Arch         string `json:"arch,omitempty"`
		User         string `json:"user,omitempty"`
		PrivilegeKey string `json:"privilege_key,omitempty"`
		Timestamp    int64  `json:"timestamp,omitempty"`
		RunID        string `json:"run_id,omitempty"`
		PoolCount    int    `json:"pool_count,omitempty"`
	}

	LoginResp struct {
		RunID string `json:"run_id,omitempty"`
		Error string `json:"error,omitempty"`
	}

	NewProxy struct {
		ProxyName     string   `json:"proxy_name,omitempty"`
		ProxyType     string   `json:"proxy_type,omitempty"`
		RemotePort    int      `json:"remote_port,omitempty"`
		CustomDomains []string `json:"custom_domains,omitempty"`
		SubDomain     string   `json:"subdomain,omitempty"`
		HTTPUser      string   `json:"http_user,omitempty"`
		HTTPPwd       string   `json:"http_pwd,omitempty"`
		HostHeaderRw  string   `json:"host_header_rewrite,omitempty"`
	}

	NewProxyResp struct {
		ProxyName string `json:"proxy_name,omitempty"`
		Error     string `json:"error,omitempty"`
	}

	NewWorkConn struct {
		RunID        string `json:"run_id,omitempty"`
		PrivilegeKey string `json:"privilege_key,omitempty"`
		Timestamp    int64  `json:"timestamp,omitempty"`
	}

	StartWorkConn struct {
		ProxyName string `json:"proxy_name,omitempty"`
		Error     string `json:"error,omitempty"`
	}

	Ping struct {
		PrivilegeKey string `json:"privilege_key,omitempty"`
		Timestamp    int64  `json:"timestamp,omitempty"`
	}

	Pong struct {
		Error string `json:"error,omitempty"`
	}
)

var magicBytes = []byte("FRP\x00\x02\r\n")

func readV2Frame(r io.Reader) (uint16, []byte, error) {
	var hdr [8]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return 0, nil, err
	}
	typ := binary.BigEndian.Uint16(hdr[0:2])
	plen := binary.BigEndian.Uint32(hdr[4:8])
	if plen > 0 {
		body := make([]byte, plen)
		if _, err := io.ReadFull(r, body); err != nil {
			return 0, nil, err
		}
		return typ, body, nil
	}
	return typ, nil, nil
}

func writeV2Frame(w io.Writer, typ uint16, body []byte) error {
	var hdr [8]byte
	binary.BigEndian.PutUint16(hdr[0:2], typ)
	binary.BigEndian.PutUint32(hdr[4:8], uint32(len(body)))
	if _, err := w.Write(hdr[:]); err != nil {
		return err
	}
	if len(body) > 0 {
		_, err := w.Write(body)
		return err
	}
	return nil
}

func readV2Msg(r io.Reader) (uint16, []byte, error) {
	typ, body, err := readV2Frame(r)
	if err != nil {
		return 0, nil, err
	}
	if typ != FrameTypeMessage {
		return 0, nil, fmt.Errorf("frpx: unexpected frame type %d", typ)
	}
	if len(body) < 2 {
		return 0, nil, fmt.Errorf("frpx: message frame too short")
	}
	msgType := binary.BigEndian.Uint16(body[0:2])
	return msgType, body[2:], nil
}

func writeV2Msg(w io.Writer, msgType uint16, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	buf := make([]byte, 2+len(body))
	binary.BigEndian.PutUint16(buf[0:2], msgType)
	copy(buf[2:], body)
	return writeV2Frame(w, FrameTypeMessage, buf)
}
