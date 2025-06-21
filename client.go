package httpx

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

// Doer is a backward compatible definition for [Client].
type Doer = Client

// Client is an interface for [http.Client].
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewClient returns [http.Client] with a sane defaults and non-shared [http.Transport].
// See [NewTransport] for more info.
func NewClient() *http.Client {
	return &http.Client{
		Transport: NewTransport(),
		Timeout:   5 * time.Second,
	}
}

// NewPooledClient returns [http.Client] which will be used for the same host(s).
func NewPooledClient() *http.Client {
	return &http.Client{
		Transport: NewPooledTransport(),
		Timeout:   5 * time.Second,
	}
}

// NewTransport returns [http.Transport] with idle connections and keepalives disabled.
func NewTransport() *http.Transport {
	transport := NewPooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// NewPooledTransport returns [http.Transport] which will be used for the same host(s).
func NewPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}
	return transport
}
