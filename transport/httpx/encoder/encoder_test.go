package encoder

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCompressor struct {
	header http.Header
}

func (m *mockCompressor) RequestHeader() http.Header { return m.header }
func (m *mockCompressor) Compress(writer io.Writer, src []byte) error {
	_, err := writer.Write(src)
	return err
}

func TestJSONEncoder(t *testing.T) {
	t.Run("JSONEncoder 测试", func(t *testing.T) {
		t.Run("Encode 正常数据", func(t *testing.T) {
			e := JSONEncoder()
			payload := map[string]string{"key": "value"}
			reader, err := e.Encode(payload)
			assert.Nil(t, err)
			assert.NotNil(t, reader)
		})

		t.Run("RequestHeader 返回正确的Content-Type", func(t *testing.T) {
			e := JSONEncoder()
			h := e.RequestHeader()
			assert.Equal(t, "application/json", h.Get("Content-Type"))
		})

		t.Run("JSONEncoderWithCompressor 压缩数据", func(t *testing.T) {
			var buf bytes.Buffer
			mc := &mockCompressor{
				header: http.Header{"Content-Encoding": []string{"gzip"}},
			}
			e := JSONEncoderWithCompressor(mc)
			payload := map[string]string{"key": "value"}
			reader, err := e.Encode(payload)
			assert.Nil(t, err)
			_ = buf
			_, err = io.ReadAll(reader)
			assert.Nil(t, err)
		})
	})
}

func TestGobEncoder(t *testing.T) {
	t.Run("GobEncoder 测试", func(t *testing.T) {
		t.Run("Encode struct数据", func(t *testing.T) {
			e := GobEncoder()
			type testStruct struct {
				Name string
				Age  int
			}
			payload := testStruct{Name: "test", Age: 30}
			reader, err := e.Encode(payload)
			assert.Nil(t, err)
			assert.NotNil(t, reader)
		})

		t.Run("RequestHeader 返回正确的Content-Type", func(t *testing.T) {
			e := GobEncoder()
			h := e.RequestHeader()
			assert.Equal(t, "application/gob", h.Get("Content-Type"))
		})
	})
}

func TestCompressorInterface(t *testing.T) {
	t.Run("Compressor 接口定义验证", func(t *testing.T) {
		var c Compressor = &mockCompressor{}
		assert.NotNil(t, c)
	})
}
