package frpx

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var version = "0.1.0"

type Client struct {
	mu     sync.Mutex
	addr   string
	token  string
	user   string

	heartbeatInterval time.Duration
	reconnectBackoff  time.Duration
	proxyConfigs      []ProxyConfig

	crw     io.ReadWriter
	rawConn io.Closer
	runID   string
	proxies map[string]ProxyConfig

	wg sync.WaitGroup
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		addr:              "127.0.0.1:7000",
		heartbeatInterval: 30 * time.Second,
		reconnectBackoff:  2 * time.Second,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) Name() string                { return "frpx" }
func (c *Client) Init(_ context.Context) error { return nil }

func (c *Client) Run(ctx context.Context) error {
	for ctx.Err() == nil {
		if err := c.runOnce(ctx); err != nil {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(c.reconnectBackoff):
			}
		}
	}
	return nil
}

func (c *Client) Destroy(_ context.Context) error {
	c.mu.Lock()
	if c.rawConn != nil {
		c.rawConn.Close()
	}
	c.mu.Unlock()
	c.wg.Wait()
	return nil
}

func (c *Client) runOnce(ctx context.Context) error {
	conn, err := c.dial(ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.rawConn = conn
	c.mu.Unlock()

	crw, err := v2Handshake(conn, c.token)
	if err != nil {
		conn.Close()
		return err
	}

	c.mu.Lock()
	c.crw = crw
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.crw = nil
		c.rawConn = nil
		c.mu.Unlock()
		conn.Close()
	}()

	if err := c.login(crw); err != nil {
		return err
	}

	c.registerProxies(crw)

	return c.eventLoop(ctx, crw)
}

func (c *Client) dial(ctx context.Context) (net.Conn, error) {
	d := net.Dialer{Timeout: 10 * time.Second}
	return d.DialContext(ctx, "tcp", c.addr)
}

func authKey(token string, ts int64) string {
	h := md5.Sum([]byte(token + strconv.FormatInt(ts, 10)))
	return hex.EncodeToString(h[:])
}

func (c *Client) login(crw io.ReadWriter) error {
	ts := time.Now().Unix()
	hostname, _ := os.Hostname()

	msg := Login{
		Version:      version,
		Hostname:     hostname,
		Os:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		User:         c.user,
		PrivilegeKey: authKey(c.token, ts),
		Timestamp:    ts,
		RunID:        c.runID,
		PoolCount:    1,
	}

	if err := writeV2Msg(crw, msgTypeLogin, msg); err != nil {
		return fmt.Errorf("frpx: login write: %w", err)
	}

	_, body, err := readV2Msg(crw)
	if err != nil {
		return fmt.Errorf("frpx: login read: %w", err)
	}
	var resp LoginResp
	json.Unmarshal(body, &resp)
	if resp.Error != "" {
		return fmt.Errorf("frpx: login rejected: %s", resp.Error)
	}

	c.mu.Lock()
	c.runID = resp.RunID
	c.proxies = make(map[string]ProxyConfig)
	for _, cfg := range c.proxyConfigs {
		c.proxies[cfg.Name] = cfg
	}
	c.mu.Unlock()

	return nil
}

func (c *Client) registerProxies(crw io.ReadWriter) {
	c.mu.Lock()
	cfgs := make([]ProxyConfig, len(c.proxyConfigs))
	copy(cfgs, c.proxyConfigs)
	user := c.user
	c.mu.Unlock()

	for _, cfg := range cfgs {
		msg := cfg.toNewProxy(user)
		if err := writeV2Msg(crw, msgTypeNewProxy, msg); err != nil {
			return
		}
		if _, _, err := readV2Msg(crw); err != nil {
			return
		}
	}
}

func (c *Client) eventLoop(ctx context.Context, crw io.ReadWriter) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan error, 1)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		done <- c.readLoop(ctx, crw)
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.heartbeatLoop(ctx, crw)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) readLoop(ctx context.Context, crw io.ReadWriter) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		typ, _, err := readV2Msg(crw)
		if err != nil {
			return err
		}

		switch typ {
		case msgTypePong:
		case msgTypeReqWorkConn:
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.handleReqWorkConn(ctx)
			}()
		}
	}
}

func (c *Client) heartbeatLoop(ctx context.Context, crw io.ReadWriter) {
	c.mu.Lock()
	interval := c.heartbeatInterval
	c.mu.Unlock()

	if interval <= 0 {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ts := time.Now().Unix()
			writeV2Msg(crw, msgTypePing, Ping{
				PrivilegeKey: authKey(c.token, ts),
				Timestamp:    ts,
			})
		}
	}
}
