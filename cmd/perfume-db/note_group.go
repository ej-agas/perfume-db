package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ej-agas/perfume-db/internal"
	"github.com/ej-agas/perfume-db/postgresql"
)

type createNoteGroupRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	ImageUrl    string `json:"image_url" validate:"omitempty,url"`
}

type updateNoteGroupRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	ImageUrl    string `json:"image_url" validate:"omitempty,url"`
}

func (app *application) listNoteGroups(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	var id = 0

	if cursor != "" {
		decrypted, err := app.Decrypt(cursor)
		if err != nil {
			id = 0
		}

		convertedID, err := strconv.Atoi(string(decrypted))
		if err != nil {
			id = 0
		}
		id = convertedID
	}

	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil || perPage < 0 || perPage > 100 {
		perPage = 25
	}

	noteGroups, err := app.services.NoteGroup.List(id, perPage)
	var newCursor string
	if len(noteGroups) == perPage {
		lastNoteGroup := noteGroups[len(noteGroups)-1]
		newCursor, _ = app.Encrypt([]byte(strconv.Itoa(lastNoteGroup.ID)))
	}

	res := Paginated[internal.NoteGroup]{
		Data: noteGroups,
		Next: newCursor,
	}

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.JSONResponse(w, res, 200, nil)
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

	noteGroup, err := app.factory.NewNoteGroup(requestData.Name, requestData.Description, requestData.ImageUrl)
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

func (app *application) updateNoteGroupByPublicId(w http.ResponseWriter, r *http.Request) {
	var requestData updateNoteGroupRequest

	noteGroup, err := app.services.NoteGroup.Find(r.PathValue("publicId"))

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		app.logger.Error(err.Error())
		app.BadRequest(w)
		return
	}

	if err := app.validator.Struct(requestData); err != nil {
		res := CreateResponseFromErrors(err)
		app.JSONResponse(w, res, http.StatusUnprocessableEntity, nil)
		return
	}

	if requestData.Name != "" {
		noteGroup.Name = requestData.Name
		noteGroup.Slug = internal.CreateSlug(requestData.Name)
	}

	if requestData.Description != "" {
		noteGroup.Description = requestData.Description
	}

	if requestData.ImageUrl != "" {
		noteGroup.ImageURL = requestData.ImageUrl
	}

	if err := app.services.NoteGroup.Save(noteGroup); err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.NoContent(w, http.StatusOK)
}
