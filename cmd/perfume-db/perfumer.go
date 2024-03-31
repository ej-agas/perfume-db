package main

import (
	"encoding/json"
	"errors"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/ej-agas/perfume-db/postgresql"
	"net/http"
	"strconv"
	"time"
)

type createPerfumerRequest struct {
	Name        string `json:"name" validate:"required"`
	Nationality string `json:"nationality" validate:"required"`
	ImageUrl    string `json:"image_url" validate:"required,url"`
	BirthDate   string `json:"birth_date" validate:"required,ymd-date-format"`
}

type updatePerfumerRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Nationality string `json:"nationality" validate:"omitempty"`
	ImageUrl    string `json:"image_url" validate:"omitempty,url"`
	BirthDate   string `json:"birth_date" validate:"omitempty,ymd-date-format"`
}

func (app *application) createPerfumerHandler(w http.ResponseWriter, r *http.Request) {
	var req createPerfumerRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.logger.Error(err.Error())
		app.BadRequest(w)
		return
	}

	if err := app.validator.Struct(req); err != nil {
		res := CreateResponseFromErrors(err)
		app.JSONResponse(w, res, http.StatusUnprocessableEntity, nil)
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		app.logger.Error(err.Error())
		res := NewValidationErrors()
		res.AddError("birth_date", "Invalid birth date.")
		app.JSONResponse(w, res, http.StatusUnprocessableEntity, nil)
		return
	}

	perfumer, err := app.factory.NewPerfumer(req.Name, req.Nationality, req.ImageUrl, birthDate)

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	err = app.services.Perfumer.Save(perfumer)

	if err == nil {
		app.NoContent(w, http.StatusCreated)
		return
	}

	if errors.Is(err, postgresql.ErrPerfumerAlreadyExists) {
		app.JSONResponse(w, ResponseMessage{Message: "Perfumer already exists.", StatusCode: http.StatusUnprocessableEntity}, http.StatusUnprocessableEntity, nil)
		return
	}

	app.logger.Error(err.Error())
	app.ServerError(w)
}

func (app *application) updatePerfumerByPublicIdHandler(w http.ResponseWriter, r *http.Request) {
	var req updatePerfumerRequest
	perfumer, err := app.services.Perfumer.Find(r.PathValue("publicId"))

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.logger.Error(err.Error())
		app.BadRequest(w)
		return
	}

	if err := app.validator.Struct(req); err != nil {
		res := CreateResponseFromErrors(err)
		app.JSONResponse(w, res, http.StatusUnprocessableEntity, nil)
		return
	}

	if req.Name != "" {
		perfumer.Name = req.Name
		perfumer.Slug = internal.CreateSlug(req.Name)
	}

	if req.Nationality != "" {
		perfumer.Nationality = req.Nationality
	}

	if req.ImageUrl != "" {
		perfumer.ImageURL = req.ImageUrl
	}

	if req.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			app.logger.Error(err.Error())
			res := NewValidationErrors()
			res.AddError("birth_date", "Invalid birth date.")
			app.JSONResponse(w, res, http.StatusUnprocessableEntity, nil)
			return
		}
		perfumer.BirthDate = birthDate
	}

	if err := app.services.Perfumer.Save(perfumer); err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.NoContent(w, http.StatusOK)
}

func (app *application) listPerfumersHandler(w http.ResponseWriter, r *http.Request) {
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

	perfumers, err := app.services.Perfumer.List(id, perPage)
	var newCursor string
	if len(perfumers) == perPage {
		lastHouse := perfumers[len(perfumers)-1]
		newCursor, _ = app.Encrypt([]byte(strconv.Itoa(lastHouse.ID)))
	}

	res := Paginated[internal.Perfumer]{
		Data: perfumers,
		Next: newCursor,
	}

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.JSONResponse(w, res, 200, nil)
}

func (app *application) showPerfumerBySlugHandler(w http.ResponseWriter, r *http.Request) {
	perfumer, err := app.services.Perfumer.FindBySlug(r.PathValue("slug"))

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, perfumer, http.StatusOK, nil)
}
