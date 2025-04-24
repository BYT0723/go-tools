package httpx

import (
	"net/http"
)

// Param is the function type used to configure Request.
// Users can set parameters for requests through a series of Param functions, such as Header or Payload.
type Param func(*request)

// WithHeader Returns a Param to set the requested header.
// Can be used to set custom request headers, such as Content-Type, Authorization, etc.
func WithHeader(header http.Header) Param {
	return func(r *request) {
		r.header = header
	}
}

// WithPayload Returns a Param to set the requested Payload (request body).
// payload map[string]any | struct | bytes | io.reader | string
func WithPayload(payload any) Param {
	return func(r *request) {
		r.payload = payload
	}
}
