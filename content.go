package httpx

import (
	"path/filepath"
)

// ContentTypeByExt or 'application/octet-stream' if not found.
func ContentTypeByExt(file string) string {
	if typ := ctypes[filepath.Ext(file)]; typ != "" {
		return typ
	}
	return ctypes[".bin"]
}

var ctypes = map[string]string{
	".avif": "image/avif",
	".bin":  "application/octet-stream",
	".css":  "text/css; charset=utf-8",
	".csv":  "text/csv; charset=utf-8",
	".gif":  "image/gif",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".ico":  "image/png",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".js":   "text/javascript; charset=utf-8",
	".json": "application/json; charset=utf-8",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".txt":  "text/plain; charset=utf-8",
	".webp": "image/webp",
	".xml":  "application/xml",
	".zip":  "application/zip",
}
