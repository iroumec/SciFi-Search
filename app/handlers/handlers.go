package handlers

import (
	"fmt"
	"net/http"

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

	// Se registran los handlers correspondientes al manejo de usuarios (registro y login).
	registerUserHandlers()

	// Se registran los handlers correspondientes al perfil de usuario.
	registerProfileHandlers()

	registerSearchHandlers()

	registerReviewHandlers()

	fmt.Println("Handlers registrados con éxito.")
}

/*

func RegisterHandlers(queryObject *sqlc.Queries) {
	r := chi.NewRouter()

	// Middlewares básicos (logs, CORS, etc)
	r.Use(handlers.LoggingMiddleware)

	// Rutas de usuario
	r.Get("/login", handlers.LoginPageHandler) // GET -> muestra formulario
	r.Post("/login", handlers.LogInHandler)    // POST -> procesa login
	r.Get("/register", handlers.RegisterPageHandler)
	r.Post("/signIn", handlers.SignInHandler)

	// Perfil de usuario
	r.Get("/profile", handlers.ProfileHandler)

	// Ejemplo de rutas dinámicas tipo Letterboxd
	r.Get("/film/{slug}", handlers.FilmHandler)
	r.Get("/film/{slug}/crew", handlers.CrewHandler)

	// Servir archivos estáticos (CSS, JS, imágenes)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", r)
}

*/
