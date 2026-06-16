package handlers

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/requests"
	"go-echo-starter/internal/responses"
	"go-echo-starter/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type roleService interface {
	GetRoles(ctx context.Context) ([]models.Role, error)
	GetRolePaginated(
		ctx context.Context,
		pagination domain.Pagination,
		searchConditions []utils.SearchCondition,
		searchJoin string,
	) ([]models.Role, int64, error)
	Create(ctx context.Context, role *models.Role) error
	Update(ctx context.Context, request domain.UpdateRoleRequest) (*models.Role, error)
	Delete(ctx context.Context, request domain.DeleteRoleRequest) error
	Authorize(ctx context.Context, request domain.AuthorizeRoleRequest) (*models.Role, error)
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

func (h *RoleHandlers) GetRolePaginated(c echo.Context) error {
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

	// Parse search parameters
	searchConditions, err := utils.ParseSearchQuery(c.QueryParam("search"))
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid search format: %v", err),
		)
	}

	// Parse search join (default to "and")
	searchJoin := c.QueryParam("searchJoin")
	if searchJoin == "" {
		searchJoin = "and"
	}

	// Convert to lowercase for consistency
	searchJoin = strings.ToLower(searchJoin)
	if searchJoin != "and" && searchJoin != "or" {
		return responses.ErrorResponse(
			c,
			http.StatusBadRequest,
			"searchJoin must be 'and' or 'or'",
		)
	}

	pagination := domain.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	roles, total, err := h.roleService.GetRolePaginated(
		c.Request().Context(),
		pagination,
		searchConditions,
		searchJoin,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get roles",
		)
	}

	return responses.Response(c, http.StatusOK, map[string]any{
		"data":     responses.NewRoleResponse(roles),
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
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
		Name:        createRequest.Name,
		Description: createRequest.Description,
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
		RoleID:      uint(roleID),
		Name:        updateRequest.Name,
		Description: updateRequest.Description,
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

func (p *RoleHandlers) AuthorizeRole(c echo.Context) error {
	idParam := c.Param("id")
	roleID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
	}

	var authorizeRequest requests.AuthorizeRoleRequest
	if err := c.Bind(&authorizeRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := authorizeRequest.Validate(); err != nil {
		return responses.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"Validation failed",
			responses.ParseValidationErrors(err),
		)
	}

	data := domain.AuthorizeRoleRequest{
		RoleID:  uint(roleID),
		MenuIDs: authorizeRequest.Menu_IDs,
		ApiIDs:  authorizeRequest.Api_IDs,
	}

	if _, err := p.roleService.Authorize(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to authorize role: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Role successfully authorized")
}
