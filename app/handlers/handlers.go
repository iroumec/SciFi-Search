package handlers

import (
	"net/http"
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
	//registrarIndexHTML()

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

	http.Handle("/", fileServer)
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
