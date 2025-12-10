package spider

import (
	"testing"
)

func TestGetRandUserAgent(t *testing.T) {
	t.Run("empty_check", func(t *testing.T) {
		ua := GetRandUserAgent()
		if ua == "" {
			t.Error("ua is empty")
		}
	})
}
