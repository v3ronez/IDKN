package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/v3ronez/IDKN/internal/data"
	"github.com/v3ronez/IDKN/internal/validator"
)

func (app *application) createAutheticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePaswordPlainText(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	t, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthenticaton)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.writeJSON(responseEnvelope{"authentication_token:": t}, w, http.StatusCreated, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
