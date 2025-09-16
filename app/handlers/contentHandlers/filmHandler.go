package contenthandlers

import (
	"html/template"
	"net/http"
)

func FilmHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	// Aquí buscás la película en tu DB o cache
	film, err := queries.GetFilmBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/film.html"))
	tmpl.Execute(w, film)
}
