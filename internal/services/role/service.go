package role

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/utils"
	"log/slog"
)

type roleRepository interface {
	GetById(ctx context.Context, id uint) (models.Role, error)
	GetRoles(ctx context.Context) ([]models.Role, error)
	GetRolePaginated(
		ctx context.Context,
		pagination domain.Pagination,
		searchConditions []utils.SearchCondition,
		searchJoin string,
	) ([]models.Role, int64, error)
	Create(ctx context.Context, dept *models.Role) error
	Update(ctx context.Context, dept *models.Role) error
	Delete(ctx context.Context, post *models.Role) error
	AssignMenus(ctx context.Context, roleID uint, MenuIDs []int) error
	AssignApis(ctx context.Context, roleID uint, ApiIDs []int) error
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

func (s *Service) GetRolePaginated(
	ctx context.Context,
	pagination domain.Pagination,
	searchConditions []utils.SearchCondition,
	searchJoin string,
) ([]models.Role, int64, error) {

	roles, total, err := s.roleRepository.GetRolePaginated(
		ctx,
		pagination,
		searchConditions,
		searchJoin,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"get roles from repository: %w",
			err,
		)
	}

	return roles, total, nil
}

func (s *Service) Create(ctx context.Context, dept *models.Role) error {
	if err := s.roleRepository.Create(ctx, dept); err != nil {
		return fmt.Errorf("create role in repository: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, request domain.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepository.GetById(ctx, request.RoleID)
	if err != nil {
		return nil, fmt.Errorf("get stored role from repository: %w", err)
	}

	role.Name = request.Name
	role.Description = request.Description

	if err := s.roleRepository.Update(ctx, &role); err != nil {
		return nil, fmt.Errorf("update role in repository: %w", err)
	}

	return &role, nil
}

func (s *Service) Delete(ctx context.Context, request domain.DeleteRoleRequest) error {
	role, err := s.roleRepository.GetById(ctx, request.RoleID)
	if err != nil {
		return fmt.Errorf("get stored role from repository: %w", err)
	}

	if err := s.roleRepository.Delete(ctx, &role); err != nil {
		return fmt.Errorf("delete role in repository: %w", err)
	}

	return nil
}

func (s *Service) Authorize(ctx context.Context, request domain.AuthorizeRoleRequest) (*models.Role, error) {
	role, err := s.roleRepository.GetById(ctx, request.RoleID)
	if err != nil {
		return nil, fmt.Errorf("get stored role from repository: %w", err)
	}

	if err := s.roleRepository.AssignMenus(ctx, role.ID, request.MenuIDs); err != nil {
		// Optionally: log error but don't fail the user creation?
		// Or return error to rollback the transaction
		slog.Error("assign menus to role: %w", err)
		return nil, fmt.Errorf("assign menus to role: %w", err)
	}

	if err := s.roleRepository.AssignApis(ctx, role.ID, request.ApiIDs); err != nil {
		// Optionally: log error but don't fail the user creation?
		// Or return error to rollback the transaction
		slog.Error("assign apis to role: %w", err)
		return nil, fmt.Errorf("assign apis to role: %w", err)
	}

	return &role, nil
}
