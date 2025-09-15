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

	fmt.Println("Handlers registrados con éxito.")
}
