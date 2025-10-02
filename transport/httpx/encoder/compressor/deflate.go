package compressor

import (
	"compress/flate"
	"io"
	"net/http"
)

type deflateCompressor struct {
	level int
	dict  []byte
}

func (c *deflateCompressor) RequestHeader() http.Header {
	return http.Header{
		"Content-Encoding": []string{"deflate"},
	}
}

func (c *deflateCompressor) Compress(writer io.Writer, src []byte) error {
	var (
		w   *flate.Writer
		err error
	)
	if len(c.dict) > 0 {
		w, err = flate.NewWriterDict(writer, c.level, c.dict)
	} else {
		w, err = flate.NewWriter(writer, c.level)
	}
	if err != nil {
		return err
	}
	_, err = w.Write(src)
	return err
}

func DeflateCompressor(level int) deflateCompressor {
	return DeflateCompressorWithDict(level, nil)
}

func DeflateCompressorWithDict(level int, dict []byte) deflateCompressor {
	return deflateCompressor{level: level, dict: dict}
}
