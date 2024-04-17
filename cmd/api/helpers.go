package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (app *application) readIDParam(r *http.Request) (int, error) {
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
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dest)
	if err != nil {
		//json type errors
		var syntexError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntexError):
			return fmt.Errorf("body contains badly format JSON (at caracter %d)", syntexError.Offset)

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON field (at caracter %d)", unmarshalTypeError.Offset)
			//validade unknown fields
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fielName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fielName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
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

	//validade if the body has just one JSON
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
