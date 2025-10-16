package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"tpe/web/app/handlers"
	"tpe/web/app/utils"

	sqlc "tpe/web/app/database"
)

// ------------------------------------------------------------------------------------------------

func main() {

	// Obtención de las variables de ambiente necesarias
	// para conectarse a la base de datos.
	port := utils.GetEnv("APP_PORT", ":8080")
	host := utils.GetEnv("DB_HOST", "db")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	user := utils.GetEnv("DB_USER", "postgres")
	password := utils.GetEnv("DB_PASSWORD", "postgres")
	dbname := utils.GetEnv("DB_NAME", "postgres")

	// Se obtiene la información necesaria para conectarnos a la base de datos a partir de
	// los datos de sesión definidos anteriormente.
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbPort, user, password, dbname)

	var err error
	// Se extablece conexión con la base de datos.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	// Independientemente de lo que ocurra, se cierra la conexión con la base de datos al final.
	defer db.Close()

	// Se obtiene un objeto que nos permita realizar las queries.
	queries := sqlc.New(db)

	// Se registran los handlers.
	handlers.RegisterHandlers(queries)

	// Se informa que el servidor está corriendo.
	fmt.Printf("\nServidor escuchando en http://localhost:%s\n", port)

	// El servidor queda a la espera de solicitudes.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

	// Nada de lo que esté acá debajo se ejecuta.
}

// -----------------------------------------------------------------------------------------------
