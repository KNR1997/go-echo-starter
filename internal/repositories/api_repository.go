package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"

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
) ([]models.Api, int64, error) {

	var apis []models.Api
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&models.Api{}).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"count apis: %w",
			err,
		)
	}

	if err := r.db.WithContext(ctx).
		Limit(pagination.PageSize).
		Offset(pagination.Offset()).
		Order("id DESC").
		Find(&apis).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"select apis: %w",
			err,
		)
	}

	return apis, total, nil
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
