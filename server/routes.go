package server

import (
	"auth-server/handlers"
	"auth-server/middleware"

	"github.com/go-chi/chi"
)

func setupRoutes(app *chi.Mux) {
	// Routes
	app.Get("/", handlers.GetIndex())
	app.Get("/db", handlers.GetDb())
	app.Get("/user", handlers.GetUser())
	app.Get("/profile", handlers.GetProfile())
	app.Get("/public", handlers.GetPublic())
	app.With(middleware.Authenticate()).
		Get("/secret", handlers.GetSecret())

	app.Route("/auth", func(app chi.Router) {
		app.Post("/signup", handlers.SignupHandler())
		app.Post("/signin", handlers.SigninHandler())
		app.With(middleware.Authenticate()).Post("/signout", handlers.SignoutHandler())
	})
}
