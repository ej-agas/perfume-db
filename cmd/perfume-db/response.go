package main

import (
	"encoding/json"
	"net/http"

	"github.com/ej-agas/perfume-db/internal"
)

type ResponseMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type Paginated[T internal.Model] struct {
	Data []T    `json:"data"`
	Next string `json:"next"`
}

func (app *application) JSONResponse(w http.ResponseWriter, data any, statusCode int, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(js)

	return nil
}

func (app *application) ServerError(w http.ResponseWriter) {
	res := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "The server encountered an error.",
		Status:  http.StatusInternalServerError,
	}

	app.JSONResponse(w, res, http.StatusInternalServerError, nil)
}

func (app *application) BadRequest(w http.ResponseWriter) {
	res := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "Invalid request body.",
		Status:  http.StatusBadRequest,
	}

	app.JSONResponse(w, res, http.StatusBadRequest, nil)
}

func (app *application) NoContent(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
}
