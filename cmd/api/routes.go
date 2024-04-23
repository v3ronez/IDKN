package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	routes := chi.NewRouter()
	routes.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})
	routes.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowedResponse(w, r)
	})
	routes.Get("/v1/healthcheck", app.healthcheckHandler)
	routes.Get("/v1/movies/{ID}", app.showMovieHandler)
	routes.Post("/v1/movies", app.createMovieHandler)
	routes.Put("/v1/movies/{ID}", app.updateMovieHandler)

	return routes
}
