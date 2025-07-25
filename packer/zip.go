package packer

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(src string, dst string) (err error) {
	defer func() {
		if err != nil {
			os.RemoveAll(dst)
		}
	}()

	rc, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	defer rc.Close()

	if err = os.MkdirAll(dst, os.ModePerm); err != nil {
		return
	}

	for _, f := range rc.File {
		fp := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fp, os.ModePerm); err != nil {
				return
			}
			continue
		}

		dstf, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer dstf.Close()

		fr, err := f.Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		if _, err := io.Copy(dstf, fr); err != nil {
			return err
		}
	}
	return nil
}
