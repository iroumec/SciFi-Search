package handlers

import (
	"html/template"
	"net/http"

	sqlc "uki/app/database/sqlc"
)

func RegisterSearchHandlers() {

	// Handler que maneja la búsqueda.
	http.HandleFunc("/search", searchHandler)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	// `template.ParseFiles` abre el archivo y lo convierte en un objeto `*template.Template`.
	// `template.Must` hace que si hay un error al leer la plantilla, el programa panic automáticamente.
	tmpl := template.Must(template.ParseFiles("template/search.html"))

	// Al usarse el método GET, la url será algo como: /search?query=inception
	// Lo que se hace a continuación es capturar ese valor `inception`.
	query := r.URL.Query().Get("query")

	// Se prepara un slice para guardar resultados.
	var results []sqlc.Work

	if query != "" {
		// Aquí se llama al proveedor de la DB, o la BD en sí.
		//results = provider.Search(query)
	}

	// La plantilla recibe dos variables:
	// .Query (el término que se buscó); y
	// .Results (la lista de contenidos encontrados).
	// Se define entonces acá una estructura con esos campos
	// y se asignan los valores.
	data := struct {
		Query   string
		Results []sqlc.Work
	}{
		Query:   query,
		Results: results,
	}

	// Rellena el html template con los valores de data y lo envía al navegador.
	tmpl.Execute(w, data)
}
