package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	resp, err := Get(srv.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"ok":true}`, string(resp.Body))
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	resp, err := Post(srv.URL, WithPayload(`{"key":"value"}`))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestDoWithHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc := r.Header.Get("Accept-Encoding")
		if enc != "gzip" {
			t.Errorf("expected gzip, got %s", enc)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	resp, err := Get(srv.URL, WithHeader(http.Header{
		"Accept-Encoding": []string{"gzip"},
	}))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDownload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer srv.Close()

	path := t.TempDir() + "/download_test.txt"
	err := Download(srv.URL, path)
	assert.Nil(t, err)
}

func TestResponseFields(t *testing.T) {
	resp := &Response{
		Code:   200,
		Header: http.Header{"Content-Type": []string{"application/json"}},
		Body:   []byte("test"),
	}
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, "test", string(resp.Body))
}
