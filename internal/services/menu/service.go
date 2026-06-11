package menu

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
)

type menuRepository interface {
	GetMenus(ctx context.Context) ([]models.Menu, error)
	GetMenuPaginated(ctx context.Context, pagination domain.Pagination) ([]models.Menu, int64, error)
	GetById(ctx context.Context, id uint) (models.Menu, error)
	Create(ctx context.Context, dept *models.Menu) error
	Update(ctx context.Context, dept *models.Menu) error
	Delete(ctx context.Context, post *models.Menu) error
}

type Service struct {
	menuRepository menuRepository
}

func NewService(menuRepository menuRepository) *Service {
	return &Service{menuRepository: menuRepository}
}

func (s *Service) GetMenus(ctx context.Context) ([]models.Menu, error) {
	menus, err := s.menuRepository.GetMenus(ctx)
	if err != nil {
		return nil, fmt.Errorf("get menus from repository: %w", err)
	}

	return menus, nil
}

func (s *Service) GetMenuPaginated(
	ctx context.Context,
	pagination domain.Pagination,
) ([]models.Menu, int64, error) {

	menus, total, err := s.menuRepository.GetMenuPaginated(
		ctx,
		pagination,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"get menus from repository: %w",
			err,
		)
	}

	return menus, total, nil
}

func (s *Service) Create(ctx context.Context, dept *models.Menu) error {
	if err := s.menuRepository.Create(ctx, dept); err != nil {
		return fmt.Errorf("create Menu in repository: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, request domain.UpdateMenuRequest) (*models.Menu, error) {
	menu, err := s.menuRepository.GetById(ctx, request.MenuID)
	if err != nil {
		return nil, fmt.Errorf("get stored Menu from repository: %w", err)
	}

	menu.Name = request.Name
	menu.MenuType = request.MenuType
	menu.Icon = request.Icon
	menu.Path = request.Path
	menu.Order = request.Order
	menu.ParentID = request.ParentID
	menu.IsHidden = request.IsHidden
	menu.Component = request.Component
	menu.Keepalive = request.Keepalive
	menu.Redirect = request.Redirect

	if err := s.menuRepository.Update(ctx, &menu); err != nil {
		return nil, fmt.Errorf("update Menu in repository: %w", err)
	}

	return &menu, nil
}

func (s *Service) Patch(ctx context.Context, request domain.PatchMenuRequest) (*models.Menu, error) {
	menu, err := s.menuRepository.GetById(ctx, request.MenuID)
	if err != nil {
		return nil, fmt.Errorf("get stored Menu from repository: %w", err)
	}

	if request.IsHidden != nil {
		menu.IsHidden = *request.IsHidden
	}
	if request.Keepalive != nil {
		menu.Keepalive = *request.Keepalive
	}

	if err := s.menuRepository.Update(ctx, &menu); err != nil {
		return nil, fmt.Errorf("update Menu in repository: %w", err)
	}

	return &menu, nil
}

func (s *Service) Delete(ctx context.Context, request domain.DeleteMenuRequest) error {
	dept, err := s.menuRepository.GetById(ctx, request.MenuID)
	if err != nil {
		return fmt.Errorf("get stored Menu from repository: %w", err)
	}

	if err := s.menuRepository.Delete(ctx, &dept); err != nil {
		return fmt.Errorf("delete Menu in repository: %w", err)
	}

	return nil
}
