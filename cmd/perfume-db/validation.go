package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ej-agas/perfume-db/internal"
	"github.com/go-playground/validator/v10"
)

type ValidationErrors struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{Message: "The given data was invalid.", Errors: make(map[string][]string)}
}

func (validationErrors *ValidationErrors) AddError(field, message string) {
	validationErrors.Errors[field] = append(validationErrors.Errors[field], message)
}

func CreateResponseFromErrors(err error) *ValidationErrors {
	response := NewValidationErrors()

	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		return response
	}

	for _, err := range validationErrors {
		field := fieldToHumanReadable(err.Field())
		jsonTag := fieldToSnakeCase(err.Field())

		switch err.Tag() {
		case "required":
			message := fmt.Sprintf("The %s field is %s.", field, err.Tag())
			response.AddError(jsonTag, message)
		case "gte":
			message := fmt.Sprintf("The %s field should be greater than %s.", field, err.Param())
			response.AddError(jsonTag, message)
		case "lte":
			message := fmt.Sprintf("The %s field should be less than %s.", field, err.Param())
			response.AddError(jsonTag, message)
		case "url":
			message := fmt.Sprintf("The %s field must be a valid URL.", field)
			response.AddError(jsonTag, message)
		case "ymd-date-format":
			message := fmt.Sprintf("The %s field must be a valid date format 'YYYY-MM-DD'.", field)
			response.AddError(jsonTag, message)
		case "min":
			message := fmt.Sprintf("The %s field must have a minimum count of %s", field, err.Param())
			response.AddError(jsonTag, message)
		case "fragranceConcentration":
			message := fmt.Sprintf("The selected %s is invalid", field)
			response.AddError(jsonTag, message)
		case "noteCategory":
			message := fmt.Sprintf("The %s field contains invalid note category.", field)
			response.AddError(jsonTag, message)
		case "noteCount":
			message := fmt.Sprintf("The %s field must have a minimum count of %s.", field, err.Param())
			response.AddError(jsonTag, message)
		}
	}

	return response
}

func fieldToHumanReadable(field string) string {
	var humanField string
	for i, char := range field {
		if i > 0 && unicode.IsUpper(char) {
			humanField += " "
		}
		humanField += string(char)
	}
	return strings.ToLower(humanField)
}

func fieldToSnakeCase(field string) string {
	var snakeCasedField string

	for i, char := range field {
		if i > 0 && unicode.IsUpper(char) {
			snakeCasedField += "_"
		}

		snakeCasedField += string(char)
	}

	return strings.ToLower(snakeCasedField)
}

type DateValidator struct{}

func (dv DateValidator) Validate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)

	return err == nil
}

type FragranceConcentrationValidator struct{}

func (validator FragranceConcentrationValidator) Validate(fl validator.FieldLevel) bool {
	str := fl.Field().String()

	_, err := internal.ConcentrationFromString(str)

	return err == nil
}

type NoteCategoriesValidator struct{}

func (validator NoteCategoriesValidator) Validate(fl validator.FieldLevel) bool {
	notes := fl.Field().Interface().(map[string][]string)

	for keys := range notes {
		if _, ok := internal.NoteCategoryMap[keys]; !ok {
			return false
		}
	}

	return true
}

type NoteCountValidator struct{}

func (validator NoteCountValidator) Validate(fl validator.FieldLevel) bool {
	notes := fl.Field().Interface().(map[string][]string)

	minCount, err := strconv.Atoi(fl.Param())

	if err != nil {
		panic(err)
	}

	for _, values := range notes {
		if len(values) < minCount {
			return false
		}
	}

	return true
}
