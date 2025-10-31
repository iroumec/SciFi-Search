package utils

// ------------------------------------------------------------------------------------------------

import (
	"log"
	"net/http"
	"time"
)

// ------------------------------------------------------------------------------------------------

// Muestra información de logging acerca de las solicitudes entrantes.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Tiempo de inicio, para medir la duración de la petición.
		start := time.Now()

		// Se imprime la información de la petición en la consola.
		// r.Method es el método HTTP (GET, POST, etc.).
		// r.URL.Path es la ruta solicitada.
		// r.RemoteAddr es la dirección IP del cliente.
		log.Printf("--> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Se invoca al siguiente handler en la cadena, el que
		// realmente trata la solicitud.
		next.ServeHTTP(w, r)

		// Finalizado el anterior handler, se muestra el tiempo que tardó en
		// responder la solicitud.
		log.Printf("<-- %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// ------------------------------------------------------------------------------------------------
