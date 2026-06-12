package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BasicRole struct {
	Name string `json:"name" validate:"required" example:"Staff"`
	Desc string `json:"desc" validate:"required" example:"Staff role"`
}

func (bp BasicRole) Validate() error {
	return validation.ValidateStruct(&bp,
		validation.Field(&bp.Name, validation.Required),
		// validation.Field(&bp.Desc, validation.Required),
	)
}

type CreateRoleRequest struct {
	BasicRole
}

type UpdateRoleRequest struct {
	BasicRole
}

type AuthorizeRoleRequest struct {
	ID       int   `json:"id" validate:"required" example:"1"`
	Menu_IDs []int `json:"menu_ids" validate:"required" example:"[1, 2, 3]"`
	Api_IDs  []int `json:"api_ids" validate:"required" example:"[1, 2, 3]"`
}

func (request AuthorizeRoleRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.Menu_IDs,
			validation.Each(validation.Min(1)), // Only validate each item if present
		),
		validation.Field(&request.Api_IDs,
			validation.Each(validation.Min(1)), // Only validate each item if present
		),
	)
}
