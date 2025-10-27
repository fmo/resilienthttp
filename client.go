package resilienthttp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var (
	defaultRetryMax = 4
	defaultBackoff  = 1 * time.Second
)

type CheckRetry func(res *http.Response) bool

// Client wraps the standard library's http.Client.
// Adds retry capability for more resilient HTTP requests.
type Client struct {
	HTTPClient *http.Client
	RetryMax   int
	CheckRetry CheckRetry
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
		CheckRetry: func(res *http.Response) bool {
			return res.StatusCode >= 500
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
	var err error

	for i := 1; i <= c.RetryMax; i++ {
		response, err = c.HTTPClient.Do(req.Request)
		if err != nil {
			slog.Error("request failed", "err", err)
			return nil, err
		}

		if response.StatusCode < 500 {
			return response, nil
		}

		if c.CheckRetry(response) {
			c.Backoff(i)
			slog.Error("request failed", "attempt", i)
			if response != nil {
				response.Body.Close()
			}
			continue
		}
	}

    slog.Error("all retries failed", "retries", c.RetryMax, "last_status", response.StatusCode)
	return response, fmt.Errorf("request failed after %d retries: last status=%d", c.RetryMax, response.StatusCode)
}
