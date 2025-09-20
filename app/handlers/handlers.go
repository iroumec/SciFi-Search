package handlers

import (
	"html/template"
	"net/http"
	"os"
	"slices"

	"tpe/web/app/utils"

	sqlc "tpe/web/app/database"

	_ "github.com/lib/pq"
)

// ------------------------------------------------------------------------------------------------
// Constantes del Paquete
// ------------------------------------------------------------------------------------------------

// Ruta a partir de la cual se servirán los archivos estáticos.
const (
	fileDir = "./static"
)

// ------------------------------------------------------------------------------------------------
// variables Globales al Paquete
// ------------------------------------------------------------------------------------------------

var queries *sqlc.Queries

// ------------------------------------------------------------------------------------------------

// registerHandlers registra todos los endpoints
func RegisterHandlers(queryObject *sqlc.Queries) {

	// Se guarda el objeto de consultas como variable global
	// para poder utilizarlo en todos los handlers que lo requieran.
	queries = queryObject

	// Se registra el hander para los archivos estáticos.
	registrarHandlerStatic()

	// Se registra el handler para el index.html.
	registrarIndexHTML()

	// Se registran los handlers para la página de consultas.
	registrarHandlersConsultas()

	// Se registran los handlers correspondientes al manejo de usuarios (registro y login).
	registrarHandlersUsuarios()

	// Se registran los handlers correspondientes al perfil de usuario.
	registrarHandlersPerfiles()

	// Se registran los handlers correspondientes a las noticias.
	registrarHandlersNoticias()

	// Se registran los handlers correspondientes a servir las páginas de las facultades.
	registerHandlersFacultades()

	// Se registran los handlers correspondientes al área de ayuda/soporte/información.
	registrarHandlersAyuda()
}

// ------------------------------------------------------------------------------------------------

func registrarHandlerStatic() {

	// Se crea un manejador (handler) de servidor de archivos.
	fileServer := http.FileServer(http.Dir(fileDir))

	// Se sirven archivos estáticos en /static/, comprimidos en gzip si el navegador así lo acepta.
	http.Handle("/static/", http.StripPrefix("/static/", utils.GzipMiddleware(fileDir, fileServer)))
}

// ------------------------------------------------------------------------------------------------
// Render Template
// ------------------------------------------------------------------------------------------------

func renderizeTemplate(w http.ResponseWriter, htmlPath string, data map[string]any, funcs template.FuncMap) {

	// Se aplica el layout y las funciones correspondientes a la plantilla.
	tmpl := applyLayout(htmlPath, funcs)

	// Se garantiza que el navegador interprete la página como html y con codificación utf-8.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Se renderiza la plantilla.
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Error al renderizar la plantilla.", http.StatusInternalServerError)
	}
}

// ------------------------------------------------------------------------------------------------
// Aplicación de Layout
// ------------------------------------------------------------------------------------------------

// Esta función aplica el layout a la página HTML.
func applyLayout(htmlPath string, funcs template.FuncMap) *template.Template {

	tmpl := template.New("layout")

	// Si vienen funciones, se aplican.
	if funcs != nil {
		tmpl = tmpl.Funcs(funcs)
	}

	// `template.ParseFiles` abre el archivo y lo convierte en un objeto `*template.Template`.
	// `template.Must` hace que si hay un error al leer la plantilla, el programa panic automáticamente.
	return template.Must(
		tmpl.ParseFiles(
			"template/layout/layout.html",
			"template/layout/header.html",
			"template/layout/footer.html",
			htmlPath,
		),
	)
}

// ------------------------------------------------------------------------------------------------
// Registro de Index
// ------------------------------------------------------------------------------------------------

func registrarIndexHTML() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Se definen las facultades.
		data := struct {
			Facultades []string
		}{
			Facultades: []string{
				"agronomía", "sociales", "humanas", "exactas",
				"ingeniería", "salud", "económicas", "derecho",
				"veterinarias", "arte",
			},
		}

		// Se define una función que establezca los títulos de los botones,
		// la cual se renderizará junto a la página.
		funcs := template.FuncMap{
			"title": func(s string) string {

				// Si no tiene ningún carácter...
				if len(s) == 0 {
					return s
				}

				// Se capitaliza la primera letra de la facultad.
				return string(s[0]-32) + s[1:]
			},
		}

		// Se renderiza la plantilla.
		renderizeTemplate(w, "template/index.html", map[string]any{
			"Facultades": data.Facultades,
		}, funcs)
	})
}

// ------------------------------------------------------------------------------------------------
// Obtener Fotos
// ------------------------------------------------------------------------------------------------

func obtenerFotos(path string) []string {

	// Se obtienen todas las entradas del directorio.
	files, err := os.ReadDir(path)
	if err != nil {
		// No se hallaron fotos.
		return nil
	}

	// TODO: si el directorio tiene un archivo que no sea un
	// directorio o una foto, esto se rompe. Solucionarlo.

	var fotos []string
	for _, file := range files {
		// Si la entrada no es un directorio, la agrega a la lista de fotos.
		if !file.IsDir() {
			fotos = append(fotos, path+file.Name())
		}
	}

	return fotos
}

// ------------------------------------------------------------------------------------------------
// Verificación de campos
// ------------------------------------------------------------------------------------------------

func hayCampoIncompleto(campos ...string) bool {

	return slices.Contains(campos, "")
}
