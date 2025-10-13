package packer

import (
	"compress/bzip2"
	"compress/gzip"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

type Uncompressor func(io.Reader) (io.Reader, error)

func GzipUncompressor() Uncompressor {
	return func(r io.Reader) (io.Reader, error) {
		return gzip.NewReader(r)
	}
}

func Bzip2Uncompressor() Uncompressor {
	return func(r io.Reader) (io.Reader, error) {
		return bzip2.NewReader(r), nil
	}
}

func XzUncompressor() Uncompressor {
	return func(r io.Reader) (io.Reader, error) {
		return xz.NewReader(r)
	}
}

func ZstdUncompressor() Uncompressor {
	return func(r io.Reader) (io.Reader, error) {
		return zstd.NewReader(r)
	}
}
