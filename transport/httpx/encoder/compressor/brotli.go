package compressor

import (
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
)

type brotliCompressor struct {
	level int
}

func (c *brotliCompressor) RequestHeader() http.Header {
	return http.Header{
		"Content-Encoding": []string{"br"},
	}
}

func (c *brotliCompressor) Compress(writer io.Writer, src []byte) error {
	_, err := brotli.NewWriterOptions(writer, brotli.WriterOptions{Quality: c.level}).Write(src)
	return err
}

func BrotliCompressor() *brotliCompressor {
	return BrotliCompressorWithLevel(brotli.DefaultCompression)
}

func BrotliCompressorWithLevel(level int) *brotliCompressor {
	return &brotliCompressor{level: level}
}
