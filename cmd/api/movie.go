package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/v3ronez/IDKN/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieID, err := app.readUIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Printf("%+v\n", movieID)
	movie := data.Movie{
		ID:       123,
		Title:    "a cool movie",
		Runtime:  102,
		Genres:   []string{"action", "idk"},
		Version:  1,
		CreateAt: time.Now(),
	}
	if err := app.writeJSON(movie, w, http.StatusOK, nil); err != nil {
		app.logger.Println("sla")
		return
	}

}
