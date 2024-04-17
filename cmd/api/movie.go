package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/v3ronez/IDKN/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// var movie data.Movie
	var input struct {
		Title   string        `json:"title"`
		Year    int32         `json:"year"`
		Runtime data.Runtinme `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Printf("%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieID, err := app.readIDParam(r)
	_ = movieID
	if err != nil {
		app.notFoundResponse(w, r)
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
