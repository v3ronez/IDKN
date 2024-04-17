package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type responseEnvelope map[string]any

func (app *application) writeJSON(v responseEnvelope, w http.ResponseWriter, httpStatus int, headers http.Header) error {
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {

	err := json.NewDecoder(r.Body).Decode(dest)
	if err == nil {
		return err
	}

	//json type errors
	var syntexError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	case errors.As(err, &syntexError):
		return fmt.Errorf("body contains badly format JSON (at caracter %d)", syntexError.Offset)

	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON field (at caracter %d)", unmarshalTypeError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly-format JSON")

	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")

	case errors.As(err, &invalidUnmarshalError):
		panic(err)

	default:
		return err
	}
}
