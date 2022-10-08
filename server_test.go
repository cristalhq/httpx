package httpx

import (
	"context"
	"net/http"
	"net/http/httputil"
	"testing"
)

func TestServer(t *testing.T) {
	cfg := &ServerConfig{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(207)
			w.Header().Del("Date")
			w.Write([]byte("test"))
		}),
	}
	srv, err := NewServer(cfg)
	if err != nil {
		t.Fatal(err)
	}

	srv.srv.ErrorLog = nil

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		srv.Run(ctx)
	}()

	resp, err := http.Get("http://:8080")
	if err != nil {
		t.Fatal(err)
	}
	body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	cancel()

	if _, err := http.Get("http://:8080"); err == nil {
		t.Fatal("should fail")
	}
}
