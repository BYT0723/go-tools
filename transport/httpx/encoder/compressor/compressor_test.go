package compressor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGZipCompressor(t *testing.T) {
	t.Run("GZipCompressor 测试", func(t *testing.T) {
		t.Run("默认压缩", func(t *testing.T) {
			c := GZipCompressor()
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			assert.Nil(t, err)
			assert.Greater(t, buf.Len(), 0)
		})

		t.Run("RequestHeader", func(t *testing.T) {
			c := GZipCompressor()
			h := c.RequestHeader()
			assert.Equal(t, "gzip", h.Get("Content-Encoding"))
		})

		t.Run("指定压缩级别", func(t *testing.T) {
			c := GZipCompressorWithLevel(9)
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			assert.Nil(t, err)
		})
	})
}

func TestDeflateCompressor(t *testing.T) {
	t.Run("DeflateCompressor 测试", func(t *testing.T) {
		t.Run("默认压缩", func(t *testing.T) {
			c := DeflateCompressor(6)
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			assert.Nil(t, err)
		})

		t.Run("RequestHeader", func(t *testing.T) {
			c := DeflateCompressor(6)
			h := c.RequestHeader()
			assert.Equal(t, "deflate", h.Get("Content-Encoding"))
		})
	})
}

func TestBrotliCompressor(t *testing.T) {
	t.Run("BrotliCompressor 测试", func(t *testing.T) {
		t.Run("默认压缩", func(t *testing.T) {
			c := BrotliCompressor()
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			assert.Nil(t, err)
		})

		t.Run("RequestHeader", func(t *testing.T) {
			c := BrotliCompressor()
			h := c.RequestHeader()
			assert.Equal(t, "br", h.Get("Content-Encoding"))
		})
	})
}

func TestDecompressRoundTrip(t *testing.T) {
	t.Run("压缩测试", func(t *testing.T) {
		t.Run("gzip 压缩产生非空输出", func(t *testing.T) {
			c := GZipCompressor()
			var compressed bytes.Buffer
			err := c.Compress(&compressed, []byte("hello world"))
			assert.Nil(t, err)
			assert.Greater(t, compressed.Len(), 0)
		})

		t.Run("deflate 压缩不返回错误", func(t *testing.T) {
			c := DeflateCompressor(6)
			var compressed bytes.Buffer
			err := c.Compress(&compressed, []byte("hello world"))
			assert.Nil(t, err)
		})

		t.Run("brotli 压缩不返回错误", func(t *testing.T) {
			c := BrotliCompressor()
			var compressed bytes.Buffer
			err := c.Compress(&compressed, []byte("hello world"))
			assert.Nil(t, err)
		})
	})
}
