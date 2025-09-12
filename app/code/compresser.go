package code

import (
	"compress/gzip"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// gzipMiddleware comprime la respuesta si el cliente acepta gzip y el archivo existe.
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Se verifica si el cliente acepta gzip.
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Se verifica si el archivo existe.
		path := filepath.Join(fileDir, filepath.Clean(r.URL.Path))
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {

			// Archivo no existe o es directorio: se sirve sin compresión.
			next.ServeHTTP(w, r)
			return
		}

		// Archivo existe: se aplica gzip.
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Se envuelve ResponseWriter para comprimir la salida.
		gzw := gzipResponseWriter{ResponseWriter: w, writer: gz}
		next.ServeHTTP(&gzw, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}
