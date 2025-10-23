package handlers

import (
	"html/template"
	"net/http"
	"os"
	"slices"

	"tpe/web/app/utils"
	"tpe/web/app/views"

	sqlc "tpe/web/app/database"

	"github.com/nats-io/nats.go"
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
var nat *nats.Conn

// ------------------------------------------------------------------------------------------------

// registerHandlers registra todos los endpoints
func RegisterHandlers(queryObject *sqlc.Queries, natObject *nats.Conn) {

	nat = natObject

	// Se guarda el objeto de consultas como variable global
	// para poder utilizarlo en todos los handlers que lo requieran.
	queries = queryObject

	// Se registra el hander para los archivos estáticos.
	registrarHandlerStatic()

	// Se registra el handler para el index.html.
	registrarIndexHTML()

	// Se registran los handlers correspondientes al manejo de usuarios (registro y login).
	registrarHandlersUsuarios()

	// Se registran los handlers correspondientes al perfil de usuario.
	registrarHandlersPerfiles()

	// Se registran los handlers correspondientes al área de ayuda/soporte/información.
	//registrarHandlersAyuda()

	registrarLogIn()

	http.HandleFunc("/health", healthCheckHandler)
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

		// Se crea una instancia de la componente de página.
		component := views.IndexPage()

		// Se renderiza la componente.
		component.Render(r.Context(), w)
	})
}

// ------------------------------------------------------------------------------------------------
// Registro de LogIn
// ------------------------------------------------------------------------------------------------

func registrarLogIn() {

	http.HandleFunc("/log-in", func(w http.ResponseWriter, r *http.Request) {

		component := views.LoginPage("string")

		component.Render(r.Context(), w)
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

// healthCheckHandler responde con un simple 200 OK.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	// Solo responde a peticiones GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Se establece el código de estado 200 OK.
	// A esto lo busca `curl -f` cuando se levanta el servidor.
	w.WriteHeader(http.StatusOK)

	// Cuerpo simple para saber que funciona si se abre desde un navegador.
	w.Write([]byte("Servidor OK"))
}
