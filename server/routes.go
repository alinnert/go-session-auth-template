package server

import (
	"auth-server/handlers"
	"auth-server/middleware"
	"auth-server/values"

	"github.com/dgraph-io/badger"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

func routes(app *chi.Mux, db *badger.DB) {
	// Global middlewares
	app.Use(chiMiddleware.Logger)
	app.Use(chiMiddleware.Recoverer)
	app.Use(middleware.BadgerDB(db))
	app.Use(values.GetSession().LoadAndSave)

	// Routes
	app.Get("/", handlers.GetIndex())
	app.Get("/user", handlers.GetUser())
	app.Get("/public", handlers.GetPublic())
	app.With(middleware.Authenticate()).
		Get("/secret", handlers.GetSecret())

	app.Route("/auth", func(app chi.Router) {
		app.Post("/signup", handlers.SignupHandler())
		app.Post("/signin", handlers.SigninHandler())
		app.With(middleware.Authenticate()).Post("/signout", handlers.SignoutHandler())
	})
}
