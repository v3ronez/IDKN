package main

import (
	"fmt"
	"net/http"
)

func (app *application) getPermissionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	permissions, err := app.models.Permissions.GetAllForUser(int64(userID))
	fmt.Println(permissions)
}
