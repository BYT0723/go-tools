package http

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
)

type Encoder func(context.Context, any) ([]byte, error)

func JsonEncoder(ctx context.Context, payload any) ([]byte, error) {
	return json.Marshal(payload)
}

func GobEncoder(ctx context.Context, payload any) ([]byte, error) {
	var buf bytes.Buffer
	defer buf.Reset()
	err := gob.NewEncoder(&buf).Encode(payload)
	return buf.Bytes(), err
}
