package decoder

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONDecoder(t *testing.T) {
	t.Run("JSONDecoder 测试", func(t *testing.T) {
		t.Run("无效 Content-Type 返回错误", func(t *testing.T) {
			d := JSONDecoder()
			payload := make(map[string]string)
			reader := bytes.NewBufferString(`{"key":"value"}`)
			header := http.Header{"Content-Type": []string{"text/plain"}}
			err := d.Decode(reader, header, &payload)
			assert.Equal(t, ErrInvalidContentType, err)
		})
	})
}

func TestErrorVars(t *testing.T) {
	t.Run("错误变量测试", func(t *testing.T) {
		assert.Equal(t, "invalid content type", ErrInvalidContentType.Error())
		assert.Equal(t, "not match compress type", ErrNotMatchCompressType.Error())
	})
}
