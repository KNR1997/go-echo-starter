package handlers

import (
	"context"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
}

type UserHandlers struct {
	userService userService
}

func NewUserHandlers(userService userService) *UserHandlers {
	return &UserHandlers{userService: userService}
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Retrieve all users from database
//	@ID				get-users
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{array}		responses.Data
//	@Failure		404		{object}	responses.Error
//	@Router			/users [get]
func (h *UserHandlers) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all users: "+err.Error())
	}

	response := responses.NewUserResponse(users)
	return responses.Response(c, http.StatusOK, response)
}
