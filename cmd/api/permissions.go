package main

import (
	"errors"
	"net/http"

	"github.com/v3ronez/IDKN/internal/data"
)

func (app *application) getPermissionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	permissions, err := app.models.Permissions.GetAllForUser(int64(userID))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	if err = app.writeJSON(responseEnvelope{"permissions": permissions}, w, http.StatusOK, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
