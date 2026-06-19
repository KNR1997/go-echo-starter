package repositories

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/utils"

	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, auditLog *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(auditLog).Error
}

func (r *AuditLogRepository) GetAuditLogPaginated(
	ctx context.Context,
	pagination domain.Pagination,
	searchConditions []utils.SearchCondition,
	searchJoin string,
	method string,
) ([]models.AuditLog, int64, error) {

	var auditLogs []models.AuditLog
	var total int64

	// Build the base query
	query := r.db.WithContext(ctx).Model(&models.AuditLog{})

	// Apply method filter if provided
	if method != "" {
		// Validate method to prevent SQL injection
		validMethods := map[string]bool{
			"GET":     true,
			"POST":    true,
			"PUT":     true,
			"DELETE":  true,
			"PATCH":   true,
			"HEAD":    true,
			"OPTIONS": true,
		}

		// Convert to uppercase for case-insensitive comparison
		// methodUpper := method
		// if len(method) > 0 {
		// 	methodUpper = method // You can use strings.ToUpper(method) if you import "strings"
		// }

		if validMethods[method] {
			query = query.Where("method = ?", method)
		}
		// If method is invalid, you might want to ignore it or return an error
		// For now, we'll ignore invalid methods
	}

	// Apply search conditions if any
	if len(searchConditions) > 0 {
		query = r.applySearchConditions(query, searchConditions, searchJoin)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count audit_logs: %w", err)
	}

	// Apply pagination and get results
	if err := query.
		Limit(pagination.PageSize).
		Offset(pagination.Offset()).
		Order("id DESC").
		Find(&auditLogs).Error; err != nil {
		return nil, 0, fmt.Errorf("select audit_logs: %w", err)
	}

	return auditLogs, total, nil
}

func (r *AuditLogRepository) applySearchConditions(
	query *gorm.DB,
	conditions []utils.SearchCondition,
	joinOperator string,
) *gorm.DB {
	// Create a new query with the conditions
	for i, condition := range conditions {
		// Validate field to prevent SQL injection
		// Only allow specific fields that exist in the department table
		validFields := map[string]bool{
			"username": true,
			"module":   true,
			"summary":  true,
			"path":     true,
			"status":   true,
			// Add other valid fields as needed
		}

		if !validFields[condition.Field] {
			// Skip invalid fields or handle error
			continue
		}

		// Build the condition
		var conditionExpr string
		var conditionValue interface{}

		switch condition.Operator {
		case "LIKE":
			conditionExpr = fmt.Sprintf("%s LIKE ?", condition.Field)
			conditionValue = fmt.Sprintf("%%%s%%", condition.Value)
		case "EQ":
			conditionExpr = fmt.Sprintf("%s = ?", condition.Field)
			conditionValue = condition.Value
		default:
			conditionExpr = fmt.Sprintf("%s LIKE ?", condition.Field)
			conditionValue = fmt.Sprintf("%%%s%%", condition.Value)
		}

		// Apply the condition with the appropriate join operator
		if i == 0 {
			query = query.Where(conditionExpr, conditionValue)
		} else {
			if joinOperator == "or" {
				query = query.Or(conditionExpr, conditionValue)
			} else {
				query = query.Where(conditionExpr, conditionValue)
			}
		}
	}

	return query
}
