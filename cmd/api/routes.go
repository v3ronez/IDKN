package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	routes := chi.NewRouter()
	routes.Get("/v1/healthcheck", app.healthcheckHandler)
	routes.Post("/v1/movies", app.createMovieHandler)
	routes.Get("/v1/movies/{ID}", app.showMovieHandler)

	return routes
}
