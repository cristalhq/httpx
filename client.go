package httpx

import "net/http"

// Doer is an interface for http.Client.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}
