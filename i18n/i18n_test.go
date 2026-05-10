package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWithParam(t *testing.T) {
	Convey("WithParam 测试", t, func() {
		opt := WithParam(map[string]string{"name": "world"})
		So(opt, ShouldNotBeNil)
	})
}

func TestWithPluralCount(t *testing.T) {
	Convey("WithPluralCount 测试", t, func() {
		opt := WithPluralCount(5)
		So(opt, ShouldNotBeNil)
	})
}

func TestLocaleOption(t *testing.T) {
	Convey("LocaleOption 类型测试", t, func() {
		var opt LocaleOption = func(lc *i18n.LocalizeConfig) {}
		So(opt, ShouldNotBeNil)
	})
}

func TestOptionType(t *testing.T) {
	Convey("Option 类型测试", t, func() {
		var opt Option = func(ls *langset) {}
		So(opt, ShouldNotBeNil)
	})
}
