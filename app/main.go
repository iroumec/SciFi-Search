package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"scifi-search/app/bootstrap"
	"scifi-search/app/http/middlewares"
	"scifi-search/app/utils"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {

	appPort := utils.GetEnv("APP_PORT", "8080")

	app, err := bootstrap.Boot()
	if err != nil {
		log.Fatal(err)
	}

	// Contexto de shutdown.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Apagando aplicación...")

		if app.Resources.DB != nil {
			log.Println("Cerrando base de datos...")
			app.Resources.DB.Close()
		}

		os.Exit(0)
	}()

	log.Println("Servidor iniciado")
	log.Printf("Escuchando en http://localhost:%s\n", appPort)

	// El servidor queda a la espera de solicitudes, trabajando en conjunto con un LoggingMiddleware.
	err = http.ListenAndServe(
		":"+appPort,
		supertokens.Middleware(
			middlewares.LoggingMiddleware(http.DefaultServeMux),
		),
	)

	// Nada de lo que esté acá debajo se ejecuta mientras el servidor espere solicitudes.

	if err != nil {
		log.Println("Servidor detenido:", err)
	}
}
