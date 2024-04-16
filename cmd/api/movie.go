package main

import (
	"net/http"
	"time"

	"github.com/v3ronez/IDKN/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieID, err := app.readUIDParam(r)
	_ = movieID
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}
	movie := data.Movie{
		ID:       123,
		Title:    "a cool movie",
		Runtime:  101,
		Genres:   []string{"action", "idk"},
		Version:  1,
		CreateAt: time.Now(),
	}
	respEnvelope := map[string]any{
		"movie": movie,
	}
	if err := app.writeJSON(respEnvelope, w, http.StatusOK, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
