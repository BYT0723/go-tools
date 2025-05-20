package decoder

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
)

type defaultDecoder struct{}

func (defaultDecoder) Decode(reader io.Reader, header http.Header, payload any) (err error) {
	var (
		ct = header.Get("Content-Type")
		ce = header.Get("Content-Encoding")
	)
	switch ce {
	case "gzip":
		r, err := gzip.NewReader(reader)
		if err != nil {
			return fmt.Errorf("gzip decode failed: %w", err)
		}
		defer r.Close()
		reader = r

	case "deflate":
		r := flate.NewReader(reader)
		defer r.Close()
		reader = r

	case "br":
		r := brotli.NewReader(reader)
		reader = r
	case "":
	default:
		return fmt.Errorf("unsupported content-encoding: %s", ct)
	}

	switch ct {
	case "application/json":
		err = json.NewDecoder(reader).Decode(payload)
	}
	return
}

func DefaultDecoder() defaultDecoder { return defaultDecoder{} }
