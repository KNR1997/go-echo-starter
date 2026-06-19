package handlers

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/responses"
	"go-echo-starter/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type auditLogSerivce interface {
	GetAuditLogPaginated(
		ctx context.Context,
		pagination domain.Pagination,
		searchConditions []utils.SearchCondition,
		searchJoin string,
		method string,
	) ([]models.AuditLog, int64, error)
}

type AudtiLogHandlers struct {
	auditLogSerivce auditLogSerivce
}

func NewAudtiLogHandlers(service auditLogSerivce) *AudtiLogHandlers {
	return &AudtiLogHandlers{auditLogSerivce: service}
}

func (h AudtiLogHandlers) GetAuditLogPaginated(c echo.Context) error {
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

	method := c.QueryParam("method")

	fmt.Println("method================: ", method)

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

	auditLogs, total, err := h.auditLogSerivce.GetAuditLogPaginated(
		c.Request().Context(),
		pagination,
		searchConditions,
		searchJoin,
		method,
	)
	if err != nil {
		return responses.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get auditLogs",
		)
	}

	return responses.Response(c, http.StatusOK, map[string]any{
		"data":     responses.NewAuditLogResponse(auditLogs),
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}
