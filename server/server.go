package server

import (
	"auth-server/globals"
	"auth-server/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

// StartServer is the main entry point of the application.
func StartServer() {
	// #region Setup global dependencies
	db := getDatabase()
	defer db.Close()
	sessionManager := getSessionManager(db)
	validator := getValidator()
	// #endregion Setup global dependencies

	// #region setup routes and global middleware
	app := chi.NewRouter()
	app.Use(chiMiddleware.Logger)
	app.Use(chiMiddleware.Recoverer)
	app.Use(middleware.ContextData(map[globals.ContextKey]interface{}{
		globals.DBContext:       db,
		globals.SessionContext:  sessionManager,
		globals.ValidatorContext: validator,
	}))
	app.Use(middleware.BadgerDB(db))
	app.Use(sessionManager.LoadAndSave)
	setupRoutes(app)
	// #endregion setup routes and global middleware

	// #region Setup server
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
	// #endregion Setup server

	// #region Graceful shutdown with [ctrl] + [c]
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

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
	// #endregion Graceful shutdown with [ctrl] + [c]
}
