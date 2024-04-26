package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	routes := chi.NewRouter()
	routes.Use(app.recoverPanic)
	// routes.Use(app.rateLimit)
	routes.Use(app.rateLimitPerClient)
	routes.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})
	routes.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowedResponse(w, r)
	})
	routes.Get("/v1/healthcheck", app.healthcheckHandler)
	routes.Get("/v1/movies/{ID}", app.showMovieHandler)
	routes.Get("/v1/movies", app.listMoviesHandler)
	// routes.Put("/v1/movies/{ID}", app.updateMovieHandler)
	routes.Patch("/v1/movies/{ID}", app.updateMovieHandler)
	routes.Post("/v1/movies", app.createMovieHandler)
	routes.Delete("/v1/movies/{ID}", app.deleteMovieHandler)

	return routes
}
