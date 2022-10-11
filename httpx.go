package httpx

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Error for a request.
type Error struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
	Type  string `json:"type,omitempty"`
}

// ErrorResponse return and error wrapped into JSON.
func ErrorResponse(w http.ResponseWriter, code int, err error) {
	data := Error{
		Code: code,
		// TODO: add Type
	}

	if err != nil {
		data.Error = err.Error()
	}

	raw, errMarsh := json.Marshal(data)
	if errMarsh != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("JSON marshal failed: %v", errMarsh)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(prettyJSON(raw))
}

// TODO: enable by compile flag?
func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", "  ")
	return out.Bytes()
}
