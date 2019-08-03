package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
)

// StartServer is the main entry point of the application.
func StartServer() {
	db := database()
	defer db.Close()

	app := chi.NewRouter()
	routes(app, db)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	addr := ":8080"
	server := &http.Server{
		Addr:         addr,
		Handler:      app,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Graceful shutdown with [ctrl] + [c]
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}

		close(done)
	}()

	logger.Println("Server is ready to handle requests at", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server on addr \"%s\": %v", addr, err)
	}

	<-done
	logger.Println("Server stopped")
}
