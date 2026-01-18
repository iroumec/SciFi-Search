package main

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------------
// Entry Point
// ------------------------------------------------------------------------------------------------

func main() {

	appPort := utils.GetEnv("APP_PORT", "8080")

	// Boot.
	app, err := bootstrap.Boot()
	if err != nil {
		log.Fatal(err)
	}

	// Shutdown context.
	// It ensures the shutdown of everything in case of: Ctrl + C, Docker Down, etc.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down...")

		if app.Resources.DB != nil {
			log.Println("Closing database...")
			app.Resources.DB.Close()
		}

		os.Exit(0)
	}()

	log.Println("Server running")
	log.Printf("Listening on http://localhost:%s\n", appPort)

	// The server waits for requests.
	err = http.ListenAndServe(
		":"+appPort,
		supertokens.Middleware(
			middlewares.LoggingMiddleware(http.DefaultServeMux),
		),
	)

	// The following lines are only executed when the server is no longer awaiting for requests.
	if err != nil {
		log.Println("Server stopped:", err)
	}
}

// ------------------------------------------------------------------------------------------------
