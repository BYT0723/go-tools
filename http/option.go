package http

import (
	"net/http"
)

type Option func(*Client)

func WithEncoder(encoder Encoder) Option {
	return func(c *Client) {
		c.encoder = encoder
	}
}

func WithDecoder(decoder Decoder) Option {
	return func(c *Client) {
		c.decoder = decoder
	}
}

func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.innerClient = client
	}
}
