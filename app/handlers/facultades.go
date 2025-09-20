package handlers

import (
	"fmt"
	"net/http"
)

func manejarFacultades(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	manejarGETFacultad(w, r)
}

func manejarGETFacultad(w http.ResponseWriter, r *http.Request) {

	// Se lee la facultad de la query.
	facultad := r.URL.Query().Get("facultad")

	htmlPath := fmt.Sprintf("template/facultades/%s.html", facultad)

	renderizeTemplate(w, htmlPath, nil, nil)
}
