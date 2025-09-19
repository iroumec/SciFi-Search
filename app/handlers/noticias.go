package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/microcosm-cc/bluemonday"

	sqlc "uki/app/database/sqlc"
)

/*
Quizás solo deba haber un POST a noticias.
*/

func manejarNoticias(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		manejarGETNoticias(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func manejarGETNoticias(w http.ResponseWriter, r *http.Request) {

	// Se lee el offset del query.
	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	noticias := obtenerNoticias(r, offset)

	data := map[string]any{
		"Results": noticias,
		"Offset":  offset,
	}

	// Permite ir sumándole 5 all offset. Y usar safeHTML.
	funcs := template.FuncMap{
		"add":      func(a, b int) int { return a + b },
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}

	renderizeTemplate(w, "template/noticias/noticias.html", data, funcs)
}

func obtenerNoticias(r *http.Request, offset int) []sqlc.Noticia {

	noticias, err := queries.ListarNoticias(r.Context(), int32(offset))
	if err != nil {
		noticias = []sqlc.Noticia{}
	}

	return noticias
}

func manejarCargaNoticias(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		manejarGETCargaNoticias(w)
	case http.MethodPost:
		manejarPOSTCargaNoticias(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func manejarGETCargaNoticias(w http.ResponseWriter) {

	renderizeTemplate(w, "template/noticias/cargar-noticia.html", nil, nil)
}

func manejarPOSTCargaNoticias(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	titulo := r.FormValue("título")
	contenidoRaw := r.FormValue("contenido")
	tiempoStr := r.FormValue("tiempo_estimado_lectura")

	// Sanitizar HTML (permitir solo etiquetas seguras para usuarios finales)
	p := bluemonday.UGCPolicy()
	contenido := p.Sanitize(contenidoRaw)

	var tiempo sql.NullTime
	if tiempoStr != "" {
		// Si se maneja como duración, conviene cambiar la columna a int.
		parsed, _ := time.Parse("15:04:05", tiempoStr)
		tiempo = sql.NullTime{Time: parsed, Valid: true}
	}

	_, err := queries.CrearNoticia(r.Context(), sqlc.CrearNoticiaParams{
		Titulo:                titulo,
		Contenido:             contenido,
		PublicadaEn:           sql.NullTime{Time: time.Now(), Valid: true},
		TiempoLecturaEstimado: tiempo,
	})
	if err != nil {
		http.Error(w, "Error guardando noticia", http.StatusInternalServerError)
		return
	}

	renderizeTemplate(w, "template/noticias/carga-exitosa.html", nil, nil)
}
