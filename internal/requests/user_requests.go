package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	minPathLength = 8
)

type BasicAuth struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

func (ba BasicAuth) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(
			&ba.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Invalid email address"),
		),
		validation.Field(
			&ba.Password,
			validation.Required.Error("Password is required"),
			validation.Length(minPathLength, 0).
				Error("Password must be at least 8 characters"),
		),
		// validation.Field(
		// 	&ba.Email,
		// 	validation.Required,
		// 	is.Email,
		// ),
		// validation.Field(
		// 	&ba.Password,
		// 	validation.Required,
		// 	validation.Length(minPathLength, 0),
		// ),
	)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

func (ba LoginRequest) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(
			&ba.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Invalid email address"),
		),
		validation.Field(
			&ba.Password,
			validation.Required.Error("Password is required"),
		),
	)
}

type RegisterRequest struct {
	BasicAuth
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Username string `json:"username" validate:"required" example:"johnDoe"`
}

func (rr RegisterRequest) Validate() error {
	err := rr.BasicAuth.Validate()
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&rr,
		validation.Field(&rr.Name, validation.Required),
		validation.Field(&rr.Username, validation.Required),
	)
}

type OAuthRequest struct {
	Token string `json:"token" validate:"required"`
}

func (oar OAuthRequest) Validate() error {
	return validation.ValidateStruct(&oar,
		validation.Field(&oar.Token, validation.Required),
	)
}

type RefreshRequest struct {
	Token string `json:"token" validate:"required" example:"refresh_token"`
}

type CreateUserRequest struct {
	// Name        string `json:"name" validate:"required" example:"John Doe"`
	Username    string `json:"username" validate:"required" example:"johnDoe"`
	Email       string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password    string `json:"password" validate:"required" example:"123456"`
	IsSuperUser bool   `json:"is_superuser" example:"false"`
	IsActive    bool   `json:"is_active" example:"true"`
	RoleIds     []int  `json:"role_ids" validate:"required" example:"[1]"`
	DeptId      int    `json:"dept_id" validate:"required" example:"1"`
}

func (createUserRequest CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&createUserRequest,
		// validation.Field(&createUserRequest.Name, validation.Required),
		validation.Field(&createUserRequest.Username, validation.Required),
		validation.Field(&createUserRequest.Email, validation.Required),
		validation.Field(&createUserRequest.Password, validation.Required),
		// validation.Field(&createUserRequest.IsSuperUser, validation.Required),
		// validation.Field(&createUserRequest.IsActive, validation.Required),
		// validation.Field(&createUserRequest.RoleIds, validation.Required),
		// validation.Field(&createUserRequest.DeptId, validation.Required),
	)
}

type UpdateUserRequest struct {
	Username    string `json:"username" validate:"required" example:"johnDoe"`
	Email       string `json:"email" validate:"required" example:"john.doe@example.com"`
	IsSuperUser bool   `json:"is_superuser" example:"false"`
	IsActive    bool   `json:"is_active" example:"true"`
	RoleIds     []int  `json:"role_ids" example:"[1]"`
	DeptId      int    `json:"dept_id" example:"1"`
}

func (request UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Username, validation.Required),
		validation.Field(&request.Email, validation.Required),
		validation.Field(&request.IsSuperUser, validation.Required),
		validation.Field(&request.IsActive, validation.Required),
		// validation.Field(&request.RoleIds, validation.Required),
		// validation.Field(&request.DeptId, validation.Required),
	)
}

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"required" example:"johnDoe"`
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
}

func (request UpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Username, validation.Required),
		validation.Field(&request.Email, validation.Required),
	)
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required" example:"123456"`
	NewPassword string `json:"newPassword" validate:"required" example:"333444"`
}

func (request UpdatePasswordRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.OldPassword, validation.Required),
		validation.Field(&request.NewPassword, validation.Required),
	)
}
