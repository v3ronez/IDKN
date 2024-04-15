package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"env":     app.config.envMode,
		"status":  "available",
		"version": version,
	}
	if err := app.writeJSON(data, w, http.StatusOK, nil); err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
