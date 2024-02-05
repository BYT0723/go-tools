package http

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
)

type Decoder func(context.Context, *http.Response, any) error

func JsonDecoder(ctx context.Context, resp *http.Response, payload any) error {
	return json.NewDecoder(resp.Body).Decode(payload)
}

func GobDecoder(ctx context.Context, resp *http.Response, payload any) error {
	return gob.NewDecoder(resp.Body).Decode(payload)
}
