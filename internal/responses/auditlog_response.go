package responses

import (
	"encoding/json"
	"go-echo-starter/internal/models"
	"time"
)

type AuditLogResponse struct {
	ID           uint            `json:"id" example:"1"`
	Username     string          `json:"username" example:"IT"`
	Module       string          `json:"module" example:"notes"`
	Summary      string          `json:"summary" example:"notes"`
	Method       string          `json:"method" example:"notes"`
	Path         string          `json:"path" example:"notes"`
	Status       int             `json:"status" example:"notes"`
	ResponseTime int             `json:"response_time" example:"notes"`
	RequestArgs  json.RawMessage `json:"request_body" example:"notes"`
	ResponseBody json.RawMessage `json:"response_body" example:"notes"`
	CreatedAt    time.Time       `json:"created_at" created_at:"notes"`
}

func NewAuditLogResponse(auditLogs []models.AuditLog) *[]AuditLogResponse {
	auditLogResponse := make([]AuditLogResponse, 0)

	for i := range auditLogs {
		auditLogResponse = append(auditLogResponse, AuditLogResponse{
			ID:           auditLogs[i].ID,
			Username:     auditLogs[i].Username,
			Module:       auditLogs[i].Module,
			Summary:      auditLogs[i].Summary,
			Method:       auditLogs[i].Method,
			Path:         auditLogs[i].Path,
			Status:       auditLogs[i].Status,
			ResponseTime: auditLogs[i].ResponseTime,
			RequestArgs:  auditLogs[i].RequestArgs,
			ResponseBody: auditLogs[i].ResponseBody,
			CreatedAt:    auditLogs[i].CreatedAt,
		})
	}

	return &auditLogResponse
}
