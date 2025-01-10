package zerologger

import (
	"io"

	"github.com/rs/zerolog"
)

type LevelFilter func(zerolog.Level) bool

type LevelWriter struct {
	io.Writer
	filter LevelFilter
}

func (w *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if w.filter(level) {
		return w.Writer.Write(p)
	}
	return len(p), nil
}

func NewLevelWriter(writer io.Writer, filter LevelFilter) zerolog.LevelWriter {
	return &LevelWriter{Writer: writer, filter: filter}
}
