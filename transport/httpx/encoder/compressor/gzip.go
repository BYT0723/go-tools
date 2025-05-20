package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

type gzipCompressor struct {
	level int
}

func (c *gzipCompressor) RequestHeader() http.Header {
	return http.Header{
		"Content-Encoding": []string{"gzip"},
	}
}

func (c *gzipCompressor) Compress(writer io.Writer, src []byte) error {
	w, err := gzip.NewWriterLevel(writer, c.level)
	if err != nil {
		return err
	}
	_, err = w.Write(src)
	return err
}

func GZipCompressor() *gzipCompressor {
	return GZipCompressorWithLevel(gzip.DefaultCompression)
}

func GZipCompressorWithLevel(level int) *gzipCompressor {
	return &gzipCompressor{level: level}
}
