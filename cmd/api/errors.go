package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String()})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	response := responseEnvelope{"error": message}
	err := app.writeJSON(response, w, status, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(status)
	}
}
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	response := responseEnvelope{"error": "rate limit exceeded"}
	err := app.writeJSON(response, w, http.StatusTooManyRequests, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusTooManyRequests)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
