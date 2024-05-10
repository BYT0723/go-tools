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

func (l *Localizer) MsgWithParam(id string, param any) (string, error) {
	return l.l.Localize(&i18n.LocalizeConfig{
		MessageID:      id,
		TemplateData:   param,
		TemplateParser: l.tempParser,
	})
}

// Plura必须为int或float类型，否则会出现panic
// Plura为句子中包含的数量值，i18n通过这个判断使用哪种类型的句式
func (l *Localizer) MsgWithPluraParam(id string, param, plura any) (string, error) {
	return l.l.Localize(&i18n.LocalizeConfig{
		MessageID:      id,
		TemplateData:   param,
		PluralCount:    plura,
		TemplateParser: l.tempParser,
	})
}
