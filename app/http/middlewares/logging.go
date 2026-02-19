package middlewares

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"log"
	"net/http"
	"time"
)

// ------------------------------------------------------------------------------------------------
// Structures
// ------------------------------------------------------------------------------------------------

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// ------------------------------------------------------------------------------------------------

// WriteHeader implements the interface http.ResponseWriter.
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// It shows logging information about the entering requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Start time, so the time of completion of the request can be measure.
		start := time.Now()

		// The `ResponseWriter` is wrapped so the external middlewares responses
		// are not modified.
		wrappedResponseWriter := newLoggingResponseWriter(w)

		// Request information is printed in the console.
		// r.Method -> HTTP method (GET, POST, etc.).
		// r.URL.Path -> Path requested.
		// r.RemoteAddr -> Client's IP direction.
		log.Printf("--> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// The next handlers is invoked.
		next.ServeHTTP(wrappedResponseWriter, r)

		// The request has finished.
		// The total time passed is shown.
		log.Printf("<-- %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func newLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	// Status is initialized, by default, with HTTP 200 OK.
	return &LoggingResponseWriter{w, http.StatusOK}
}

// ------------------------------------------------------------------------------------------------
