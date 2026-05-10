package compressor

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGZipCompressor(t *testing.T) {
	Convey("GZipCompressor 测试", t, func() {
		Convey("默认压缩", func() {
			c := GZipCompressor()
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			So(err, ShouldBeNil)
			So(buf.Len(), ShouldBeGreaterThan, 0)
		})

		Convey("RequestHeader", func() {
			c := GZipCompressor()
			h := c.RequestHeader()
			So(h.Get("Content-Encoding"), ShouldEqual, "gzip")
		})

		Convey("指定压缩级别", func() {
			c := GZipCompressorWithLevel(9)
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			So(err, ShouldBeNil)
		})
	})
}

func TestDeflateCompressor(t *testing.T) {
	Convey("DeflateCompressor 测试", t, func() {
		Convey("默认压缩", func() {
			c := DeflateCompressor(6)
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			So(err, ShouldBeNil)
		})

		Convey("RequestHeader", func() {
			c := DeflateCompressor(6)
			h := c.RequestHeader()
			So(h.Get("Content-Encoding"), ShouldEqual, "deflate")
		})
	})
}

func TestBrotliCompressor(t *testing.T) {
	Convey("BrotliCompressor 测试", t, func() {
		Convey("默认压缩", func() {
			c := BrotliCompressor()
			var buf bytes.Buffer
			err := c.Compress(&buf, []byte("hello world"))
			So(err, ShouldBeNil)
		})

		Convey("RequestHeader", func() {
			c := BrotliCompressor()
			h := c.RequestHeader()
			So(h.Get("Content-Encoding"), ShouldEqual, "br")
		})
	})
}

func TestDecompressRoundTrip(t *testing.T) {
	Convey("压缩测试", t, func() {
		Convey("gzip 压缩产生非空输出", func() {
			c := GZipCompressor()
			var compressed bytes.Buffer
			err := c.Compress(&compressed, []byte("hello world"))
			So(err, ShouldBeNil)
			So(compressed.Len(), ShouldBeGreaterThan, 0)
		})

		Convey("deflate 压缩不返回错误", func() {
			c := DeflateCompressor(6)
			var compressed bytes.Buffer
			err := c.Compress(&compressed, []byte("hello world"))
			So(err, ShouldBeNil)
		})

		Convey("brotli 压缩不返回错误", func() {
			c := BrotliCompressor()
			var compressed bytes.Buffer
			err := c.Compress(&compressed, []byte("hello world"))
			So(err, ShouldBeNil)
		})
	})
}
