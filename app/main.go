package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"uki/app/handlers"
	"uki/app/utils"

	sqlc "uki/app/database"

	_ "github.com/lib/pq"
)

var port = os.Getenv("APP_PORT")

func main() {

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "db"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "postgres"
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, dbPort, user, password, dbname)

	var err error
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer db.Close()

	queries := sqlc.New(db)

	handlers.RegisterHandlers(queries)

	fmt.Printf("\nServidor escuchando en http://localhost%s\n", port)

	utils.ValidarConstancia()

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

	// Nada de lo que esté acá debajo se ejecuta.
}
