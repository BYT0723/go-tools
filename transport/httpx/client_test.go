package httpx

import (
	"fmt"
	"net/http"
	"testing"
)

func TestClientDo(t *testing.T) {
	resp, err := Get("http://localhost:9090/metrics", WithHeader(http.Header{
		"Accept-Encoding": []string{"gzip", "br", "flate"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("resp.Code: %v\n", resp.Code)
	fmt.Printf("resp.Header: %v\n", resp.Header)
	fmt.Printf("resp.Body: %s\n", resp.Body)
}
