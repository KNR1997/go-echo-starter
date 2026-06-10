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

type departmentService interface {
	GetDepartments(ctx context.Context) ([]models.Department, error)
	GetDepartmentPaginated(
		ctx context.Context,
		pagination domain.Pagination,
		searchConditions []utils.SearchCondition,
		searchJoin string,
	) ([]models.Department, int64, error)
	Create(ctx context.Context, department *models.Department) error
	Update(ctx context.Context, request domain.UpdateDepartmentRequest) (*models.Department, error)
	Delete(ctx context.Context, request domain.DeleteDepartmentRequest) error
}

type DepartmentHandlers struct {
	departmentService departmentService
}

func NewDepartmentHandlers(departmentService departmentService) *DepartmentHandlers {
	return &DepartmentHandlers{departmentService: departmentService}
}

func (h *DepartmentHandlers) GetDepartments(c echo.Context) error {
	departments, err := h.departmentService.GetDepartments(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all departments: "+err.Error())
	}

	response := responses.NewDeptResponse(departments)
	return responses.Response(c, http.StatusOK, response)
}

// handlers/department_handler.go
func (h *DepartmentHandlers) GetDepartmentPaginated(c echo.Context) error {
	// Parse pagination
	page := 1
	pageSize := 5

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

	departments, total, err := h.departmentService.GetDepartmentPaginated(
		c.Request().Context(),
		pagination,
		searchConditions,
		searchJoin,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get departments",
		)
	}

	return responses.Response(c, http.StatusOK, map[string]any{
		"data":     responses.NewDeptResponse(departments),
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (p *DepartmentHandlers) CreateDepartment(c echo.Context) error {
	var createRequest requests.CreateDeptRequest
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

	department := &models.Department{
		Name: createRequest.Name,
		Desc: createRequest.Desc,
	}

	if err := p.departmentService.Create(c.Request().Context(), department); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return responses.ErrorResponse(c, http.StatusConflict, err.Error())
		}

		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create department")
	}

	return responses.MessageResponse(c, http.StatusCreated, "Department successfully created")
}

func (p *DepartmentHandlers) UpdateDepartment(c echo.Context) error {
	idParam := c.Param("id")
	departmentID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid department ID")
	}

	var updateRequest requests.UpdateDeptRequest
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

	data := domain.UpdateDepartmentRequest{
		DeptID: uint(departmentID),
		Name:   updateRequest.Name,
		Desc:   updateRequest.Desc,
	}

	if _, err := p.departmentService.Update(c.Request().Context(), data); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return responses.ErrorResponse(c, http.StatusConflict, err.Error())
		}

		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update department: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Department successfully updated")
}

func (p *DepartmentHandlers) DeleteDepartment(c echo.Context) error {
	idParam := c.Param("id")
	departmentID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid department ID")
	}

	data := domain.DeleteDepartmentRequest{
		DeptID: uint(departmentID),
	}

	if err := p.departmentService.Delete(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to delete department: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Department successfully deleted")
}
