package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/ej-agas/perfume-db/postgresql"
	"net/http"
	"strconv"
	"time"
)

type createHouseRequest struct {
	Name        string `json:"name" validate:"required"`
	Country     string `json:"country" validate:"required"`
	Description string `json:"description" validate:"required"`
	YearFounded int    `json:"year_founded" validate:"required,gte=1000,lte=9999"`
}

type updateHouseRequest struct {
	Name        string `json:"name" validate:"omitempty"`
	Country     string `json:"country" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	YearFounded int    `json:"year_founded" validate:"omitempty,gte=1000,lte=9999"`
}

func (app *application) createHouseHandler(w http.ResponseWriter, r *http.Request) {
	var requestData createHouseRequest

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
	house, err := app.factory.NewHouse(requestData.Name, requestData.Country, requestData.Description, yearFounded)

	if err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

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

func (app *application) showHouseBySlug(w http.ResponseWriter, r *http.Request) {
	house, err := app.services.House.FindBySlug(r.PathValue("slug"))

	if err != nil {
		app.NoContent(w, http.StatusNotFound)
		return
	}

	app.JSONResponse(w, house, http.StatusOK, nil)
}

func (app *application) updateHouseBySlug(w http.ResponseWriter, r *http.Request) {
	var requestData updateHouseRequest
	house, err := app.services.House.FindBySlug(r.PathValue("slug"))
	fmt.Printf("%#v\n", house)
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
		house.Name = requestData.Name
		house.Slug = internal.CreateSlug(requestData.Name)
	}

	if requestData.Description != "" {
		house.Description = requestData.Description
	}

	if requestData.Country != "" {
		house.Country = requestData.Country
	}

	if requestData.YearFounded != 0 {
		house.YearFounded = time.Date(requestData.YearFounded, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	if err := app.services.House.Save(house); err != nil {
		app.logger.Error(err.Error())
		app.ServerError(w)
		return
	}

	app.NoContent(w, http.StatusOK)
}
