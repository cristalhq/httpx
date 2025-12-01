package httpx

import (
	"io"
	"io/fs"
	"log"
	"net/http"
)

func ResponseOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func ResponseJSON(w http.ResponseWriter, data any) {
	MarshalResponse(w, http.StatusOK, data)
}

// DiscardResponse reads and closes [http.Response.Body].
func DiscardResponse(resp *http.Response) {
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func Is1xx(code int) bool { return code >= 100 && code < 200 }
func Is2xx(code int) bool { return code >= 200 && code < 300 }
func Is3xx(code int) bool { return code >= 300 && code < 400 }
func Is4xx(code int) bool { return code >= 400 && code < 500 }
func Is5xx(code int) bool { return code >= 500 && code < 600 }

func ServeFile(fs fs.FS, file string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		f, err := fs.Open(file)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", ContentTypeByExt(file))
		w.Header().Set("Cache-Control", "max-age=31536000, public")
		w.WriteHeader(http.StatusOK)

		if _, err := io.Copy(w, f); err != nil {
			log.Printf("error sending file: %v", err)
		}
	}
}
