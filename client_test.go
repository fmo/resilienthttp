package resilienthttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("response comes here"))
	}))
	defer ts.Close()

	httpReq, _ := http.NewRequest("GET", ts.URL, nil)
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
