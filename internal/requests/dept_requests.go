package requests

import validation "github.com/go-ozzo/ozzo-validation/v4"

type BasicDept struct {
	Name  string `json:"name" validate:"required" example:"Staff"`
	Desc  string `json:"desc" example:"Staff role"`
	Order string `json:"order" example:"1"`
}

func (bp BasicDept) Validate() error {
	return validation.ValidateStruct(&bp,
		validation.Field(&bp.Name, validation.Required),
	)
}

type CreateDeptRequest struct {
	BasicDept
}

type UpdateDeptRequest struct {
	BasicDept
}
