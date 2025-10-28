package resilienthttp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var (
	defaultRetryMax = 4
	defaultBackoff  = 1 * time.Second
	maxBackoff      = 16 * time.Second
)

type CheckRetry func(res *http.Response) bool

type BackoffFunc func(attempt int)

// Client extens http.Client with configurable retry logic for improved request resilience
type Client struct {
	HTTPClient *http.Client
	RetryMax   int
	CheckRetry CheckRetry
	Backoff    BackoffFunc
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
		// Limit the number of retries to give the service time to recover
		RetryMax: defaultRetryMax,
		Backoff: func(attempt int) {
			backoffTime := time.Duration(math.Min(float64(defaultBackoff)*math.Pow(2, float64(attempt)), float64(maxBackoff)))
			jitter := time.Duration(rand.Float64() * float64(backoffTime) * 0.5)
			slog.Info("exponential backoff", "time", backoffTime+jitter)
			time.Sleep(backoffTime + jitter)
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

func (c *Client) Do(req *Request) (*http.Response, error) {
	var response *http.Response
	var err error

	for i := 0; i < c.RetryMax; i++ {
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
			slog.Error("request failed", "attempt", i+1)
			if response != nil {
				response.Body.Close()
			}
			continue
		}
	}

	slog.Error("all retries failed", "retries", c.RetryMax, "last_status", response.StatusCode)
	return response, fmt.Errorf("request failed after %d retries: last status=%d", c.RetryMax, response.StatusCode)
}

func Get(url string) (*http.Response, error) {
	c := NewClient()
	r, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(r)
}
