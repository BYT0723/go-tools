package decoder

import (
	"encoding/json"
	"io"
	"net/http"
)

type jsonDecoder struct{}

func (d jsonDecoder) Decode(reader io.Reader, header http.Header, payload any) (err error) {
	ct := header.Get("Content-Type")

	switch ct {
	case "application/json":
		err = json.NewDecoder(reader).Decode(payload)
	default:
		err = ErrInvalidContentType
	}

	if err != nil {
		return err
	}

	return json.NewDecoder(reader).Decode(payload)
}

func JSONDecoder() jsonDecoder { return jsonDecoder{} }
