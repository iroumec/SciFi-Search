package handlers

/*
func searchHandler(w http.ResponseWriter, r *http.Request) {

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
	data := map[string]any{
		"Query":   query,
		"Results": results,
	}

	// Se rellena el html template con los valores de data y lo envía al navegador.
	renderizeTemplate(w, "template/search.html", data)
}
*/
