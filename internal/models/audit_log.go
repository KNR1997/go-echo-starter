package models

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	UserID       int             `gorm:"index;not null" json:"user_id"`
	Username     string          `gorm:"size:64;index;default:''" json:"username"`
	Module       string          `gorm:"size:64;index;default:''" json:"module"`
	Summary      string          `gorm:"size:128;index;default:''" json:"summary"`
	Method       string          `gorm:"size:10;index;not null" json:"method"`
	Path         string          `gorm:"size:255;index;default:''" json:"path"`
	Status       int             `gorm:"index;default:-1" json:"status"`
	ResponseTime int             `gorm:"index;default:0" json:"response_time"`
	RequestArgs  json.RawMessage `gorm:"type:jsonb" json:"request_args"`
	ResponseBody json.RawMessage `gorm:"type:jsonb" json:"response_body"`
	CreatedAt    time.Time       `gorm:"index;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
