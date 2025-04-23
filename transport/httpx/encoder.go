package httpx

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"io"
	"net/http"
)

type Encoder struct {
	hs func(http.Header)
	f  func(context.Context, any) (io.Reader, error)
}

func JsonEncoder() Encoder {
	return Encoder{
		hs: func(h http.Header) {
			h.Set("Content-Type", "application/json")
		},
		f: func(ctx context.Context, payload any) (io.Reader, error) {
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(payload)
			return &buf, err
		},
	}
}

func GobEncoder() Encoder {
	return Encoder{
		hs: func(h http.Header) {
			h.Set("Content-Type", "application/gob")
		},
		f: func(ctx context.Context, payload any) (io.Reader, error) {
			var buf bytes.Buffer
			err := gob.NewEncoder(&buf).Encode(payload)
			return &buf, err
		},
	}
}
