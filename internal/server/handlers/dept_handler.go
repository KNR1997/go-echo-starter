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

type departmentService interface {
	GetDepartments(ctx context.Context) ([]models.Department, error)
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
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create department: "+err.Error())
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
