package main

import (
	"encoding/json"
	"errors"
	"github.com/ej-agas/perfume-db/postgresql"
	"net/http"
)

type createNoteGroupRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	PhotoUrl    string `json:"photo_url" validate:"omitempty,url"`
}

func (app *application) createNoteGroupHandler(w http.ResponseWriter, r *http.Request) {
	var requestData createNoteGroupRequest

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

	noteGroup, err := app.factory.NewNoteGroup(requestData.Name, requestData.Description, requestData.PhotoUrl)
	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	err = app.services.NoteGroup.Save(noteGroup)

	if err == nil {
		app.NoContent(w, http.StatusCreated)
		return
	}

	if errors.Is(err, postgresql.ErrNoteGroupAlreadyExists) {
		app.JSONResponse(w, ResponseMessage{Message: "Note group already exists.", StatusCode: http.StatusUnprocessableEntity}, http.StatusUnprocessableEntity, nil)
		return
	}

	app.logger.Error(err.Error())
	app.ServerError(w)
}

func (app *application) showNoteGroupBySlug(w http.ResponseWriter, r *http.Request) {
	noteGroup, err := app.services.NoteGroup.FindBySlug(r.PathValue("slug"))

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, noteGroup, http.StatusOK, nil)
}
