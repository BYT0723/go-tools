package http

import (
	"context"
	"testing"
)

func TestClientDo(t *testing.T) {
	cli := NewClient()

	cli.Get(context.Background(), "http://www.baidu.com", nil, map[string]any{
		"serial[]": []string{"aaa", "bbbb", "ccc"},
	})
}
