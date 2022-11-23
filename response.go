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
