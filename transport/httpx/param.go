package httpx

import (
	"net/http"
)

type Param func(*Request)

func WithHeader(header http.Header) Param {
	return func(r *Request) {
		r.Header = header
	}
}

func WithPayload(payload any) Param {
	return func(r *Request) {
		r.Payload = payload
	}
}
