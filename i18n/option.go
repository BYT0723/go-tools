package i18n

import (
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Option func(*LangSet)

func WithUnmarshaler(ft string, unmarshalFunc i18n.UnmarshalFunc) Option {
	return func(ls *LangSet) {
		ls.b.RegisterUnmarshalFunc(ft, unmarshalFunc)
	}
}

func WithMsgFiles(paths ...string) Option {
	return func(ls *LangSet) {
		for _, p := range paths {
			ls.b.MustLoadMessageFile(p)
		}
	}
}

func WithMsgFilesFs(fs fs.FS, paths ...string) Option {
	return func(ls *LangSet) {
		for _, p := range paths {
			_, err := ls.b.LoadMessageFileFS(fs, p)
			if err != nil {
				panic(err)
			}
		}
	}
}

func WithTemplateDelim(left, right string) Option {
	return func(ls *LangSet) {
		ls.tempParser.LeftDelim = left
		ls.tempParser.RightDelim = right
	}
}

func WithTemplateFunc(name string, f any) Option {
	return func(ls *LangSet) {
		ls.tempParser.Funcs[name] = f
	}
}
