package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
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
// If the err is httpx.Error, code parameter is ignoted and httpx.Error.Code is used.
// If the err is nil just status code is returned.
func ErrorResponse(w http.ResponseWriter, code int, err error) {
	if err == nil {
		w.WriteHeader(code)
		return
	}

	var resp *Error
	if e, ok := err.(*Error); ok {
		resp = e
		code = resp.Code
	} else {
		resp = &Error{
			Code:    code,
			Message: err.Error(),
		}
	}

	// TODO: indent under compile flag?
	raw, _ := json.MarshalIndent(resp, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(raw)
}

func DumpRequest(r *http.Request) string {
	raw, err := httputil.DumpRequest(r, true)
	if err != nil {
		return fmt.Errorf("dump error: %#v", err).Error()
	}
	return fmt.Sprintf("%#v", string(raw))
}

// NoopHandler just return 200 OK. Useful for prototyping.
func NoopHandler(http.ResponseWriter, *http.Request) {}
