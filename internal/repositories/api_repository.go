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

type ApiRepository struct {
	db *gorm.DB
}

func NewApiRepository(db *gorm.DB) *ApiRepository {
	return &ApiRepository{db: db}
}

func (r *ApiRepository) GetApis(ctx context.Context) ([]models.Api, error) {
	var apis []models.Api

	result := r.db.WithContext(ctx).Find(&apis)

	if result.Error != nil {
		return nil, fmt.Errorf("executes select apis query: %w", result.Error)
	}

	return apis, nil
}

func (r *ApiRepository) GetApiPaginated(
	ctx context.Context,
	pagination domain.Pagination,
	searchConditions []utils.SearchCondition,
	searchJoin string,
) ([]models.Api, int64, error) {

	var apis []models.Api
	var total int64

	// Build the base query
	query := r.db.WithContext(ctx).Model(&models.Api{})

	// Apply search conditions if any
	if len(searchConditions) > 0 {
		query = r.applySearchConditions(query, searchConditions, searchJoin)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count apis: %w", err)
	}

	// Apply pagination and get results
	if err := query.
		Limit(pagination.PageSize).
		Offset(pagination.Offset()).
		Order("id DESC").
		Find(&apis).Error; err != nil {
		return nil, 0, fmt.Errorf("select apis: %w", err)
	}

	return apis, total, nil
}

func (r *ApiRepository) applySearchConditions(
	query *gorm.DB,
	conditions []utils.SearchCondition,
	joinOperator string,
) *gorm.DB {
	// Create a new query with the conditions
	for i, condition := range conditions {
		// Validate field to prevent SQL injection
		// Only allow specific fields that exist in the department table
		validFields := map[string]bool{
			"path": true,
			"tags": true,
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

func (r *ApiRepository) GetById(ctx context.Context, id uint) (models.Api, error) {
	var api models.Api
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&api).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Api{}, errors.Join(models.ErrPostNotFound, err)
	} else if err != nil {
		return models.Api{}, fmt.Errorf("execute select api by id query: %w", err)
	}

	return api, nil
}

func (r *ApiRepository) Create(ctx context.Context, api *models.Api) error {
	if err := r.db.WithContext(ctx).Create(api).Error; err != nil {
		return fmt.Errorf("execute insert api query: %w", err)
	}

	return nil
}

func (r *ApiRepository) Update(ctx context.Context, api *models.Api) error {
	if err := r.db.WithContext(ctx).Save(api).Error; err != nil {
		return fmt.Errorf("execute update api query: %w", err)
	}

	return nil
}

func (r *ApiRepository) Delete(ctx context.Context, api *models.Api) error {
	if err := r.db.WithContext(ctx).Delete(api).Error; err != nil {
		return fmt.Errorf("execute delete api query: %w", err)
	}

	return nil
}
