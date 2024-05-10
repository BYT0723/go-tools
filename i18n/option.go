package i18n

import (
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Option func(*LangSet)

func RegisterUnmarshaler(ft string, unmarshalFunc i18n.UnmarshalFunc) Option {
	return func(ls *LangSet) {
		ls.b.RegisterUnmarshalFunc(ft, unmarshalFunc)
	}
}

func LoadMsgFile(path string) Option {
	return func(ls *LangSet) {
		ls.b.MustLoadMessageFile(path)
	}
}

func LoadMsgFileFs(fs fs.FS, path string) Option {
	return func(ls *LangSet) {
		_, err := ls.b.LoadMessageFileFS(fs, path)
		if err != nil {
			panic(err)
		}
	}
}

func SetTemplateDelim(left, right string) Option {
	return func(ls *LangSet) {
		ls.tempParser.LeftDelim = left
		ls.tempParser.RightDelim = right
	}
}

func AddTemplateFunc(name string, f any) Option {
	return func(ls *LangSet) {
		ls.tempParser.Funcs[name] = f
	}
}
