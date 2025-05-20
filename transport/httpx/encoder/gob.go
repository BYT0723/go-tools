package encoder

import (
	"bytes"
	"encoding/gob"
	"io"
	"net/http"
)

type gobEncoder struct {
	ct string
}

func (d *gobEncoder) RequestHeader() http.Header {
	return http.Header{
		"Content-Type": []string{d.ct},
	}
}

func (d *gobEncoder) Encode(payload any) (io.Reader, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(payload)
	return &buf, err
}

func GobEncoder() *gobEncoder {
	return &gobEncoder{
		ct: "application/gob",
	}
}
