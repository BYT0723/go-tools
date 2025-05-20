package encoder

import (
	"io"
	"net/http"
)

type Compressor interface {
	RequestHeader() http.Header
	Compress(writer io.Writer, src []byte) error
}
