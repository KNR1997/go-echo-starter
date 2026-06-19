package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go-echo-starter/internal/models"
	"go-echo-starter/internal/repositories"
	"go-echo-starter/internal/services/audit"
	"go-echo-starter/internal/slogx"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// auditLogMiddleware logs all requests to the audit log
type auditLogMiddleware struct {
	auditService *audit.Service
}

func NewAuditLogMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	auditRepo := repositories.NewAuditLogRepository(db)
	auditService := audit.NewService(auditRepo)
	return (&auditLogMiddleware{auditService: auditService}).handle
}

func (m *auditLogMiddleware) handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Skip audit logging for OPTIONS requests (CORS preflight)
		// if c.Request().Method == http.MethodOptions {
		// 	return next(c)
		// }

		start := time.Now()

		// Debug: Log the actual request method at the start
		slog.DebugContext(c.Request().Context(), "AuditLogMiddleware - Request received",
			"method", c.Request().Method,
			"path", c.Path(),
			"url", c.Request().URL.String(),
		)

		// Get request body
		requestBody, err := m.getRequestBody(c)
		if err != nil {
			slog.ErrorContext(c.Request().Context(), "Failed to read request body for audit log", "error", err)
		}

		// Create response storer to capture response body
		storer := newResponseStorer(c.Response().Writer)
		c.Response().Writer = storer

		// Process the request
		errNext := next(c)

		// Calculate response time
		responseTime := int(time.Since(start).Milliseconds())

		// Get user info from context
		ctx := c.Request().Context()
		userID := slogx.UserIDFromContext(ctx)
		username := slogx.UsernameFromContext(ctx)

		// Build audit log entry
		auditLog := &models.AuditLog{
			UserID:       userID,
			Username:     username,
			Module:       m.getModule(c),
			Summary:      m.getSummary(c),
			Method:       c.Request().Method,
			Path:         c.Path(),
			Status:       c.Response().Status,
			ResponseTime: responseTime,
			RequestArgs:  requestBody,
			ResponseBody: m.getResponseBody(storer),
		}

		// Save audit log asynchronously with a background context
		go func() {
			// Use context.Background() instead of the request context
			ctx := context.Background()

			// Add trace ID for logging correlation if needed
			// if traceID := slogx.TraceIDFromContext(c.Request().Context()); traceID != "" {
			// 	ctx = slogx.ContextWithTraceID(ctx, traceID)
			// }

			if err := m.auditService.Create(ctx, auditLog); err != nil {
				slog.ErrorContext(ctx, "Failed to save audit log", "error", err)
			}
		}()

		if errNext != nil {
			return fmt.Errorf("handle request with audit log: %w", errNext)
		}

		return nil
	}
}

func (m *auditLogMiddleware) getRequestBody(c echo.Context) (json.RawMessage, error) {
	if c.Request().Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}

	// Restore the body so it can be read again by handlers
	c.Request().Body = io.NopCloser(bytes.NewReader(body))

	// Check if body is empty
	if len(body) == 0 {
		return nil, nil
	}

	// Only return JSON bodies
	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
		// Validate that it's valid JSON
		var js json.RawMessage
		if err := json.Unmarshal(body, &js); err != nil {
			// Not valid JSON, don't store it
			slog.WarnContext(c.Request().Context(), "Request body is not valid JSON",
				"error", err,
				"content_type", contentType,
			)
			return nil, nil
		}
		return json.RawMessage(body), nil
	}

	// For non-JSON bodies, return nil
	return nil, nil
}

func (m *auditLogMiddleware) getResponseBody(storer *responseStorer) json.RawMessage {
	if storer.storedResponse == nil || len(storer.storedResponse) == 0 {
		return nil
	}

	contentType := storer.Header().Get(echo.HeaderContentType)
	if !strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
		return nil
	}

	// Validate that it's valid JSON
	var js json.RawMessage
	if err := json.Unmarshal(storer.storedResponse, &js); err != nil {
		// Not valid JSON, don't store it
		return nil
	}

	return json.RawMessage(storer.storedResponse)
}

func (m *auditLogMiddleware) getUserInfo(c echo.Context) (int, string) {
	ctx := c.Request().Context()

	// Try to get user ID from context (set by auth middleware)
	userID := slogx.UserIDFromContext(ctx)

	// Try to get username from JWT claims or context
	username := ""
	if userID > 0 {
		// You might want to fetch username from context or JWT claims
		// For now, we'll try to get it from the context
		if val := ctx.Value("username"); val != nil {
			if str, ok := val.(string); ok {
				username = str
			}
		}
	}

	return userID, username
}

func (m *auditLogMiddleware) getModule(c echo.Context) string {
	// Extract module from path
	path := c.Path()

	// You can define module mapping based on path patterns
	if strings.HasPrefix(path, "/users") {
		return "user"
	}
	if strings.HasPrefix(path, "/roles") {
		return "role"
	}
	if strings.HasPrefix(path, "/departments") {
		return "department"
	}
	if strings.HasPrefix(path, "/apis") {
		return "api"
	}
	if strings.HasPrefix(path, "/menus") {
		return "menu"
	}
	if strings.HasPrefix(path, "/posts") {
		return "post"
	}
	if strings.HasPrefix(path, "/base") {
		return "base"
	}
	if strings.HasPrefix(path, "/auditLogs") {
		return "auditLog"
	}
	if strings.HasPrefix(path, "/login") || strings.HasPrefix(path, "/register") {
		return "auth"
	}

	return "unknown"
}

func (m *auditLogMiddleware) getSummary(c echo.Context) string {
	// Generate a summary based on the request
	method := c.Request().Method
	path := c.Path()
	status := c.Response().Status

	summary := fmt.Sprintf("%s %s - %d", method, path, status)

	// Add more context for specific endpoints
	if strings.Contains(path, "paginated") {
		summary = fmt.Sprintf("List %s", m.getModule(c))
	} else if method == http.MethodPost && !strings.Contains(path, "paginated") {
		summary = fmt.Sprintf("Create %s", m.getModule(c))
	} else if method == http.MethodPut || method == http.MethodPatch {
		summary = fmt.Sprintf("Update %s", m.getModule(c))
	} else if method == http.MethodDelete {
		summary = fmt.Sprintf("Delete %s", m.getModule(c))
	}

	return summary
}
