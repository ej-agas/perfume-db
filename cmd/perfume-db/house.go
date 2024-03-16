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

type CreateHouseRequest struct {
	Name        string `json:"name" validate:"required"`
	Country     string `json:"country" validate:"required"`
	Description string `json:"description" validate:"required"`
	YearFounded int    `json:"year_founded" validate:"required,gte=1000,lte=9999"`
}

func (app *application) createHouseHandler(w http.ResponseWriter, r *http.Request) {
	var requestData CreateHouseRequest

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

	yearFounded := time.Date(requestData.YearFounded, time.January, 1, 0, 0, 0, 0, time.UTC)
	house := internal.NewHouse(requestData.Name, requestData.Country, requestData.Description, yearFounded)

	err = app.services.House.Save(house)

	if err == nil {
		app.NoContent(w, http.StatusCreated)
		return
	}

	if errors.Is(err, postgresql.ErrHouseAlreadyExists) {
		app.JSONResponse(w, ResponseMessage{Message: "House already exists.", StatusCode: http.StatusUnprocessableEntity}, http.StatusUnprocessableEntity, nil)
		return
	}

	app.logger.Error(err.Error())
	app.ServerError(w)
}

func (app *application) listHouses(w http.ResponseWriter, r *http.Request) {
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

	houses, err := app.services.House.List(id, perPage)
	var newCursor string
	if len(houses) == perPage {
		lastHouse := houses[len(houses)-1]
		newCursor, _ = app.Encrypt([]byte(strconv.Itoa(lastHouse.ID)))
	}

	res := Paginated[internal.House]{
		Data: houses,
		Next: newCursor,
	}

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.JSONResponse(w, res, 200, nil)
}

func (app *application) showHouseById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	house, err := app.services.House.Find(id)

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, house, http.StatusOK, nil)
}

func (app *application) showHouseBySlug(w http.ResponseWriter, r *http.Request) {
	house, err := app.services.House.FindBySlug(r.PathValue("slug"))

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, house, http.StatusOK, nil)
}
