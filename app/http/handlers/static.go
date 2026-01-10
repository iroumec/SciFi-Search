package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/http/middlewares"
)

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

const (
	StaticFilesRoute     = "./static"
	StaticFilesDirectory = "/static/"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Registra la ruta que sirve los archivos estáticos.
func RegisterStatic() {

	// Se crea un manejador (handler) de servidor de archivos.
	fileServer := http.FileServer(http.Dir(StaticFilesRoute))

	// Se sirven archivos estáticos en /static/,
	// comprimidos en gzip si el navegador así lo acepta.
	http.Handle(
		StaticFilesDirectory,
		http.StripPrefix(
			StaticFilesDirectory,
			middlewares.GzipMiddleware(StaticFilesRoute, fileServer),
		),
	)
}

// ------------------------------------------------------------------------------------------------
