package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
)

type (
	Localizer struct {
		l          *i18n.Localizer
		tempParser *template.TextParser
	}
	LocaleOption func(*i18n.LocalizeConfig)
)

func (l *Localizer) Msg(id string, opts ...LocaleOption) (string, error) {
	lc := &i18n.LocalizeConfig{MessageID: id, TemplateParser: l.tempParser}
	for _, opt := range opts {
		opt(lc)
	}
	return l.l.Localize(lc)
}

func (l *Localizer) MustMsg(id string, opts ...LocaleOption) string {
	s, err := l.Msg(id, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

func WithParam(param any) LocaleOption {
	return func(lc *i18n.LocalizeConfig) {
		lc.TemplateData = param
	}
}

// Plura的一些使用场景
// 英语：
// 复数形式：通常情况下，当数字为 1 时使用单数形式，其他情况使用复数形式。
// 示例：
// 1 apple
// 2 apples
// 0 apples
//
// 西班牙语：
// 复数形式：通常情况下，当数字为 1 时使用单数形式，其他情况使用复数形式。
// 示例：
// 1 manzana (苹果)
// 2 manzanas
// 0 manzanas
//
// 法语：
//
// 复数形式：通常情况下，当数字为 0 或 1 时使用单数形式，其他情况使用复数形式。
// 示例：
// 1 pomme (苹果)
// 2 pommes
// 0 pommes
//
// 德语：
//
// 复数形式：通常情况下，当数字为 1 时使用单数形式，其他情况使用复数形式。
// 示例：
// 1 Apfel (苹果)
// 2 Äpfel
// 0 Äpfel
//
// 俄语：
//
// 复数形式：通常情况下，1 对应单数形式，2-4 对应复数形式，5 及以上对应另一种复数形式。
// 示例：
// 1 яблоко (苹果)
// 2 яблока
// 5 яблок
//
// 意大利语：
//
// 复数形式：通常情况下，1 对应单数形式，其他情况使用复数形式。
// 示例：
// 1 mela (苹果)
// 2 mele
// 0 mele
//
// 葡萄牙语：
//
// 复数形式：通常情况下，当数字为 1 时使用单数形式，其他情况使用复数形式。
// 示例：
// 1 maçã (苹果)
// 2 maçãs
// 0 maçãs
//
// 阿拉伯语：
//
// 复数形式：规则较复杂，通常与数字的形式、性别和状态等有关。
// 示例：
// ١ تفاحة (苹果)
// ٢ تفاحة
// ٥ تفاحات
// -------------------------------------------
// Plura必须为int或float类型，否则会出现panic
// Plura为句子中包含的数量值，i18n通过这个判断使用哪种类型的句式
func WithPluralCount(pluralCount any) LocaleOption {
	return func(lc *i18n.LocalizeConfig) {
		lc.PluralCount = pluralCount
	}
}
