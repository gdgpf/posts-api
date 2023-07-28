package dto

import (
	"encoding/json"
	"io"

	"github.com/gdgpf/posts-api/adapter/validate"
	"github.com/go-playground/validator/v10"
)

type PostRequestBody struct {
	Title       string   `json:"title" validate:"required,lte=100"`
	Description string   `json:"description" validate:"required,lte=500"`
	Latitude    *float32 `json:"latitude"`
	Longitude   *float32 `json:"longitude"`
}

func FromJsonPostRequestBody(body io.ReadCloser, validator *validator.Validate) (*PostRequestBody, error) {
	var postRequestBody PostRequestBody

	if err := json.NewDecoder(body).Decode(&postRequestBody); err != nil {
		return nil, err
	}
	if err := validator.Struct(postRequestBody); err != nil {
		return nil, validate.HandleValidatorFieldError(&postRequestBody, err)
	}

	return &postRequestBody, nil
}
