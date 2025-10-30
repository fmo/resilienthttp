package resilienthttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoWith404(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(404)
    }))

    c := NewClient()
    req, _ := http.NewRequest("GET", ts.URL, nil)

    resilientReq := &Request{Request: req}

    resp, _ := c.Do(resilientReq)

    if resp.StatusCode != 404 {
        t.Errorf("unexpected status code")
    }
}


func TestDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("response comes here"))
	}))
	defer ts.Close()

	httpReq, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	resilientReq := &Request{Request: httpReq}

	client := NewClient()

	resp, err := client.Do(resilientReq)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if string(body) != "response comes here" {
		t.Error("no matching response")
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))

	resp, err := Get(ts.URL)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if string(body) != "hello world" {
		t.Error("not expected response")
	}
}
