package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/http/middlewares"
)

// ------------------------------------------------------------------------------------------------
// Constants
// ------------------------------------------------------------------------------------------------

const (
	StaticFilesRoute     = "./static"
	StaticFilesDirectory = "/static/"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Endpoints.
func RegisterStatic() {

	// Static files routes handler.
	fileServer := http.FileServer(http.Dir(StaticFilesRoute))

	// The files are compressed if the browser accepts it.
	http.Handle(
		StaticFilesDirectory,
		http.StripPrefix(
			StaticFilesDirectory,
			middlewares.GzipMiddleware(StaticFilesRoute, fileServer),
		),
	)
}

// ------------------------------------------------------------------------------------------------
