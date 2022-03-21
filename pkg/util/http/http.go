package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Declared as variables due to ability to set custom values during build.
var (
	defaultTLSHandshakeTimeout = 2 * time.Second
	defaultMaxIdleConnections  = 10
	defaultMaxConnsPerHost     = 20
	defaultIdleConnTimeout     = 10 * time.Second
	defaultClientTimeout       = 5 * time.Second
)

// Client with default timeout 3 seconds
func NewClient(timeout time.Duration) *http.Client {
	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
			MaxIdleConns:        defaultMaxIdleConnections,
			MaxConnsPerHost:     defaultMaxConnsPerHost,
			IdleConnTimeout:     defaultIdleConnTimeout,
		},
	}

	if timeout == 0 {
		client.Timeout = defaultClientTimeout
	}

	return &client
}

func NewRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}
