package httpx

import (
	"io"
	"net/http"
)

func ReturnOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func ReturnOKJSON(w http.ResponseWriter, data interface{}) {
	MarshalResponse(w, http.StatusOK, data)
}

func ReturnNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func ReturnRedirect(w http.ResponseWriter) {
	w.WriteHeader(http.StatusPermanentRedirect)
}

func ReturnBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
}

// DiscardResponseBody reads and closes http.Response.Body.
func DiscardResponseBody(resp *http.Response) {
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func Is1xx(code int) bool { return code >= 100 && code < 200 }
func Is2xx(code int) bool { return code >= 200 && code < 300 }
func Is3xx(code int) bool { return code >= 300 && code < 400 }
func Is4xx(code int) bool { return code >= 400 && code < 500 }
func Is5xx(code int) bool { return code >= 500 && code < 600 }
