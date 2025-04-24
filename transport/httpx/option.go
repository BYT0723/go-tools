package httpx

import (
	"net/http"
)

// Option configures the Client.
type Option func(*Client)

// WithEncoder sets a custom encoder for the Client.
//
// The encoder is responsible for converting the payload to the request body.
func WithEncoder(encoder Encoder) Option {
	return func(c *Client) {
		c.encoder = encoder
	}
}

// WithDecoder sets a custom decoder for the Client.
//
// The decoder is responsible for parsing the response body into the result object.
func WithDecoder(decoder Decoder) Option {
	return func(c *Client) {
		c.decoder = decoder
	}
}

// WithHttpClient sets the underlying http.Client used for making requests.
func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.cli = client
	}
}
