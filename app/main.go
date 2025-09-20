package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"tpe/web/app/handlers"

	sqlc "tpe/web/app/database"

	_ "github.com/lib/pq"
)

var port = os.Getenv("APP_PORT")

func main() {

	// Obtención de las variables de ambiente necesarias
	// para conectarse a la base de datos.
	host := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "postgres")

	// Se obtiene la información necesaria para conectarnos a la base de datos a partir de los datos de sesión anteriores.
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, dbPort, user, password, dbname)

	var err error
	// Establecemos conexión con la base de datos.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		// Ocurrió un error...
		log.Fatal("cannot connect to db:", err)
	}
	// Independientemente de lo que ocurra, se cierra la conexión con la base de datos al final.
	defer db.Close()

	// Se obtiene un objeto que nos permita realizar las queries.
	queries := sqlc.New(db)

	// Se registran los handlers.
	handlers.RegisterHandlers(queries)

	// Se informa que el servidor está corriendo.
	fmt.Printf("\nServidor escuchando en http://localhost%s\n", port)

	// El servidor queda a la espera de solicitudes.
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

	// Nada de lo que esté acá debajo se ejecuta.
}

/*
Permite obtener una variable de ambiente o
un valor por defecto, en caso de no hallar la primera.
*/
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
