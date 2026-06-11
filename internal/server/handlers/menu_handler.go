package handlers

import (
	"context"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/requests"
	"go-echo-starter/internal/responses"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type menuService interface {
	GetMenus(ctx context.Context) ([]models.Menu, error)
	GetMenuPaginated(ctx context.Context, pagination domain.Pagination) ([]models.Menu, int64, error)
	Create(ctx context.Context, menu *models.Menu) error
	Update(ctx context.Context, request domain.UpdateMenuRequest) (*models.Menu, error)
	Patch(ctx context.Context, request domain.PatchMenuRequest) (*models.Menu, error)
	Delete(ctx context.Context, request domain.DeleteMenuRequest) error
}

type MenuHandlers struct {
	menuService menuService
}

func NewMenuHandlers(menuService menuService) *MenuHandlers {
	return &MenuHandlers{menuService: menuService}
}

func (h *MenuHandlers) GetMenus(c echo.Context) error {
	menus, err := h.menuService.GetMenus(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all menus: "+err.Error())
	}

	response := responses.NewMenuResponse(menus)
	return responses.Response(c, http.StatusOK, response)
}

func (h *MenuHandlers) GetMenuPaginated(c echo.Context) error {
	page := 1
	pageSize := 10

	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := c.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	pagination := domain.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	menus, total, err := h.menuService.GetMenuPaginated(
		c.Request().Context(),
		pagination,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get menus",
		)
	}

	return responses.Response(c, http.StatusOK, map[string]any{
		"data":     responses.NewMenuTreeResponseOptimized(menus),
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (p *MenuHandlers) CreateMenu(c echo.Context) error {
	var createRequest requests.CreateMenuRequest
	if err := c.Bind(&createRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := createRequest.Validate(); err != nil {
		return responses.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"Validation failed",
			responses.ParseValidationErrors(err),
		)
	}

	menu := &models.Menu{
		Name:      createRequest.Name,
		Remark:    createRequest.Remark,
		MenuType:  createRequest.MenuType,
		Icon:      createRequest.Icon,
		Path:      createRequest.Path,
		Order:     createRequest.Order,
		ParentID:  createRequest.ParentID,
		IsHidden:  createRequest.IsHidden,
		Component: createRequest.Component,
		Keepalive: createRequest.Keepalive,
		Redirect:  createRequest.Redirect,
	}

	if err := p.menuService.Create(c.Request().Context(), menu); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create menu: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Menu successfully created")
}

func (p *MenuHandlers) UpdateMenu(c echo.Context) error {
	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid menu ID")
	}

	var updateRequest requests.UpdateMenuRequest
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

	data := domain.UpdateMenuRequest{
		MenuID:    uint(menuID),
		Name:      updateRequest.Name,
		Remark:    updateRequest.Remark,
		MenuType:  updateRequest.MenuType,
		Icon:      updateRequest.Icon,
		Path:      updateRequest.Path,
		Order:     updateRequest.Order,
		ParentID:  updateRequest.ParentID,
		IsHidden:  updateRequest.IsHidden,
		Component: updateRequest.Component,
		Keepalive: updateRequest.Keepalive,
		Redirect:  updateRequest.Redirect,
	}

	if _, err := p.menuService.Update(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update menu: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Menu successfully updated")
}

func (p *MenuHandlers) PatchMenu(c echo.Context) error {
	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid menu ID")
	}

	var updateRequest requests.PatchMenuRequest
	if err := c.Bind(&updateRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	data := domain.PatchMenuRequest{
		MenuID:    uint(menuID),
		IsHidden:  updateRequest.IsHidden,
		Keepalive: updateRequest.Keepalive,
	}

	if _, err := p.menuService.Patch(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update menu: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Menu successfully updated")
}

func (p *MenuHandlers) DeleteMenu(c echo.Context) error {
	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid menu ID")
	}

	data := domain.DeleteMenuRequest{
		MenuID: uint(menuID),
	}

	if err := p.menuService.Delete(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to delete menu: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Menu successfully deleted")
}
