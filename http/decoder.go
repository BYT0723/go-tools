package http

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"io"
)

type Decoder func(context.Context, io.Reader, any) error

func JsonDecoder(ctx context.Context, reader io.Reader, payload any) error {
	return json.NewDecoder(reader).Decode(payload)
}

func GobDecoder(ctx context.Context, reader io.Reader, payload any) error {
	return gob.NewDecoder(reader).Decode(payload)
}
