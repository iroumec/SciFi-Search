package middlewares

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"compress/gzip"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Structures
// ------------------------------------------------------------------------------------------------

// It wraps `ResponseWriter` so the response can be compressed.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// ------------------------------------------------------------------------------------------------

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Compresses the response if the client accepts it and a file exists.
func GzipMiddleware(fileDir string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Does the client accept gzip?
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Does the file exists?
		path := filepath.Join(fileDir, filepath.Clean(r.URL.Path))
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {

			// File doesn't exist or it is a directory -> No compression.
			next.ServeHTTP(w, r)
			return
		}

		// Compression via gzip is applied to the file.
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		// The `ResponseWriter` is wrapped so it can be compressed.
		gzw := gzipResponseWriter{ResponseWriter: w, writer: gz}
		next.ServeHTTP(&gzw, r)
	})
}

// ------------------------------------------------------------------------------------------------
