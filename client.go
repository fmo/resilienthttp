package resilienthttp

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var (
	defaultRetryMax = 4
	defaultBackoff  = 1 * time.Second
)

// Client wraps the standard library's http.Client.
// Adds retry capability for more resilient HTTP requests.
type Client struct {
	HTTPClient *http.Client
	RetryMax   int
	Backoff    BackoffFunc
}

type BackoffFunc func(attempt int)

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
		RetryMax:   defaultRetryMax,
		Backoff: func(attempt int) {
			time.Sleep(defaultBackoff * time.Duration(attempt))
		},
	}
}

type Request struct {
	*http.Request
}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return &Request{Request: r}, nil
}

func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	return &Request{Request: r}, err
}

func Get(url string) (*http.Response, error) {
	c := NewClient()
	r, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(r)
}

func (c *Client) Do(req *Request) (*http.Response, error) {
	var response *http.Response
	var doErr error

	for i := 1; ; i++ {
		response, doErr = c.HTTPClient.Do(req.Request)
		if doErr != nil {
			if i >= c.RetryMax {
				return nil, doErr
			} else {
				c.Backoff(i)
				slog.Error("request failed", "attempt", i, "error", doErr)
				continue
			}
		}
		break
	}

	return response, nil
}
