package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/models"

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
