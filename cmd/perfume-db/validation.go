package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
	"unicode"
)

type ValidationErrors struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func CreateResponseFromErrors(err error) *ValidationErrors {
	response := &ValidationErrors{
		Message: "The given data was invalid.",
		Errors:  make(map[string][]string),
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return response
	}

	for _, err := range validationErrors {
		field := fieldToHumanReadable(err.Field())
		jsonTag := fieldToSnakeCase(err.Field())

		switch err.Tag() {
		case "required":
			message := fmt.Sprintf("The %s field is %s.", field, err.Tag())
			response.Errors[jsonTag] = append(response.Errors[jsonTag], message)
		case "gte":
			message := fmt.Sprintf("The %s field should be greater than %s.", field, err.Param())
			response.Errors[jsonTag] = append(response.Errors[jsonTag], message)
		case "lte":
			message := fmt.Sprintf("The %s field should be less than %s.", field, err.Param())
			response.Errors[jsonTag] = append(response.Errors[jsonTag], message)
		case "url":
			message := fmt.Sprintf("The %s field must be a valid URL.", field)
			response.Errors[jsonTag] = append(response.Errors[jsonTag], message)
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
