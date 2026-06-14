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

type userService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, request *requests.CreateUserRequest) error
	Update(ctx context.Context, request domain.UpdateUserRequest) (*models.User, error)
	Patch(ctx context.Context, request domain.PatchUserRequest) (*models.User, error)
	GetUserPaginated(
		ctx context.Context,
		pagination domain.Pagination,
		searchConditions []utils.SearchCondition,
		searchJoin string,
		deptId int,
	) ([]models.User, int64, error)
	Delete(ctx context.Context, request domain.DeleteUserRequest) error
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

func (h *UserHandlers) GetUserPaginated(c echo.Context) error {
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

	deptId, err := strconv.Atoi(c.QueryParam("dept_id"))

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

	users, total, err := h.userService.GetUserPaginated(
		c.Request().Context(),
		pagination,
		searchConditions,
		searchJoin,
		deptId,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get users",
		)
	}

	return responses.Response(c, http.StatusOK, map[string]any{
		"data":     responses.NewUserResponse(users),
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (p *UserHandlers) CreateUser(c echo.Context) error {
	var createRequest requests.CreateUserRequest
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

	if err := p.userService.Create(c.Request().Context(), &createRequest); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return responses.ErrorResponse(c, http.StatusConflict, err.Error())
		}

		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create user")
	}

	return responses.MessageResponse(c, http.StatusCreated, "User successfully created")
}

func (p *UserHandlers) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid department ID")
	}

	var updateRequest requests.UpdateUserRequest
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
		UserID:      uint(userID),
		UserName:    updateRequest.Username,
		Email:       updateRequest.Email,
		IsSuperUser: updateRequest.IsSuperUser,
		IsActive:    updateRequest.IsActive,
		RoleIds:     updateRequest.RoleIds,
		DeptId:      updateRequest.DeptId,
	}

	if _, err := p.userService.Update(c.Request().Context(), data); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return responses.ErrorResponse(c, http.StatusConflict, err.Error())
		}

		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update user: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "User successfully updated")
}

func (p *UserHandlers) PatchUser(c echo.Context) error {
	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	var updateRequest requests.PatchUserRequest
	if err := c.Bind(&updateRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	data := domain.PatchUserRequest{
		UserID:   uint(menuID),
		IsActive: updateRequest.IsActive,
	}

	if _, err := p.userService.Patch(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update user: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "User successfully updated")
}

func (p *UserHandlers) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	data := domain.DeleteUserRequest{
		UserID: uint(userID),
	}

	if err := p.userService.Delete(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to delete user: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "User successfully deleted")
}
