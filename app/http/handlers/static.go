package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/http/middlewares"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Registra la ruta que sirve los archivos estáticos.
func RegisterStatic() {

	// Se crea un manejador (handler) de servidor de archivos.
	fileServer := http.FileServer(http.Dir(fileDir))

	// Se sirven archivos estáticos en /static/,
	// comprimidos en gzip si el navegador así lo acepta.
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			middlewares.GzipMiddleware(fileDir, fileServer),
		),
	)
}

// ------------------------------------------------------------------------------------------------
