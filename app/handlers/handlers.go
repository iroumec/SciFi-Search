package handlers

// ------------------------------------------------------------------------------------------------

import (
	"net/http"

	"scifi-search/app/utils"
	"scifi-search/app/views"

	sqlc "scifi-search/app/database"
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
func Init(queryObject *sqlc.Queries) {

	// Se guarda el objeto de consultas para que pueda ser utilizado
	// en todos los handlers que lo requierean.
	queries = queryObject

	// Se registra el hander para los archivos estáticos.
	registrarHandlerStatic()

	// Se registra el handler para el index.html.
	registerIndexHTML()

	// Se registran los handlers correspondientes a la búsqueda.
	registerSearchHandlers()

	// Se registran los handlers correspondientes a la autenticación (EN PROCESO).
	registerAuthenticationHandlers()

	// Se registran los handlers correspondientes a la configuración (EN PROCESO).
	registerSettingsHandlers()

	// Se registran los handlers correspondientes al historial del usuario.
	registerHistoryHandlers()

	// Se registran los handlers correspondientes al manejo de avatars.
	registerAvatarHandlers()

	// Se registran los handlers correspondiente al manejo de tendencias.
	registerTrendsHandler()

	// Se registran los handlers correspondientes al manejo de idiomas.
	registerLanguageHandlers()

	// Se registra un handler que informe el estado del servidor.
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
// Registro de Index
// ------------------------------------------------------------------------------------------------

func registerIndexHTML() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Se crea una instancia de la componente de página.
		component := views.IndexPage(isUserAuthenticated(w, r), utils.GetTranslatorFromRequest(r))

		// Se renderiza la componente.
		component.Render(r.Context(), w)

	})
}

// ------------------------------------------------------------------------------------------------

// Responde con un simple 200 OK. Se utiliza para saber si la aplicación ya se levantó.
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

// ------------------------------------------------------------------------------------------------
