package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BasicApi struct {
	Path    string `json:"path" validate:"required" example:"Staff"`
	Method  string `json:"method" validate:"required" example:"GET"`
	Summary string `json:"summary" validate:"required" example:"Staff role"`
	Tags    string `json:"tags" validate:"required" example:"Staff role"`
}

func (bp BasicApi) Validate() error {
	return validation.ValidateStruct(&bp,
		validation.Field(&bp.Path, validation.Required),
		validation.Field(&bp.Method,
			validation.Required,
			validation.In("GET", "POST", "PUT", "DELETE", "PATCH").Error("method must be one of: GET, POST, PUT, DELETE, PATCH"),
		),
		validation.Field(&bp.Summary, validation.Required),
		validation.Field(&bp.Tags, validation.Required),
	)
}

type CreateApiRequest struct {
	BasicApi
}

type UpdateApiRequest struct {
	BasicApi
}
