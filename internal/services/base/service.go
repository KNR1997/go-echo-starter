package base

import (
	"context"
	"fmt"
	"go-echo-starter/internal/models"
)

type userRepository interface {
	GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) // Changed roleID to userID
}

type roleRepository interface {
	GetRoleMenus(ctx context.Context, roleID uint) ([]models.Menu, error)
}

type Service struct {
	userRepository userRepository
	roleRepository roleRepository
}

func NewService(userRepository userRepository, roleRepository roleRepository) *Service { // Added userRepository parameter
	return &Service{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

func (s *Service) GetUserMenus(ctx context.Context, userID uint) ([]models.Menu, error) {
	// Get all roles for the user
	roles, err := s.userRepository.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles from repository: %w", err)
	}

	// Collect unique menus from all roles
	menuMap := make(map[uint]models.Menu) // Using map to deduplicate menus by ID

	for _, role := range roles {
		menus, err := s.roleRepository.GetRoleMenus(ctx, role.ID) // Use role.ID, not userID
		if err != nil {
			return nil, fmt.Errorf("get role menus for role %d: %w", role.ID, err)
		}

		// Add menus to map (automatically deduplicates by ID)
		for _, menu := range menus {
			menuMap[menu.ID] = menu
		}
	}

	// Convert map values to slice
	uniqueMenus := make([]models.Menu, 0, len(menuMap))
	for _, menu := range menuMap {
		uniqueMenus = append(uniqueMenus, menu)
	}

	return uniqueMenus, nil
}
