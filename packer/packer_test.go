package packer

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGzipUncompressor(t *testing.T) {
	Convey("GzipUncompressor 测试", t, func() {
		u := GzipUncompressor()
		So(u, ShouldNotBeNil)
	})
}

func TestBzip2Uncompressor(t *testing.T) {
	Convey("Bzip2Uncompressor 测试", t, func() {
		u := Bzip2Uncompressor()
		So(u, ShouldNotBeNil)
	})
}

func TestXzUncompressor(t *testing.T) {
	Convey("XzUncompressor 测试", t, func() {
		u := XzUncompressor()
		So(u, ShouldNotBeNil)
	})
}

func TestZstdUncompressor(t *testing.T) {
	Convey("ZstdUncompressor 测试", t, func() {
		u := ZstdUncompressor()
		So(u, ShouldNotBeNil)
	})
}

func TestUncompressorType(t *testing.T) {
	Convey("Uncompressor 类型测试", t, func() {
		var u Uncompressor = GzipUncompressor()
		So(u, ShouldNotBeNil)
	})
}

func TestUnzip(t *testing.T) {
	Convey("Unzip 测试", t, func() {
		Convey("解压不存在的文件返回错误", func() {
			err := Unzip("/nonexistent/zip/file.zip", "/tmp/test_unzip")
			So(err, ShouldNotBeNil)
		})

		Convey("解压有效的zip文件", func() {
			tmpDir := t.TempDir()
			zipPath := filepath.Join(tmpDir, "test.zip")

			zf, err := os.Create(zipPath)
			So(err, ShouldBeNil)
			zw := zip.NewWriter(zf)
			w, err := zw.Create("test.txt")
			So(err, ShouldBeNil)
			_, err = w.Write([]byte("hello world"))
			So(err, ShouldBeNil)
			zw.Close()
			zf.Close()

			destDir := filepath.Join(tmpDir, "output")
			err = Unzip(zipPath, destDir)
			So(err, ShouldBeNil)

			content, err := os.ReadFile(filepath.Join(destDir, "test.txt"))
			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, "hello world")
		})
	})
}
