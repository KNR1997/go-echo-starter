package requests

import validation "github.com/go-ozzo/ozzo-validation/v4"

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
