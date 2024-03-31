package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
	"unicode"
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
