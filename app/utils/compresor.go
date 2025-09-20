package utils

import (
	"compress/gzip"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/*
Comprime la respuesta si el cliente acepta gzip y el archivo existe.
*/
func GzipMiddleware(fileDir string, next http.Handler) http.Handler {
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

			// Si el archivo no existe o es un directorio,
			// se sirve la solicitud sin compresión.
			next.ServeHTTP(w, r)
			return
		}

		// Si el archivo existe, se aplica la compresión mediante gzip.
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Se envuelve el ResponseWriter para comprimir la salida.
		gzw := gzipResponseWriter{ResponseWriter: w, writer: gz}
		next.ServeHTTP(&gzw, r)
	})
}

// ------------------------------------------------------------------------------------------------

// Se envuelve ResponseWriter para comprimir la salida.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// ------------------------------------------------------------------------------------------------

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}
