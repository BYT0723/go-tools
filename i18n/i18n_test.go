package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/stretchr/testify/assert"
)

func TestWithParam(t *testing.T) {
	t.Run("WithParam 测试", func(t *testing.T) {
		opt := WithParam(map[string]string{"name": "world"})
		assert.NotNil(t, opt)
	})
}

func TestWithPluralCount(t *testing.T) {
	t.Run("WithPluralCount 测试", func(t *testing.T) {
		opt := WithPluralCount(5)
		assert.NotNil(t, opt)
	})
}

func TestLocaleOption(t *testing.T) {
	t.Run("LocaleOption 类型测试", func(t *testing.T) {
		var opt LocaleOption = func(lc *i18n.LocalizeConfig) {}
		assert.NotNil(t, opt)
	})
}

func TestOptionType(t *testing.T) {
	t.Run("Option 类型测试", func(t *testing.T) {
		var opt Option = func(ls *langset) {}
		assert.NotNil(t, opt)
	})
}
