package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ej-agas/perfume-db/internal"
	"github.com/ej-agas/perfume-db/postgresql"
)

type createNoteRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	ImageUrl    string `json:"image_url" validate:"omitempty,url"`
	NoteGroupId string `json:"note_group_id" validate:"required"`
}

type updateNoteRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	ImageUrl    string `json:"image_url" validate:"omitempty,url"`
	NoteGroupId string `json:"note_group_id" validate:"omitempty"`
}

func (app *application) listNotes(w http.ResponseWriter, r *http.Request) {
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

	notes, err := app.services.Note.List(id, perPage)
	var newCursor string
	if len(notes) == perPage {
		lastNoteGroup := notes[len(notes)-1]
		newCursor, _ = app.Encrypt([]byte(strconv.Itoa(lastNoteGroup.ID)))
	}

	res := Paginated[internal.Note]{
		Data: notes,
		Next: newCursor,
	}

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.JSONResponse(w, res, 200, nil)
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

	if errors.Is(err, postgresql.ErrNoteGroupNotFound) {
		res := NewValidationErrors()
		res.AddError("note_group_id", "Note Group does not exist.")

		app.JSONResponse(w, res, http.StatusUnprocessableEntity, nil)
		return
	}

	app.logger.Error(err.Error())
	app.ServerError(w)
}

func (app *application) showNoteBySlug(w http.ResponseWriter, r *http.Request) {
	note, err := app.services.Note.FindBySlug(r.PathValue("slug"))
	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, note, http.StatusOK, nil)
}

func (app *application) updateNoteByPublicId(w http.ResponseWriter, r *http.Request) {
	var requestData updateNoteRequest
	note, err := app.services.Note.Find(r.PathValue("publicId"))

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
		note.Name = requestData.Name
		note.Slug = internal.CreateSlug(requestData.Name)
	}

	if requestData.Description != "" {
		note.Description = requestData.Description
	}

	if requestData.ImageUrl != "" {
		note.ImageURL = requestData.ImageUrl
	}

	if requestData.NoteGroupId != "" {
		note.NoteGroupId = requestData.NoteGroupId
	}

	if err := app.services.Note.Save(note); err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.NoContent(w, http.StatusOK)
}
