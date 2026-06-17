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

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role

	result := r.db.WithContext(ctx).Find(&roles)

	if result.Error != nil {
		return nil, fmt.Errorf("executes select roles query: %w", result.Error)
	}

	return roles, nil
}

func (r *RoleRepository) GetRolePaginated(
	ctx context.Context,
	pagination domain.Pagination,
	searchConditions []utils.SearchCondition,
	searchJoin string,
) ([]models.Role, int64, error) {

	var roles []models.Role
	var total int64

	// Build the base query
	query := r.db.WithContext(ctx).Model(&models.Role{}).
		Preload("Menus").
		Preload("Apis")

	// Apply search conditions if any
	if len(searchConditions) > 0 {
		query = r.applySearchConditions(query, searchConditions, searchJoin)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count roles: %w", err)
	}

	// Apply pagination and get results
	if err := query.
		Limit(pagination.PageSize).
		Offset(pagination.Offset()).
		Order("id DESC").
		Find(&roles).Error; err != nil {
		return nil, 0, fmt.Errorf("select roles: %w", err)
	}

	return roles, total, nil
}

func (r *RoleRepository) applySearchConditions(
	query *gorm.DB,
	conditions []utils.SearchCondition,
	joinOperator string,
) *gorm.DB {
	// Create a new query with the conditions
	for i, condition := range conditions {
		// Validate field to prevent SQL injection
		// Only allow specific fields that exist in the department table
		validFields := map[string]bool{
			"name": true,
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

func (r *RoleRepository) GetById(ctx context.Context, id uint) (models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Role{}, errors.Join(models.ErrPostNotFound, err)
	} else if err != nil {
		return models.Role{}, fmt.Errorf("execute select role by id query: %w", err)
	}

	return role, nil
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return fmt.Errorf("execute insert role query: %w", err)
	}

	return nil
}

func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	if err := r.db.WithContext(ctx).Save(role).Error; err != nil {
		return fmt.Errorf("execute update role query: %w", err)
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, role *models.Role) error {
	if err := r.db.WithContext(ctx).Delete(role).Error; err != nil {
		return fmt.Errorf("execute delete role query: %w", err)
	}

	return nil
}

func (r *RoleRepository) AssignMenus(ctx context.Context, roleID uint, menuIDs []int) error {
	// First, verify the role exists
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// If no menus to assign, clear all existing menus
	if len(menuIDs) == 0 {
		return r.db.WithContext(ctx).Model(&role).Association("Menus").Clear()
	}

	// Find the menus
	var menus []models.Menu
	if err := r.db.WithContext(ctx).Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
		return fmt.Errorf("failed to find menus: %w", err)
	}

	// Check if all requested menus exist
	if len(menus) != len(menuIDs) {
		return fmt.Errorf("some menus were not found")
	}

	// Replace all menus (this will remove old ones and add new ones)
	if err := r.db.WithContext(ctx).Model(&role).Association("Menus").Replace(&menus); err != nil {
		return fmt.Errorf("failed to replace menus: %w", err)
	}

	return nil
}

func (r *RoleRepository) AssignApis(ctx context.Context, roleID uint, apiIDs []int) error {
	// First, verify the role exists
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// If no apis to assign, clear all existing apis
	if len(apiIDs) == 0 {
		return r.db.WithContext(ctx).Model(&role).Association("Apis").Clear()
	}

	// Find the apis
	var apis []models.Api
	if err := r.db.WithContext(ctx).Where("id IN ?", apiIDs).Find(&apis).Error; err != nil {
		return fmt.Errorf("failed to find apis: %w", err)
	}

	// Check if all requested apis exist
	if len(apis) != len(apiIDs) {
		return fmt.Errorf("some apis were not found")
	}

	// Replace all apis (this will remove old ones and add new ones)
	if err := r.db.WithContext(ctx).Model(&role).Association("Apis").Replace(&apis); err != nil {
		return fmt.Errorf("failed to replace apis: %w", err)
	}

	return nil
}

func (r *RoleRepository) GetRoleMenus(ctx context.Context, roleID uint) ([]models.Menu, error) {
	var role models.Role
	var menus []models.Menu

	// First, load the role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return nil, err
	}

	// Then get the menus association
	err := r.db.WithContext(ctx).Model(&role).Association("Menus").Find(&menus)
	return menus, err
}

func (r *RoleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var role models.Role

	err := r.db.WithContext(ctx).
		Select("id").
		Where("name = ?", name).
		Take(&role).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("execute exists by name query: %w", err)
	}

	return true, nil
}
