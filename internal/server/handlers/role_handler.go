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

type roleService interface {
	GetRoles(ctx context.Context) ([]models.Role, error)
	Create(ctx context.Context, role *models.Role) error
	Update(ctx context.Context, request domain.UpdateRoleRequest) (*models.Role, error)
	Delete(ctx context.Context, request domain.DeleteRoleRequest) error
}

type RoleHandlers struct {
	roleService roleService
}

func NewRoleHandlers(roleService roleService) *RoleHandlers {
	return &RoleHandlers{roleService: roleService}
}

func (h *RoleHandlers) GetRoles(c echo.Context) error {
	roles, err := h.roleService.GetRoles(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all roles: "+err.Error())
	}

	response := responses.NewRoleResponse(roles)
	return responses.Response(c, http.StatusOK, response)
}

func (p *RoleHandlers) CreateRole(c echo.Context) error {
	var createRequest requests.CreateRoleRequest
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

	role := &models.Role{
		Name: createRequest.Name,
		Desc: createRequest.Desc,
	}

	if err := p.roleService.Create(c.Request().Context(), role); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create role: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Role successfully created")
}

func (p *RoleHandlers) UpdateRole(c echo.Context) error {
	idParam := c.Param("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
	}

	var updateRequest requests.UpdateRoleRequest
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

	data := domain.UpdateRoleRequest{
		RoleID: uint(roleID),
		Name:   updateRequest.Name,
		Desc:   updateRequest.Desc,
	}

	if _, err := p.roleService.Update(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update role: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Role successfully updated")
}

func (p *RoleHandlers) DeleteRole(c echo.Context) error {
	idParam := c.Param("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
	}

	data := domain.DeleteRoleRequest{
		RoleID: uint(roleID),
	}

	if err := p.roleService.Delete(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to delete role: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Role successfully deleted")
}
