package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// NewRequest returns a new http.Request.
func NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// MustNewRequest returns a new http.Request or panics on error.
func MustNewRequest(ctx context.Context, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		panic(fmt.Sprintf("httpx: create request: %v", err))
	}
	return req
}
