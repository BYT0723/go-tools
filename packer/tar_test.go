package packer

import (
	"os"
	"testing"
)

func TestUntar(t *testing.T) {
	uncompressDir := "./uncompress"
	defer os.RemoveAll(uncompressDir)

	type args struct {
		src           string
		dst           string
		uncompressors []Uncompressor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "uncompress tar file",
			args: args{
				src:           "./compress/test.tar",
				dst:           "./uncompress/tar/",
				uncompressors: []Uncompressor{},
			},
			wantErr: false,
		},
		{
			name: "uncompress tar.gz file",
			args: args{
				src:           "./compress/test.tar.gz",
				dst:           "./uncompress/tar.gz/",
				uncompressors: []Uncompressor{GzipUncompressor()},
			},
			wantErr: false,
		},
		{
			name: "uncompress tar.bz2 file",
			args: args{
				src:           "./compress/test.tar.bz2",
				dst:           "./uncompress/tar.bz2/",
				uncompressors: []Uncompressor{Bzip2Uncompressor()},
			},
			wantErr: false,
		},
		{
			name: "uncompress tar.xz file",
			args: args{
				src:           "./compress/test.tar.xz",
				dst:           "./uncompress/tar.xz/",
				uncompressors: []Uncompressor{XzUncompressor()},
			},
			wantErr: false,
		},
		{
			name: "uncompress tar.zst file",
			args: args{
				src:           "./compress/test.tar.zst",
				dst:           "./uncompress/tar.zst/",
				uncompressors: []Uncompressor{ZstdUncompressor()},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Untar(tt.args.src, tt.args.dst, tt.args.uncompressors...); (err != nil) != tt.wantErr {
				t.Errorf("Untar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
