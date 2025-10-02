package decoder

import "errors"

var (
	ErrInvalidContentType   = errors.New("invalid content type")
	ErrNotMatchCompressType = errors.New("not match compress type")
)
