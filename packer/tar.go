package packer

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func Untar(src string, dst string, uncompressors ...Uncompressor) (err error) {
	defer func() {
		if err != nil {
			os.RemoveAll(dst)
		}
	}()

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	var reader io.Reader = f

	for _, c := range uncompressors {
		r, err := c(reader)
		if err != nil {
			return err
		}
		reader = r
	}

	tr := tar.NewReader(reader)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dst, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
