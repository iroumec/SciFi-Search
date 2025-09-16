package handlers

import (
	"html/template"
	"net/http"

	sqlc "uki/app/database/sqlc"
)

func RegisterSearchHandlers() {
	// Handler que maneja el acceso al perfil.
	http.HandleFunc("/search", searchHandler)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/search.html"))

	query := r.URL.Query().Get("query")
	var results []sqlc.Work

	if query != "" {
		// Aquí se llama al proveedor de la DB, o la BD en sí.
		//results = provider.Search(query)
	}

	data := struct {
		Query   string
		Results []sqlc.Work
	}{
		Query:   query,
		Results: results,
	}

	tmpl.Execute(w, data)
}
