package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
)

type langset struct {
	b          *i18n.Bundle
	tempParser *template.TextParser
}

func (l *langset) GetLocalizer(lang string) *Localizer {
	return &Localizer{
		l:          i18n.NewLocalizer(l.b, lang),
		tempParser: l.tempParser,
	}
}
