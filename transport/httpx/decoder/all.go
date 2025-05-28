package decoder

import (
	"encoding/json"
	"io"
	"net/http"
)

type defaultDecoder struct{}

func (defaultDecoder) Decode(reader io.Reader, header http.Header, payload any) (err error) {
	ct := header.Get("Content-Type")

	switch ct {
	case "application/json":
		err = json.NewDecoder(reader).Decode(payload)
	}
	return
}

func DefaultDecoder() defaultDecoder { return defaultDecoder{} }
