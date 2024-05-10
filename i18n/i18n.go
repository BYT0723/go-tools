package i18n

import (
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
	"golang.org/x/text/language"
)

var (
	ls   *LangSet
	once sync.Once
)

func Init(opts ...Option) {
	once.Do(func() {
		ls = &LangSet{
			b:          i18n.NewBundle(language.English),
			tempParser: &template.TextParser{Funcs: map[string]any{}},
		}

		for _, opt := range opts {
			opt(ls)
		}
	})
}

func GetLocalizer(lang string) *Localizer {
	return ls.GetLocalizer(lang)
}
