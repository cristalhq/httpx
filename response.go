package httpx

import (
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
