package handlers

import (
	"context"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type baseService interface {
	GetUserMenus(ctx context.Context, userId uint) ([]models.Menu, error)
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
