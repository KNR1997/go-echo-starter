package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetMenus(ctx context.Context) ([]models.Menu, error) {
	var menus []models.Menu

	result := r.db.WithContext(ctx).Find(&menus)

	if result.Error != nil {
		return nil, fmt.Errorf("executes select menus query: %w", result.Error)
	}

	return menus, nil
}

func (r *MenuRepository) GetMenuPaginated(
	ctx context.Context,
	pagination domain.Pagination,
) ([]models.Menu, int64, error) {

	var menus []models.Menu
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&models.Menu{}).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"count menus: %w",
			err,
		)
	}

	if err := r.db.WithContext(ctx).
		Limit(pagination.PageSize).
		Offset(pagination.Offset()).
		Order("id DESC").
		Find(&menus).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"select menus: %w",
			err,
		)
	}

	return menus, total, nil
}

func (r *MenuRepository) GetById(ctx context.Context, id uint) (models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Menu{}, errors.Join(models.ErrPostNotFound, err)
	} else if err != nil {
		return models.Menu{}, fmt.Errorf("execute select menu by id query: %w", err)
	}

	return menu, nil
}

func (r *MenuRepository) Create(ctx context.Context, menu *models.Menu) error {
	if err := r.db.WithContext(ctx).Create(menu).Error; err != nil {
		return fmt.Errorf("execute insert menu query: %w", err)
	}

	return nil
}

func (r *MenuRepository) Update(ctx context.Context, menu *models.Menu) error {
	if err := r.db.WithContext(ctx).Save(menu).Error; err != nil {
		return fmt.Errorf("execute update menu query: %w", err)
	}

	return nil
}

func (r *MenuRepository) Delete(ctx context.Context, menu *models.Menu) error {
	if err := r.db.WithContext(ctx).Delete(menu).Error; err != nil {
		return fmt.Errorf("execute delete menu query: %w", err)
	}

	return nil
}
