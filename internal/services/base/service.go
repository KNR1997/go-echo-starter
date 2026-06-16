package base

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error)
	GetByID(ctx context.Context, userID uint) (models.User, error)
	Update(ctx context.Context, user *models.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type roleRepository interface {
	GetRoleMenus(ctx context.Context, roleID uint) ([]models.Menu, error)
}

type menuRepository interface {
	GetMenus(ctx context.Context) ([]models.Menu, error)
	Create(ctx context.Context, dept *models.Menu) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type Service struct {
	userRepository userRepository
	roleRepository roleRepository
	menuRepository menuRepository
}

func NewService(userRepository userRepository, roleRepository roleRepository, menuRepository menuRepository) *Service { // Added userRepository parameter
	return &Service{
		userRepository: userRepository,
		roleRepository: roleRepository,
		menuRepository: menuRepository,
	}
}

func (s *Service) GetUserMenus(ctx context.Context, userID uint) ([]models.Menu, error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	if user.IsSuperUser {
		menus, err := s.menuRepository.GetMenus(ctx)
		if err != nil {
			return nil, fmt.Errorf("get menus from repository: %w", err)
		}
		return menus, nil
	}

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

func (s *Service) GetMeDetails(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	return &user, nil
}

func (s *Service) ProfileUpdate(ctx context.Context, request domain.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepository.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	if user.Email != request.Email {
		exists, err := s.userRepository.ExistsByEmail(ctx, request.Email)
		if err != nil {
			return nil, fmt.Errorf("check user exists: %w", err)
		}

		if exists {
			return nil, fmt.Errorf("user email already exists")
		}
	}

	user.Username = request.UserName
	user.Email = request.Email

	if err := s.userRepository.Update(ctx, &user); err != nil {
		return nil, fmt.Errorf("update user in repository: %w", err)
	}

	return &user, nil
}

func (s *Service) PasswordUpdate(ctx context.Context, request domain.UpdatePasswordRequest) (*models.User, error) {
	// Get user by ID
	user, err := s.userRepository.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user from repository: %w", err)
	}

	// Verify current password if required
	if request.OldPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
		if err != nil {
			return nil, fmt.Errorf("current password is incorrect: %w", err)
		}
	}

	// Validate new password strength
	if len(request.NewPassword) < 6 {
		return nil, fmt.Errorf("new password must be at least 6 characters long")
	}

	// Encrypt the new password
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("encrypt password: %w", err)
	}

	// Update user password
	user.Password = string(encryptedPassword)

	//  Save to repository
	if err := s.userRepository.Update(ctx, &user); err != nil {
		return nil, fmt.Errorf("update user in repository: %w", err)
	}

	//  Return updated user
	user.Password = "" // Don't return password hash to client
	return &user, nil
}

func (s *Service) InitiateAdmin(ctx context.Context) error {
	user, err := s.userRepository.GetUserByEmail(ctx, "admin@demo.com")
	if err != nil {
		// Check if this is a "not found" error
		if !strings.Contains(err.Error(), "not found") { // Adjust based on your error
			return fmt.Errorf("check user exists: %w", err)
		}

		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte("demodemo"),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return fmt.Errorf("encrypt password: %w", err)
		}

		// User not found, create new admin
		newUser := &models.User{
			Email:       "admin@demo.com",
			Username:    "JohnDoe",
			Password:    string(encryptedPassword),
			IsActive:    true,
			IsSuperUser: true,
		}
		if err := s.userRepository.Create(ctx, newUser); err != nil {
			return fmt.Errorf("create user in repository: %w", err)
		}
		return nil
	}

	// User exists, update to superuser
	user.IsSuperUser = true
	if err := s.userRepository.Update(ctx, &user); err != nil {
		return fmt.Errorf("update user in repository: %w", err)
	}

	return nil
}

func (s *Service) InitiateMenus(ctx context.Context) error {
	menuTypeCatalog := "catalog"
	menuTypeMenu := "menus"
	redirect := "/settings/user"
	userIcon := "material-symbols:person-outline-rounded"
	roleIcon := "carbon:user-role"
	apiIcon := "ant-design:api-outlined"
	departmentIcon := "mingcute:department-line"
	menuIcon := "material-symbols:list-alt-outline"

	menus := []models.Menu{
		{
			ID:          1,
			Name:        "Settings",
			MenusType:   &menuTypeCatalog,
			Icon:        &menuIcon,
			Path:        "/settings",
			OrderNumber: 1,
			IsHidden:    false,
			Keepalive:   false,
			Redirect:    &redirect,
			Component:   "Layout",
			ParentID:    0,
		},
		{
			ID:          2,
			Name:        "users",
			MenusType:   &menuTypeMenu,
			Icon:        &userIcon,
			Path:        "user",
			OrderNumber: 1,
			IsHidden:    false,
			Keepalive:   false,
			Redirect:    nil,
			Component:   "/settings/user",
			ParentID:    1,
		},
		{
			ID:          3,
			Name:        "roles",
			MenusType:   &menuTypeMenu,
			Icon:        &roleIcon,
			Path:        "role",
			OrderNumber: 2,
			IsHidden:    false,
			Keepalive:   false,
			Redirect:    nil,
			Component:   "/settings/role",
			ParentID:    1,
		},
		{
			ID:          4,
			Name:        "apis",
			MenusType:   &menuTypeMenu,
			Icon:        &apiIcon,
			Path:        "api",
			OrderNumber: 3,
			IsHidden:    false,
			Keepalive:   false,
			Redirect:    nil,
			Component:   "/settings/api",
			ParentID:    1,
		},
		{
			ID:          5,
			Name:        "departments",
			MenusType:   &menuTypeMenu,
			Icon:        &departmentIcon,
			Path:        "department",
			OrderNumber: 4,
			IsHidden:    false,
			Keepalive:   false,
			Redirect:    nil,
			Component:   "/settings/department",
			ParentID:    1,
		},
		{
			ID:          6,
			Name:        "menus",
			MenusType:   &menuTypeMenu,
			Icon:        &menuIcon,
			Path:        "menu",
			OrderNumber: 5,
			IsHidden:    false,
			Keepalive:   false,
			Redirect:    nil,
			Component:   "/settings/menu",
			ParentID:    1,
		},
	}

	for _, menu := range menus {
		exists, err := s.menuRepository.ExistsByName(ctx, menu.Name)
		if err != nil {
			return fmt.Errorf("check menu exists: %w", err)
		}

		if !exists {
			if err := s.menuRepository.Create(ctx, &menu); err != nil {
				return fmt.Errorf("create Menu in repository: %w", err)
			}
		}

	}

	return nil

}
