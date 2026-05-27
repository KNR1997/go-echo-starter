package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/models"

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
