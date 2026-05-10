package decoder

import (
	"bytes"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONDecoder(t *testing.T) {
	Convey("JSONDecoder 测试", t, func() {
		Convey("无效 Content-Type 返回错误", func() {
			d := JSONDecoder()
			payload := make(map[string]string)
			reader := bytes.NewBufferString(`{"key":"value"}`)
			header := http.Header{"Content-Type": []string{"text/plain"}}
			err := d.Decode(reader, header, &payload)
			So(err, ShouldEqual, ErrInvalidContentType)
		})
	})
}

func TestErrorVars(t *testing.T) {
	Convey("错误变量测试", t, func() {
		So(ErrInvalidContentType.Error(), ShouldEqual, "invalid content type")
		So(ErrNotMatchCompressType.Error(), ShouldEqual, "not match compress type")
	})
}
