package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status:available")
	fmt.Fprintf(w, "env: %s\n", app.config.envMode)
	fmt.Fprintf(w, "versions: %s\n", version)
}
