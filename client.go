package httpx

import (
	"net/http"
	"time"
)

// Doer is an interface for http.Client.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewClient return http.Client with a sane default.
func NewClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
