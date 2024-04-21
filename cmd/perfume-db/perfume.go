package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ej-agas/perfume-db/internal"
	"github.com/ej-agas/perfume-db/postgresql"
)

type createPerfumeRequest struct {
	Name             string              `json:"name" validate:"required"`
	Description      string              `json:"description" validate:"required"`
	Concentration    string              `json:"concentration" validate:"required,fragranceConcentration"`
	YearReleased     int                 `json:"year_released" validate:"required,gte=1000,lte=9999"`
	YearDiscontinued int                 `json:"year_discontinued" validate:"omitempty,gte=1000,lte=9999"`
	HouseId          string              `json:"house_id" validate:"required"`
	Perfumers        []string            `json:"perfumers" validate:"required,min=1"`
	Notes            map[string][]string `json:"notes" validate:"required,noteCategory,noteCount=1"`
}

type updatePerfumeRequest struct {
	Name             string              `json:"name" validate:"omitempty"`
	Description      string              `json:"description" validate:"omitempty"`
	Concentration    string              `json:"concentration" validate:"omitempty,fragranceConcentration"`
	YearReleased     int                 `json:"year_released" validate:"omitempty,gte=1000,lte=9999"`
	YearDiscontinued int                 `json:"year_discontinued" validate:"omitempty,gte=1000,lte=9999"`
	HouseId          string              `json:"house_id" validate:"omitempty"`
	Perfumers        []string            `json:"perfumers" validate:"omitempty,min=1"`
	Notes            map[string][]string `json:"notes" validate:"omitempty,noteCategory"`
}

func (app *application) showPerfumeBySlug(w http.ResponseWriter, r *http.Request) {
	house, err := app.services.Perfume.FindBySlug(r.PathValue("slug"))

	if err != nil {
		fmt.Println(err)
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, house, http.StatusOK, nil)
}

func (app *application) createPerfumeHandler(w http.ResponseWriter, r *http.Request) {
	var req createPerfumeRequest
	var yearDiscontinued time.Time

	validationErrors := NewValidationErrors()

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

	house, err := app.services.House.Find(req.HouseId)

	if err != nil {
		validationErrors.AddError("house_id", "House not found.")
		app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
		return
	}

	perfumers, err := app.services.Perfumer.FindMany(req.Perfumers...)
	if err != nil {
		switch {
		case errors.Is(err, postgresql.ErrPerfumerNotFound):
			validationErrors.AddError("perfumers", err.Error())
			app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
			return
		default:
			app.logger.Error(err.Error())
			app.JSONResponse(w, err.Error(), http.StatusInternalServerError, nil)
			return
		}
	}

	yearReleased := time.Date(req.YearReleased, time.January, 1, 0, 0, 0, 0, time.UTC)
	if req.YearDiscontinued != 0 {
		yearDiscontinued = time.Date(req.YearDiscontinued, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	concentration, _ := internal.ConcentrationFromString(req.Concentration)
	notes := make(map[internal.NoteCategory][]*internal.Note, len(req.Notes))

	for key, publicIds := range req.Notes {
		category, err := internal.NoteCategoryFromString(key)

		if err != nil {
			validationErrors.AddError("notes", "Invalid note category.")
			app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
			return
		}

		notesResult, err := app.services.Note.FindMany(publicIds)
		if err != nil {
			validationErrors.AddError("notes", err.Error())
			app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
			return
		}

		notes[category] = notesResult
	}

	perfume, err := app.factory.NewPerfume(
		internal.WithName(req.Name),
		internal.WithDescription(req.Description),
		internal.WithConcentration(concentration),
		internal.WithYearReleased(yearReleased),
		internal.WithYearDiscontinued(yearDiscontinued),
		internal.WithHouse(house),
		internal.WithPerfumers(perfumers...),
		internal.WithNotes(notes),
	)

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	err = app.services.Perfume.Save(perfume)

	if err != nil {
		switch {
		case errors.Is(err, postgresql.ErrPerfumeAlreadyExists):
			app.JSONResponse(w, ResponseMessage{Message: "Perfume already exists", StatusCode: 422}, 422, nil)
		case errors.Is(err, postgresql.ErrHouseNotFound):
			validationErrors.AddError("house_id", "House not found.")
			app.JSONResponse(w, validationErrors, 422, nil)
		default:
			app.logger.Error(err.Error())
			app.ServerError(w)
		}
		return
	}

	app.JSONResponse(w, perfume, 200, nil)
}

func (app *application) updatePerfumeHandler(w http.ResponseWriter, r *http.Request) {
	var req updatePerfumeRequest
	//var yearDiscontinued time.Time

	perfume, err := app.services.Perfume.Find(r.PathValue("publicId"))

	if err != nil {
		fmt.Println(err)
		app.NoContent(w, 404)
		return
	}

	validationErrors := NewValidationErrors()

	err = json.NewDecoder(r.Body).Decode(&req)
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

	if req.Name != "" {
		perfume.Name = req.Name
	}

	if req.Description != "" {
		perfume.Description = req.Description
	}

	if req.Concentration != "" {
		concentration, _ := internal.ConcentrationFromString(req.Concentration)
		perfume.Concentration = concentration
	}

	if req.YearReleased != 0 {
		perfume.YearReleased = time.Date(req.YearReleased, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	if req.YearDiscontinued != 0 {
		perfume.YearDiscontinued = time.Date(req.YearDiscontinued, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	if req.HouseId != "" {
		house, err := app.services.House.Find(req.HouseId)

		if err != nil {
			validationErrors.AddError("house_id", "House not found.")
			app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
			return
		}

		perfume.House = house
	}

	if len(req.Perfumers) != 0 {
		perfumers, err := app.services.Perfumer.FindMany(req.Perfumers...)
		if err != nil {
			switch {
			case errors.Is(err, postgresql.ErrPerfumerNotFound):
				validationErrors.AddError("perfumers", err.Error())
				app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
				return
			default:
				app.logger.Error(err.Error())
				app.JSONResponse(w, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		perfume.Perfumers = perfumers
	}

	if req.Notes != nil {
		notes := make(map[internal.NoteCategory][]*internal.Note, len(req.Notes))

		for key, publicIds := range req.Notes {
			category, err := internal.NoteCategoryFromString(key)

			if err != nil {
				validationErrors.AddError("notes", "Invalid note category.")
				app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
				return
			}

			notesResult, err := app.services.Note.FindMany(publicIds)
			if err != nil {
				validationErrors.AddError("notes", err.Error())
				app.JSONResponse(w, validationErrors, http.StatusUnprocessableEntity, nil)
				return
			}

			notes[category] = notesResult
		}

		perfume.Notes = notes
	}

	if err := app.services.Perfume.Save(perfume); err != nil {
		app.logger.Error(err.Error())
		app.JSONResponse(w, err.Error(), 500, nil)
		return
	}

	app.JSONResponse(w, perfume, 200, nil)
}
