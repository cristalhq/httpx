package httpx

import (
	"encoding/json"
	"net/http"
)

// Error for a request.
type Error struct {
	Code    int    `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"msg,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

// ErrorResponse return and error wrapped into JSON.
func ErrorResponse(w http.ResponseWriter, code int, err error) {
	if err == nil {
		w.WriteHeader(code)
		return
	}

	var resp *Error
	if e, ok := err.(*Error); ok {
		resp = e
	} else {
		resp = &Error{
			Message: err.Error(),
		}
	}
	resp.Code = code

	// TODO: indent under compile flag?
	raw, _ := json.MarshalIndent(resp, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(raw)
}
