package handlers

import (
	"context"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/requests"
	"go-echo-starter/internal/responses"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type baseService interface {
	GetUserMenus(ctx context.Context, userId uint) ([]models.Menu, error)
	GetMeDetails(ctx context.Context, userId uint) (*models.User, error)
	ProfileUpdate(ctx context.Context, request domain.UpdateUserRequest) (*models.User, error)
	PasswordUpdate(ctx context.Context, request domain.UpdatePasswordRequest) (*models.User, error)
}

type BaseHandlers struct {
	baseService baseService
}

func NewBaseHandlers(baseService baseService) *BaseHandlers {
	return &BaseHandlers{baseService: baseService}
}

func (h *BaseHandlers) GetUserMenu(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	menus, err := h.baseService.GetUserMenus(c.Request().Context(), authClaims.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all menus: "+err.Error())
	}

	response := responses.NewMenuTreeResponse(menus)
	return responses.Response(c, http.StatusOK, response)
}

func (h *BaseHandlers) GetMeDetails(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.baseService.GetMeDetails(c.Request().Context(), authClaims.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all menus: "+err.Error())
	}

	response := responses.NewSingleUserResponse(user)
	return responses.Response(c, http.StatusOK, response)
}

func (p *BaseHandlers) ProfileUpdate(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var updateRequest requests.UpdateProfileRequest
	if err := c.Bind(&updateRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := updateRequest.Validate(); err != nil {
		return responses.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"Validation failed",
			responses.ParseValidationErrors(err),
		)
	}

	data := domain.UpdateUserRequest{
		UserID:   authClaims.ID,
		UserName: updateRequest.Username,
		Email:    updateRequest.Email,
	}

	if _, err := p.baseService.ProfileUpdate(c.Request().Context(), data); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return responses.ErrorResponse(c, http.StatusConflict, err.Error())
		}

		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update user: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "User Profile successfully updated")
}

func (p *BaseHandlers) PasswordUpdate(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var updateRequest requests.UpdatePasswordRequest
	if err := c.Bind(&updateRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	data := domain.UpdatePasswordRequest{
		UserID:      authClaims.ID,
		OldPassword: updateRequest.OldPassword,
		NewPassword: updateRequest.NewPassword,
	}

	if err := updateRequest.Validate(); err != nil {
		return responses.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"Validation failed",
			responses.ParseValidationErrors(err),
		)
	}

	if _, err := p.baseService.PasswordUpdate(c.Request().Context(), data); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return responses.ErrorResponse(c, http.StatusConflict, err.Error())
		}

		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update user password: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "User password successfully updated")
}
