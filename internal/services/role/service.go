package role

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
)

type roleRepository interface {
	GetById(ctx context.Context, id uint) (models.Role, error)
	GetRoles(ctx context.Context) ([]models.Role, error)
	Create(ctx context.Context, dept *models.Role) error
	Update(ctx context.Context, dept *models.Role) error
	Delete(ctx context.Context, post *models.Role) error
}

type Service struct {
	roleRepository roleRepository
}

func NewService(roleRepository roleRepository) *Service {
	return &Service{roleRepository: roleRepository}
}

func (s *Service) GetRoles(ctx context.Context) ([]models.Role, error) {
	roles, err := s.roleRepository.GetRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("get roles from repository: %w", err)
	}

	return roles, nil
}

func (s *Service) Create(ctx context.Context, dept *models.Role) error {
	if err := s.roleRepository.Create(ctx, dept); err != nil {
		return fmt.Errorf("create role in repository: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, request domain.UpdateRoleRequest) (*models.Role, error) {
	dept, err := s.roleRepository.GetById(ctx, request.RoleID)
	if err != nil {
		return nil, fmt.Errorf("get stored role from repository: %w", err)
	}

	dept.Name = request.Name
	dept.Desc = request.Desc

	if err := s.roleRepository.Update(ctx, &dept); err != nil {
		return nil, fmt.Errorf("update role in repository: %w", err)
	}

	return &dept, nil
}

func (s *Service) Delete(ctx context.Context, request domain.DeleteRoleRequest) error {
	dept, err := s.roleRepository.GetById(ctx, request.RoleID)
	if err != nil {
		return fmt.Errorf("get stored role from repository: %w", err)
	}

	if err := s.roleRepository.Delete(ctx, &dept); err != nil {
		return fmt.Errorf("delete role in repository: %w", err)
	}

	return nil
}
