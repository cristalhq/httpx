package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MarshalResponse writes an object to the http.ResponseWriter or a fallback error.
func MarshalResponse(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := strings.ReplaceAll(err.Error(), `"`, `\"`)
		fmt.Fprintf(w, `{"error":"%s"}`, msg)
		return
	}

	w.Write(data)
	w.WriteHeader(status)
}

const maxBodyBytes = 4096

// UnmarshalRequest as JSON with well defined error handling.
func UnmarshalRequest(w http.ResponseWriter, r *http.Request, data interface{}) (code int, err error) {
	if ct := r.Header.Get("content-type"); !strings.HasPrefix(ct, "application/json") {
		return http.StatusUnsupportedMediaType, fmt.Errorf("Content-Type header is not application/json: %s", ct)
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&data); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxErr):
			return http.StatusBadRequest, fmt.Errorf("malformed json at position %d", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return http.StatusBadRequest, errors.New("malformed json")
		case errors.As(err, &unmarshalError):
			return http.StatusBadRequest, fmt.Errorf("invalid value %q at position %d", unmarshalError.Field, unmarshalError.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return http.StatusBadRequest, errors.New("unknown field " + fieldName)
		case errors.Is(err, io.EOF):
			return http.StatusBadRequest, errors.New("body must not be empty")
		case err.Error() == "http: request body too large":
			return http.StatusRequestEntityTooLarge, err
		default:
			return http.StatusInternalServerError, fmt.Errorf("failed to decode json: %w", err)
		}
	}

	if d.More() {
		return http.StatusBadRequest, errors.New("body must contain only one JSON object")
	}
	return http.StatusOK, nil
}
