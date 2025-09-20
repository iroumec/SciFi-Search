package handlers

import (
	"fmt"
	"net/http"
)

func registerHandlersFacultades() {

	http.HandleFunc("/facultades", manejarFacultades)
}

// ------------------------------------------------------------------------------------------------
// Registro de Handlers de Facultades
// ------------------------------------------------------------------------------------------------

func manejarFacultades(w http.ResponseWriter, r *http.Request) {

	// Si el método no es GET, error.
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	manejarGETFacultad(w, r)
}

// ------------------------------------------------------------------------------------------------

func manejarGETFacultad(w http.ResponseWriter, r *http.Request) {

	// Se lee la facultad de la query.
	facultad := r.URL.Query().Get("facultad")

	// Se establece el path al html.
	htmlPath := fmt.Sprintf("template/facultades/%s.html", facultad)

	// Se renderiza la plantilla.
	renderizeTemplate(w, htmlPath, nil, nil)
}
