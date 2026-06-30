package frpx

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"time"
)

func (c *Client) handleReqWorkConn(ctx context.Context) {
	c.mu.Lock()
	addr := c.addr
	token := c.token
	runID := c.runID
	c.mu.Unlock()

	d := net.Dialer{Timeout: 10 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return
	}

	crw, err := v2Handshake(conn, token)
	if err != nil {
		conn.Close()
		return
	}

	ts := time.Now().Unix()
	if err := writeV2Msg(crw, msgTypeNewWorkConn, NewWorkConn{
		RunID:        runID,
		PrivilegeKey: authKey(token, ts),
		Timestamp:    ts,
	}); err != nil {
		conn.Close()
		return
	}

	_, body, err := readV2Msg(crw)
	if err != nil {
		conn.Close()
		return
	}
	var start StartWorkConn
	json.Unmarshal(body, &start)
	if start.Error != "" {
		conn.Close()
		return
	}

	cfg := c.findProxy(start.ProxyName)
	if cfg == nil {
		conn.Close()
		return
	}

	local, err := d.DialContext(ctx, "tcp", cfg.LocalAddr)
	if err != nil {
		conn.Close()
		return
	}
	defer local.Close()

	go func() { _, _ = io.Copy(crw, local) }()
	_, _ = io.Copy(local, crw)
	conn.Close()
}

func (c *Client) findProxy(proxyName string) *ProxyConfig {
	c.mu.Lock()
	defer c.mu.Unlock()

	if cfg, ok := c.proxies[proxyName]; ok {
		return &cfg
	}
	for _, cfg := range c.proxies {
		if proxyName == c.user+"."+cfg.Name || proxyName == "."+cfg.Name {
			return &cfg
		}
	}
	return nil
}
