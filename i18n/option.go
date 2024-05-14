package i18n

import (
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
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

// set delimiters in template
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

type Message struct {
	Language string
	ID       string
	Zero     string
	One      string
	Two      string
	Few      string
	Many     string
	Other    string
}

func WithMessages(ms ...*Message) Option {
	return func(ls *LangSet) {
		for _, m := range ms {
			lang, err := language.Parse(m.Language)
			if err != nil {
				panic(err)
			}
			if err := ls.b.AddMessages(lang, &i18n.Message{
				ID:    m.ID,
				Zero:  m.Zero,
				One:   m.One,
				Two:   m.Two,
				Few:   m.Few,
				Many:  m.Many,
				Other: m.Other,
			}); err != nil {
				panic(err)
			}
		}
	}
}
