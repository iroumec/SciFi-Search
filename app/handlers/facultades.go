package handlers

import (
	"fmt"
	"net/http"
)

func manejarFacultades(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		manejarGETFacultad(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func manejarGETFacultad(w http.ResponseWriter, r *http.Request) {

	// Se lee la facultad de la query.
	facultad := r.URL.Query().Get("facultad")

	htmlPath := fmt.Sprintf("template/facultades/%s.html", facultad)

	renderizeTemplate(w, htmlPath, nil, nil)
}
