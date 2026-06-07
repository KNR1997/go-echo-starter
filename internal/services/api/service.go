package api

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
)

type apiRepository interface {
	GetApis(ctx context.Context) ([]models.Api, error)
	GetApiPaginated(ctx context.Context, pagination domain.Pagination) ([]models.Api, int64, error)
	GetById(ctx context.Context, id uint) (models.Api, error)
	Create(ctx context.Context, dept *models.Api) error
	Update(ctx context.Context, dept *models.Api) error
	Delete(ctx context.Context, post *models.Api) error
}

type Service struct {
	apiRepository apiRepository
}

func NewService(apiRepository apiRepository) *Service {
	return &Service{apiRepository: apiRepository}
}

func (s *Service) GetApis(ctx context.Context) ([]models.Api, error) {
	apis, err := s.apiRepository.GetApis(ctx)
	if err != nil {
		return nil, fmt.Errorf("get apis from repository: %w", err)
	}

	return apis, nil
}

func (s *Service) GetApiPaginated(
	ctx context.Context,
	pagination domain.Pagination,
) ([]models.Api, int64, error) {

	apis, total, err := s.apiRepository.GetApiPaginated(
		ctx,
		pagination,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"get apis from repository: %w",
			err,
		)
	}

	return apis, total, nil
}

func (s *Service) Create(ctx context.Context, dept *models.Api) error {
	if err := s.apiRepository.Create(ctx, dept); err != nil {
		return fmt.Errorf("create Api in repository: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, request domain.UpdateApiRequest) (*models.Api, error) {
	api, err := s.apiRepository.GetById(ctx, request.ApiID)
	if err != nil {
		return nil, fmt.Errorf("get stored Api from repository: %w", err)
	}

	api.Path = request.Path
	api.Method = request.Method
	api.Summary = request.Summary
	api.Tags = request.Tags

	if err := s.apiRepository.Update(ctx, &api); err != nil {
		return nil, fmt.Errorf("update Api in repository: %w", err)
	}

	return &api, nil
}

func (s *Service) Delete(ctx context.Context, request domain.DeleteApiRequest) error {
	dept, err := s.apiRepository.GetById(ctx, request.ApiID)
	if err != nil {
		return fmt.Errorf("get stored Api from repository: %w", err)
	}

	if err := s.apiRepository.Delete(ctx, &dept); err != nil {
		return fmt.Errorf("delete Api in repository: %w", err)
	}

	return nil
}
