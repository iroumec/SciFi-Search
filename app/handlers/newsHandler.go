package handlers

/*
import (
	"database/sql"
	"log"
	"net/http"

	sqlc "uki/app/database/sqlc"
)

func noticiasHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")

	// Se prepara un slice para guardar resultados.
	var results []sqlc.Work

	if query != "" {
		// Aquí se llama al proveedor de la DB, o la BD en sí.
		//results = provider.Search(query)
	}
	noticias, err := queries.ListarNoticias(r.Context(), 0)
	if err != nil {
		if err == sql.ErrNoRows {
			logInHandleGET(w, "No hay noticias.")
			return
		}
		log.Printf("error getting user: %v", err)
		logInHandleGET(w, "Error interno del servidor.")
		return
	}

	// La plantilla recibe dos variables:
	// .Query (el término que se buscó); y
	// .Results (la lista de contenidos encontrados).
	data := map[string]any{
		"No":      query,
		"Results": results,
	}

	// Se rellena el html template con los valores de data y lo envía al navegador.

	renderizeTemplate(w, "template/news.html", data)
}
*/
