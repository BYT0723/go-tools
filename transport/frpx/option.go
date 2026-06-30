package frpx

import "time"

type Option func(*Client)

func WithServerAddr(addr string) Option {
	return func(c *Client) { c.addr = addr }
}

func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

func WithUser(user string) Option {
	return func(c *Client) { c.user = user }
}

func WithHeartbeatInterval(d time.Duration) Option {
	return func(c *Client) { c.heartbeatInterval = d }
}

func WithProxy(cfg ProxyConfig) Option {
	return func(c *Client) { c.proxyConfigs = append(c.proxyConfigs, cfg) }
}

func WithReconnectBackoff(d time.Duration) Option {
	return func(c *Client) { c.reconnectBackoff = d }
}
