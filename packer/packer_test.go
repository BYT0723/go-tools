package packer

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipUncompressor(t *testing.T) {
	t.Run("GzipUncompressor 测试", func(t *testing.T) {
		u := GzipUncompressor()
		assert.NotNil(t, u)
	})
}

func TestBzip2Uncompressor(t *testing.T) {
	t.Run("Bzip2Uncompressor 测试", func(t *testing.T) {
		u := Bzip2Uncompressor()
		assert.NotNil(t, u)
	})
}

func TestXzUncompressor(t *testing.T) {
	t.Run("XzUncompressor 测试", func(t *testing.T) {
		u := XzUncompressor()
		assert.NotNil(t, u)
	})
}

func TestZstdUncompressor(t *testing.T) {
	t.Run("ZstdUncompressor 测试", func(t *testing.T) {
		u := ZstdUncompressor()
		assert.NotNil(t, u)
	})
}

func TestUncompressorType(t *testing.T) {
	t.Run("Uncompressor 类型测试", func(t *testing.T) {
		var u Uncompressor = GzipUncompressor()
		assert.NotNil(t, u)
	})
}

func TestUnzip(t *testing.T) {
	t.Run("Unzip 测试", func(t *testing.T) {
		t.Run("解压不存在的文件返回错误", func(t *testing.T) {
			err := Unzip("/nonexistent/zip/file.zip", "/tmp/test_unzip")
			assert.NotNil(t, err)
		})

		t.Run("解压有效的zip文件", func(t *testing.T) {
			tmpDir := t.TempDir()
			zipPath := filepath.Join(tmpDir, "test.zip")

			zf, err := os.Create(zipPath)
			assert.Nil(t, err)
			zw := zip.NewWriter(zf)
			w, err := zw.Create("test.txt")
			assert.Nil(t, err)
			_, err = w.Write([]byte("hello world"))
			assert.Nil(t, err)
			zw.Close()
			zf.Close()

			destDir := filepath.Join(tmpDir, "output")
			err = Unzip(zipPath, destDir)
			assert.Nil(t, err)

			content, err := os.ReadFile(filepath.Join(destDir, "test.txt"))
			assert.Nil(t, err)
			assert.Equal(t, "hello world", string(content))
		})
	})
}
