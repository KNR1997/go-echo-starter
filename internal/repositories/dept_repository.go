package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/utils"

	"gorm.io/gorm"
)

type DeptRepository struct {
	db *gorm.DB
}

func NewDeptRepository(db *gorm.DB) *DeptRepository {
	return &DeptRepository{db: db}
}

func (r *DeptRepository) GetDepartments(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department

	result := r.db.WithContext(ctx).Find(&departments)

	if result.Error != nil {
		return nil, fmt.Errorf("executes select departments query: %w", result.Error)
	}

	return departments, nil
}

func (r *DeptRepository) GetDepartmentPaginated(
	ctx context.Context,
	pagination domain.Pagination,
	searchConditions []utils.SearchCondition,
	searchJoin string,
) ([]models.Department, int64, error) {

	var departments []models.Department
	var total int64

	// Build the base query
	query := r.db.WithContext(ctx).Model(&models.Department{})

	// Apply search conditions if any
	if len(searchConditions) > 0 {
		query = r.applySearchConditions(query, searchConditions, searchJoin)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count departments: %w", err)
	}

	// Apply pagination and get results
	if err := query.
		Limit(pagination.PageSize).
		Offset(pagination.Offset()).
		Order("id DESC").
		Find(&departments).Error; err != nil {
		return nil, 0, fmt.Errorf("select departments: %w", err)
	}

	return departments, total, nil
}

func (r *DeptRepository) applySearchConditions(
	query *gorm.DB,
	conditions []utils.SearchCondition,
	joinOperator string,
) *gorm.DB {
	// Create a new query with the conditions
	for i, condition := range conditions {
		// Validate field to prevent SQL injection
		// Only allow specific fields that exist in the department table
		validFields := map[string]bool{
			"name":        true,
			"description": true,
			"code":        true,
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

func (r *DeptRepository) GetById(ctx context.Context, id uint) (models.Department, error) {
	var dept models.Department
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&dept).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Department{}, errors.Join(models.ErrPostNotFound, err)
	} else if err != nil {
		return models.Department{}, fmt.Errorf("execute select dept by id query: %w", err)
	}

	return dept, nil
}

func (r *DeptRepository) Create(ctx context.Context, dept *models.Department) error {
	if err := r.db.WithContext(ctx).Create(dept).Error; err != nil {
		return fmt.Errorf("execute insert dept query: %w", err)
	}

	return nil
}

func (r *DeptRepository) Update(ctx context.Context, dept *models.Department) error {
	if err := r.db.WithContext(ctx).Save(dept).Error; err != nil {
		return fmt.Errorf("execute update dept query: %w", err)
	}

	return nil
}

func (r *DeptRepository) Delete(ctx context.Context, dept *models.Department) error {
	if err := r.db.WithContext(ctx).Delete(dept).Error; err != nil {
		return fmt.Errorf("execute delete dept query: %w", err)
	}

	return nil
}

func (r *DeptRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var dept models.Department

	err := r.db.WithContext(ctx).
		Select("id").
		Where("name = ?", name).
		Take(&dept).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("execute exists by name query: %w", err)
	}

	return true, nil
}
