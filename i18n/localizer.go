package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
)

type Localizer struct {
	l          *i18n.Localizer
	tempParser *template.TextParser
}

func (l *Localizer) Msg(id string) (string, error) {
	return l.l.Localize(&i18n.LocalizeConfig{MessageID: id})
}

func (l *Localizer) MsgPlura(id string, plura int) (string, error) {
	return l.l.Localize(&i18n.LocalizeConfig{
		MessageID:   id,
		PluralCount: plura,
	})
}

func (l *Localizer) MsgWithParam(id string, param any) (string, error) {
	return l.l.Localize(&i18n.LocalizeConfig{
		MessageID:      id,
		TemplateData:   param,
		TemplateParser: l.tempParser,
	})
}

func (l *Localizer) MsgWithPluraParam(id string, plura int, param any) (string, error) {
	return l.l.Localize(&i18n.LocalizeConfig{
		MessageID:      id,
		TemplateData:   param,
		TemplateParser: l.tempParser,
	})
}
