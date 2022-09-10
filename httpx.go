package httpx

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Error for a request.
type Error struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse return and error wrapped into JSON.
func ErrorResponse(w http.ResponseWriter, code int, err error) {
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(code)
		return
	}

	msg := err.Error()

	data := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  code,
		Error: msg,
	}

	raw, errMarsh := json.Marshal(data)
	if errMarsh != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("JSON marshal failed: %v", errMarsh)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = w.Write(prettyJSON(raw))
}

// TODO: enable by compile flag
func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", "  ")
	return out.Bytes()
}
