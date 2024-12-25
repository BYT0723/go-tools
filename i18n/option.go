package i18n

import (
	"io/fs"
	"maps"

	"github.com/Masterminds/sprig/v3"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Option func(*langset)

// set the Unmarshaler to parse the language file for the specified file type
func WithUnmarshaler(ft string, unmarshalFunc i18n.UnmarshalFunc) Option {
	return func(ls *langset) {
		ls.b.RegisterUnmarshalFunc(ft, unmarshalFunc)
	}
}

func WithMsgFiles(paths ...string) Option {
	return func(ls *langset) {
		for _, p := range paths {
			ls.b.MustLoadMessageFile(p)
		}
	}
}

func WithMsgFilesFs(fs fs.FS, paths ...string) Option {
	return func(ls *langset) {
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
	return func(ls *langset) {
		ls.tempParser.LeftDelim = left
		ls.tempParser.RightDelim = right
	}
}

// add a template function
func WithTemplateFunc(name string, f any) Option {
	return func(ls *langset) {
		ls.tempParser.Funcs[name] = f
	}
}

// add template functions
// eg: WithTemplateFuncs(sprig.FuncMap())
// import [sprig](https://github.com/Masterminds/sprig) template funcs
func WithTemplateFuncs(fs map[string]any) Option {
	return func(ls *langset) {
		maps.Copy(ls.tempParser.Funcs, fs)
	}
}

// import [sprig](https://github.com/Masterminds/sprig) template funcs
func WithDefaultTemplateFuncs() Option {
	return func(ls *langset) {
		maps.Copy(ls.tempParser.Funcs, sprig.TxtFuncMap())
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
	return func(ls *langset) {
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
