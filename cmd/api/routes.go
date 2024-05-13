package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	routes := chi.NewRouter()
	routes.Use(app.recoverPanic)
	// routes.Use(app.rateLimitPerClient)
	// routes.Use(app.authenticate)
	routes.Use(app.Metrics)
	routes.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})
	routes.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowedResponse(w, r)
	})
	routes.Get("/v1/healthcheck", app.requireActivatedUser(app.healthcheckHandler))
	routes.Get("/v1/movies/{ID}", app.requireActivatedUser(app.showMovieHandler))
	routes.Get("/v1/movies", app.requireActivatedUser(app.requirePermission("movie:read", app.listMoviesHandler)))
	// routes.Put("/v1/movies/{ID}", app.updateMovieHandler)
	routes.Patch("/v1/movies/{ID}", app.requireActivatedUser(app.updateMovieHandler))
	routes.Post("/v1/movies", app.requirePermission("movie:create", app.requireActivatedUser(app.createMovieHandler)))
	routes.Delete("/v1/movies/{ID}", app.requireActivatedUser(app.deleteMovieHandler))

	//user
	routes.Post("/v1/users", app.registerUserHandler)
	routes.Put("/v1/users/activated", app.activateUserHandler)

	// permissions
	routes.Get("/v1/users/permissions/{ID}", app.getPermissionsByUserID)

	//token
	routes.Post("/v1/tokens/authentication", app.createAutheticationTokenHandler)

	//metrics
	routes.Get("/debug/vars", expvar.Handler().ServeHTTP)
	return routes
}
