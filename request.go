package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

var NoBody = http.NoBody

// NewRequest returns a new http.Request.
func NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// NewRequest returns a new http.Request.
func NewGetRequest(ctx context.Context, url string) (*http.Request, error) {
	return NewRequest(ctx, http.MethodGet, url, http.NoBody)
}

// NewRequest returns a new http.Request.
func NewPostRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodPost, url, body)
}

// MustNewRequest returns a new http.Request or panics on error.
func MustNewRequest(ctx context.Context, method, url string, body io.Reader) *http.Request {
	req, err := NewRequest(ctx, method, url, body)
	if err != nil {
		panic(fmt.Sprintf("httpx: must create request: %v", err))
	}
	return req
}
