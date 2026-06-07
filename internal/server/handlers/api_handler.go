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

type apiService interface {
	GetApis(ctx context.Context) ([]models.Api, error)
	GetApiPaginated(ctx context.Context, pagination domain.Pagination) ([]models.Api, int64, error)
	Create(ctx context.Context, api *models.Api) error
	Update(ctx context.Context, request domain.UpdateApiRequest) (*models.Api, error)
	Delete(ctx context.Context, request domain.DeleteApiRequest) error
}

type ApiHandlers struct {
	apiService apiService
}

func NewApiHandlers(apiService apiService) *ApiHandlers {
	return &ApiHandlers{apiService: apiService}
}

func (h *ApiHandlers) GetApis(c echo.Context) error {
	apis, err := h.apiService.GetApis(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all apis: "+err.Error())
	}

	response := responses.NewApiResponse(apis)
	return responses.Response(c, http.StatusOK, response)
}

func (h *ApiHandlers) GetApiPaginated(c echo.Context) error {
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

	apis, total, err := h.apiService.GetApiPaginated(
		c.Request().Context(),
		pagination,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get apis",
		)
	}

	return responses.Response(c, http.StatusOK, map[string]any{
		"data":     responses.NewApiResponse(apis),
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func (p *ApiHandlers) CreateApi(c echo.Context) error {
	var createRequest requests.CreateApiRequest
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

	api := &models.Api{
		Path:    createRequest.Path,
		Method:  createRequest.Method,
		Summary: createRequest.Summary,
		Tags:    createRequest.Tags,
	}

	if err := p.apiService.Create(c.Request().Context(), api); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create api: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Api successfully created")
}

func (p *ApiHandlers) UpdateApi(c echo.Context) error {
	idParam := c.Param("id")
	apiID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid api ID")
	}

	var updateRequest requests.UpdateApiRequest
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

	data := domain.UpdateApiRequest{
		ApiID:   uint(apiID),
		Path:    updateRequest.Path,
		Method:  updateRequest.Method,
		Summary: updateRequest.Summary,
		Tags:    updateRequest.Tags,
	}

	if _, err := p.apiService.Update(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update api: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Api successfully updated")
}

func (p *ApiHandlers) DeleteApi(c echo.Context) error {
	idParam := c.Param("id")
	apiID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid api ID")
	}

	data := domain.DeleteApiRequest{
		ApiID: uint(apiID),
	}

	if err := p.apiService.Delete(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to delete api: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Api successfully deleted")
}
