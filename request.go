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

// NewGetRequest returns a new http.Request with GET method.
func NewGetRequest(ctx context.Context, url string) (*http.Request, error) {
	return NewRequest(ctx, http.MethodGet, url, http.NoBody)
}

// NewHeadRequest returns a new http.Request with HEAD method.
func NewHeadRequest(ctx context.Context, url string) (*http.Request, error) {
	return NewRequest(ctx, http.MethodHead, url, http.NoBody)
}

// NewPostRequest returns a new http.Request with POST method.
func NewPostRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodPost, url, body)
}

// NewPutRequest returns a new http.Request with PUT method.
func NewPutRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodPut, url, body)
}

// NewPatchRequest returns a new http.Request with PATCH method.
func NewPatchRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodPatch, url, body)
}

// NewDeleteRequest returns a new http.Request with DELETE method.
func NewDeleteRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodDelete, url, body)
}

// NewConnectRequest returns a new http.Request with CONNECT method.
func NewConnectRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodConnect, url, body)
}

// NewOptionsRequest returns a new http.Request with OPTIONS method.
func NewOptionsRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodOptions, url, body)
}

// NewTraceRequest returns a new http.Request with TRACE method.
func NewTraceRequest(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	return NewRequest(ctx, http.MethodTrace, url, body)
}

// MustNewRequest returns a new http.Request or panics on error.
func MustNewRequest(ctx context.Context, method, url string, body io.Reader) *http.Request {
	req, err := NewRequest(ctx, method, url, body)
	if err != nil {
		panic(fmt.Sprintf("httpx: must create request: %v", err))
	}
	return req
}

// MustGetRequest returns a new http.Request with GET method or panics on error.
func MustGetRequest(ctx context.Context, url string) *http.Request {
	req, err := NewGetRequest(ctx, url)
	if err != nil {
		panic(fmt.Sprintf("httpx: must create GET request: %v", err))
	}
	return req
}

// MustPostRequest returns a new http.Request with POST method or panics on error.
func MustPostRequest(ctx context.Context, url string, body io.Reader) *http.Request {
	req, err := NewPostRequest(ctx, url, body)
	if err != nil {
		panic(fmt.Sprintf("httpx: must create POST request: %v", err))
	}
	return req
}

// Bearer header with a give token.
func Bearer(token string) string {
	return "Bearer " + token
}
