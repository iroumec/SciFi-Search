package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"slices"

	"uki/app/utils"

	sqlc "uki/app/database/sqlc"

	_ "github.com/lib/pq"
)

const (
	fileDir = "./static"
)

var queries *sqlc.Queries

// registerHandlers registra todos los endpoints
func RegisterHandlers(queryObject *sqlc.Queries) {

	fmt.Println("Comenzando a registrar handlers...")

	queries = queryObject

	// Se crea un manejador (handler) de servidor de archivos.
	fileServer := http.FileServer(http.Dir(fileDir))

	// Se envuelve en un gzip middleware.
	http.Handle("/", utils.GzipMiddleware(fileDir, fileServer))

	http.HandleFunc("/consulta", enqueryHandler)

	http.HandleFunc("/perfil", manejarPerfil)

	http.HandleFunc("/noticias", manejarNoticias)

	setStaticHandler("/carnet-deportivo", "template/card.html")

	// Se registran los handlers correspondientes al manejo de usuarios (registro y login).
	registerUserHandlers()

	fmt.Println("Handlers registrados con éxito.")
}

// ------------------------------------------------------------------------------------------------
// Set Static Handler
// ------------------------------------------------------------------------------------------------

func setStaticHandler(path string, htmlPath string) {

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		tmpl := applyLayout(htmlPath)

		// Se garantiza que el navegador interprete la página como html y con codificación utf-8.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
			http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
		}
	})

}

// ------------------------------------------------------------------------------------------------
// Render Template
// ------------------------------------------------------------------------------------------------

func renderizeTemplate(w http.ResponseWriter, htmlPath string, data map[string]any) {

	tmpl := applyLayout(htmlPath)

	// Se garantiza que el navegador interprete la página como html y con codificación utf-8.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}
}

// ------------------------------------------------------------------------------------------------
// Aplicación de Layout
// ------------------------------------------------------------------------------------------------

func applyLayout(htmlPath string) *template.Template {

	// `template.ParseFiles` abre el archivo y lo convierte en un objeto `*template.Template`.
	// `template.Must` hace que si hay un error al leer la plantilla, el programa panic automáticamente.
	return template.Must(template.ParseFiles(
		"template/layout/layout.html",
		"template/layout/header.html",
		"template/layout/footer.html",
		htmlPath,
	))
}

// ------------------------------------------------------------------------------------------------
// Verificación de campos
// ------------------------------------------------------------------------------------------------

func hayCampoIncompleto(campos ...string) bool {

	return slices.Contains(campos, "")
}

func isThereEmptyField(fields ...string) bool {

	return slices.Contains(fields, "")
}

func atLeastOneFieldIsNotEmpty(fields ...string) bool {

	for _, field := range fields {
		if field != "" {
			return true
		}
	}

	return false
}
