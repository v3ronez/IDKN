package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	sysInfo := map[string]string{
		"env":     app.config.envMode,
		"version": version,
	}
	respEnv := responseEnvelope{
		"status":      "available",
		"system_info": sysInfo,
	}
	if err := app.writeJSON(respEnv, w, http.StatusOK, nil); err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
