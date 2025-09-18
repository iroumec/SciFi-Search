package handlers

import (
	"database/sql"
	"log"
	"net/http"

	sqlc "uki/app/database/sqlc"
)

func manejarNoticias(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		manejarGETNoticias(w, 0)
	case http.MethodPost:
		manejarPOSTNoticias(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func manejarGETNoticias(w http.ResponseWriter, offset int) {

	noticias := obtenerNoticias(offset)

	data := map[string]any{
		"Results": {
			"Título": 
		}
	}

}

func manejarPostNoticias(w http.ResponseWriter, r *http.Request) {

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

func obtenerNoticias(offset int) []sqlc.Noticia {

	noticias, err := queries.ListarNoticias(r.Context(), int32(offset))
	if err != nil {
		if err == sql.ErrNoRows {
			logInHandleGET(w, "No hay noticias.")
			return
		}
		log.Printf("error getting user: %v", err)
		logInHandleGET(w, "Error interno del servidor.")
		return nil
	}

	return noticias
}
