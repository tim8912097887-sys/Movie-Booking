package validation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var (
	validate = validator.New()
	queryDecoder = schema.NewDecoder()
)

func ValidateQueryParams[T any](r *http.Request) (T, error) {
	var queryParams T

    if err := r.ParseForm(); err != nil {
		return queryParams, err
	}

	if err := queryDecoder.Decode(&queryParams, r.Form); err != nil {
		return queryParams, err
	}

	if err := validate.Struct(queryParams); err != nil {
		return queryParams, err
	}

	return queryParams, nil
}

func ValidateRequestBody[T any](r *http.Request) (T, error) {
	var input T

	// Check if the request body is missing entirely
	if r.Body == nil || r.ContentLength == 0 {
        return input, errors.New("request body is empty")
    }

	// Bind the request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return input, err
	}

	if err := validate.Struct(input); err != nil {
        return input, err
    }

    return input, nil
}

func Validate[T any](input T) (T, error) {
	if err := validate.Struct(input); err != nil {
		return input, err
	}
	return input, nil
}