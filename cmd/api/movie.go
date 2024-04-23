package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/v3ronez/IDKN/internal/data"
	"github.com/v3ronez/IDKN/internal/validator"
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

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Movies.Insert(movie); err != nil {
		app.serverErrorResponse(w, r, err)
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err := app.writeJSON(responseEnvelope{"movie": movie}, w, http.StatusCreated, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.models.Movies.Get(int64(movieID))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	respEnvelope := map[string]any{
		"movie": movie,
	}

	if err := app.writeJSON(respEnvelope, w, http.StatusOK, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return`
	}

}
