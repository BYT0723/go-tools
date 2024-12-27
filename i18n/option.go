package i18n

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func WithMsgDir(dirs ...string) Option {
	return func(ls *langset) {
		for _, dir := range dirs {
			de, err := os.ReadDir(dir)
			if err != nil {
				panic(err)
			}
			for _, d := range de {
				ls.b.MustLoadMessageFile(filepath.Join(dir, d.Name()))
			}
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

func WithMsgFilesDirFs(fs fs.ReadDirFS, dirs ...string) Option {
	return func(ls *langset) {
		for _, dir := range dirs {
			de, err := fs.ReadDir(dir)
			if err != nil {
				panic(err)
			}

			for _, d := range de {
				_, err := ls.b.LoadMessageFileFS(fs, filepath.Join(dir, d.Name()))
				if err != nil {
					panic(err)
				}
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
