package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"tpe/web/app/handlers"
	"tpe/web/app/meili"
	"tpe/web/app/utils"

	sqlc "tpe/web/app/database"
)

// ------------------------------------------------------------------------------------------------

func main() {

	// Obtención de las variables de ambiente necesarias
	// para conectarse a la base de datos.
	appPort := utils.GetEnv("APP_PORT", ":8080")
	dbHost := utils.GetEnv("DB_HOST", "db")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "postgres")
	dbPassword := utils.GetEnv("DB_PASSWORD", "postgres")
	dbName := utils.GetEnv("DB_NAME", "postgres")

	db := openConnectionToDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)

	// Independientemente de lo que ocurra, se cierra la conexión con la base de datos al final.
	defer db.Close()

	// Se obtiene un objeto que nos permita realizar las queries.
	queries := sqlc.New(db)

	// Se registran los handlers.
	handlers.RegisterHandlers(queries)

	// Se incializan las aplicaciones de terceros.
	initThirdPartyApplication(queries)

	// Se informa que el servidor está corriendo.
	fmt.Printf("\nServidor escuchando en http://localhost:%s\n", appPort)

	// El servidor queda a la espera de solicitudes, trabajando en conjunto con un LoggingMiddleware.
	if err := http.ListenAndServe(":"+appPort, utils.LoggingMiddleware(http.DefaultServeMux)); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

	// Nada de lo que esté acá debajo se ejecuta.
}

// -----------------------------------------------------------------------------------------------

func openConnectionToDatabase(dbHost, dbPort, dbUser, dbPassword, dbName string) *sql.DB {

	// Se obtiene la información necesaria para conectarnos a la base de datos a partir de
	// los datos de sesión definidos anteriormente.
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	// Se extablece conexión con la base de datos.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	return db
}

// -----------------------------------------------------------------------------------------------

func initThirdPartyApplication(queries *sqlc.Queries) {
	meili.Init(queries)
	//supertokens.Init()
}
