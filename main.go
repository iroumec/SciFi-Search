package main

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	port    = ":8080"  // Puerto que se escucha.
	fileDir = "static" // Directorio relativo con los archivos estáticos. Relativo adonde se ejecuta go run.
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

var db = setupDB()

func main() {

	// Se crea un manejador (handler) de servidor de archivos.
	fileServer := http.FileServer(http.Dir(fileDir))

	// Se envuelve en un gzip middleware.
	http.Handle("/", gzipMiddleware(fileServer))

	registerHandlers()

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
