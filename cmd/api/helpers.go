package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) readUIDParam(r *http.Request) (int, error) {
	ID, err := strconv.Atoi(chi.URLParam(r, "ID"))
	if err != nil || ID < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return ID, nil
}

func (app *application) writeJSON(v any, w http.ResponseWriter, httpStatus int, headers http.Header) error {
	status, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		app.logger.Println(err)
		return err
	}
	status = append(status, '\n')
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write([]byte(status))
	return nil
}
