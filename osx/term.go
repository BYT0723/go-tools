package osx

import (
	"errors"
	"os"

	"golang.org/x/term"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

func GetTermSize() (w, h int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func CharmapDecode(cp uint32, raw []byte) ([]byte, error) {
	var decoder *encoding.Decoder

	switch cp {
	case 437: // US
		decoder = charmap.CodePage437.NewDecoder()
	case 850: // Latin-I
		decoder = charmap.CodePage850.NewDecoder()
	case 852: // Latin-II
		decoder = charmap.CodePage852.NewDecoder()
	case 866: // cyrillic
		decoder = charmap.CodePage866.NewDecoder()
	case 932: // shift-jis
		decoder = japanese.ShiftJIS.NewDecoder()
	case 936: // gbk
		decoder = simplifiedchinese.GBK.NewDecoder()
	case 949: // euc-kr
		decoder = korean.EUCKR.NewDecoder()
	case 950: // big5
		decoder = traditionalchinese.Big5.NewDecoder()
	case 65001: // utf8
	default:
		return raw, errors.New("unknown OEM Code Page")
	}

	if decoder != nil {
		return decoder.Bytes(raw)
	}
	return raw, nil
}
