package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	sqlc "uki/app/database/sqlc"

	_ "github.com/lib/pq"
)

const (
	port = ":8080" // Puerto que se escucha.
)

var db *sql.DB
var queries *sqlc.Queries
var ctx context.Context

func main() {

	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer db.Close()

	queries = sqlc.New(db)
	ctx = context.Background()

	registerHandlers()

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}

/*

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	sqlc "T2_E2/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer db.Close()
	queries := sqlc.New(db)
	ctx := context.Background()

	createdUser, err := queries.CreateUser(ctx,
		sqlc.CreateUserParams{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		})

	if err != nil {
		log.Fatal("error creating user:", err)
	}
	fmt.Printf("User created: %+v\n", createdUser)

	user, err := queries.GetUser(ctx, createdUser.ID) // Read One
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	fmt.Printf("Retrieved user: %+v\n", user)

	users, err := queries.ListUsers(ctx) // Read Many
	if err != nil {
		log.Fatalf("failed to list users: %v", err)
	}
	fmt.Printf("All users: %+v\n", users)

	err = queries.UpdateUser(ctx, sqlc.UpdateUserParams{ // Update
		ID:    createdUser.ID,
		Name:  "Johnny Doe",
		Email: "johnny.doe@example.com",
	})

	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}
	fmt.Println("User updated successfully")

	updatedUser, err := queries.GetUser(ctx, createdUser.ID)
	if err != nil {
		log.Fatalf("failed to get updated user: %v", err)
	}

	fmt.Printf("Updated user: %+v\n", updatedUser)
}

*/
