package encoder

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
	Convey("JSONEncoder 测试", t, func() {
		Convey("Encode 正常数据", func() {
			e := JSONEncoder()
			payload := map[string]string{"key": "value"}
			reader, err := e.Encode(payload)
			So(err, ShouldBeNil)
			So(reader, ShouldNotBeNil)
		})

		Convey("RequestHeader 返回正确的Content-Type", func() {
			e := JSONEncoder()
			h := e.RequestHeader()
			So(h.Get("Content-Type"), ShouldEqual, "application/json")
		})

		Convey("JSONEncoderWithCompressor 压缩数据", func() {
			var buf bytes.Buffer
			mc := &mockCompressor{
				header: http.Header{"Content-Encoding": []string{"gzip"}},
			}
			e := JSONEncoderWithCompressor(mc)
			payload := map[string]string{"key": "value"}
			reader, err := e.Encode(payload)
			So(err, ShouldBeNil)
			_ = buf
			_, err = io.ReadAll(reader)
			So(err, ShouldBeNil)
		})
	})
}

func TestGobEncoder(t *testing.T) {
	Convey("GobEncoder 测试", t, func() {
		Convey("Encode struct数据", func() {
			e := GobEncoder()
			type testStruct struct {
				Name string
				Age  int
			}
			payload := testStruct{Name: "test", Age: 30}
			reader, err := e.Encode(payload)
			So(err, ShouldBeNil)
			So(reader, ShouldNotBeNil)
		})

		Convey("RequestHeader 返回正确的Content-Type", func() {
			e := GobEncoder()
			h := e.RequestHeader()
			So(h.Get("Content-Type"), ShouldEqual, "application/gob")
		})
	})
}

func TestCompressorInterface(t *testing.T) {
	Convey("Compressor 接口定义验证", t, func() {
		var c Compressor = &mockCompressor{}
		So(c, ShouldNotBeNil)
	})
}
