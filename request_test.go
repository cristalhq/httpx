package httpx

import (
	"context"
	"net/http"
	"testing"
)

func TestRequestIP(t *testing.T) {
	const defaultAddr = "128.10.20.30:50"

	testCases := []struct {
		headers http.Header
		want    string
	}{
		{
			want: "128.10.20.30",
		},
		{
			headers: http.Header{
				"Cf-Connecting-Ip": []string{"192.0.2.5"},
			},
			want: "192.0.2.5",
		},
		{
			headers: http.Header{
				"True-Client-IP": []string{"192.0.20.50"},
				"X-Real-Ip":      []string{"192.0.2.5"},
			},
			want: "192.0.2.5",
		},
		{
			headers: http.Header{
				"X-Real-Ip":       []string{"20.20.20.42"},
				"X-Forwarded-For": []string{"203.0.113.42"},
			},
			want: "20.20.20.42",
		},
		{
			headers: http.Header{
				"X-Forwarded-For": []string{"203.0.113.42"},
			},
			want: "203.0.113.42",
		},
		{
			headers: http.Header{
				"X-Forwarded-For": []string{"203.0.113.42, 198.51.100.17, 192.0.2.5"},
			},
			want: "203.0.113.42",
		},
	}

	for _, tc := range testCases {
		r := MustGetRequest(context.Background(), "/")
		r.Header = tc.headers
		r.RemoteAddr = defaultAddr

		addr := RequestIP(r)

		if have := addr.String(); have != tc.want {
			t.Errorf("\nhave: %v\nwant: %v", have, tc.want)
		}
	}
}
