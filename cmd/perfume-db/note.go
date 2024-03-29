package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ej-agas/perfume-db/postgresql"
	"net/http"
)

type createNoteRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	ImageUrl    string `json:"image_url" validate:"omitempty,url"`
	NoteGroupId string `json:"note_group_id" validate:"required"`
}

func (app *application) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var requestData createNoteRequest

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		app.logger.Error(err.Error())
		app.BadRequest(w)
		return
	}

	if err := app.validator.Struct(requestData); err != nil {
		res := CreateResponseFromErrors(err)
		app.JSONResponse(w, res, http.StatusUnprocessableEntity, nil)
		return
	}

	note, err := app.factory.NewNote(
		requestData.Name,
		requestData.Description,
		requestData.ImageUrl,
		requestData.NoteGroupId,
	)

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	err = app.services.Note.Save(note)

	if err == nil {
		app.NoContent(w, http.StatusCreated)
		return
	}

	if errors.Is(err, postgresql.ErrNoteAlreadyExists) {
		app.JSONResponse(w, ResponseMessage{Message: "Note already exists.", StatusCode: http.StatusUnprocessableEntity}, http.StatusUnprocessableEntity, nil)
		return
	}

	app.logger.Error(err.Error())
	app.ServerError(w)
}

func (app *application) showNoteBySlug(w http.ResponseWriter, r *http.Request) {
	note, err := app.services.Note.FindBySlug(r.PathValue("slug"))
	fmt.Printf("%#v\n", note)
	fmt.Printf("%#v\n", err)
	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, note, http.StatusOK, nil)
}
