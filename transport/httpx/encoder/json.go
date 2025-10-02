package encoder

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"net/http"
)

type jsonEncoder struct {
	compressor Compressor
	ct         string
}

func (d jsonEncoder) RequestHeader() http.Header {
	header := http.Header{"Content-Type": []string{d.ct}}
	if d.compressor != nil {
		maps.Copy(header, d.RequestHeader())
	}
	return header
}

func (d jsonEncoder) Encode(payload any) (io.Reader, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if d.compressor == nil {
		return &buf, err
	}

	var cbuf bytes.Buffer
	err = d.compressor.Compress(&cbuf, buf.Bytes())
	return &cbuf, err
}

func JSONEncoder() jsonEncoder {
	return jsonEncoder{
		ct: "application/json",
	}
}

func JSONEncoderWithCompressor(compressor Compressor) jsonEncoder {
	return jsonEncoder{
		compressor: compressor,
		ct:         "application/json",
	}
}
