package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(chi.URLParam(r, "movieID"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Printf("%+v\n", movieID)
}
