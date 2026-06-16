package user

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/requests"
	"go-echo-starter/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

//go:generate go tool mockgen -source=$GOFILE -destination=service_mock_test.go -package=${GOPACKAGE}_test -typed=true

type userRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUserAndOAuthProvider(ctx context.Context, user *models.User, oauthProvider *models.OAuthProviders) error
	GetUsers(ctx context.Context) ([]models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	AssignRoles(ctx context.Context, userID uint, roleIDs []int) error
	GetUserPaginated(
		ctx context.Context,
		pagination domain.Pagination,
		searchConditions []utils.SearchCondition,
		searchJoin string,
		deptId int,
	) ([]models.User, int64, error)
	Delete(ctx context.Context, post *models.User) error
	GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error)
	RemoveRoles(ctx context.Context, userID uint, rolesToRemove []int) error
	ValidateRolesExist(ctx context.Context, rolesToAdd []int) ([]models.Role, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
}

type Service struct {
	userRepository userRepository
}

func NewService(userRepository userRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Register(ctx context.Context, request *requests.RegisterRequest) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	user := &models.User{
		Email:       request.Email,
		Name:        request.Name,
		Username:    request.Username,
		IsActive:    true,
		IsSuperUser: true,
		Password:    string(encryptedPassword),
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("create user in repository: %w", err)
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (models.User, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id from repository: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by email from repository: %w", err)
	}

	return user, nil
}

func (s *Service) CreateUserAndOAuthProvider(ctx context.Context, user *models.User, oauthProvider *models.OAuthProviders) error {
	err := s.userRepository.CreateUserAndOAuthProvider(ctx, user, oauthProvider)
	if err != nil {
		return fmt.Errorf("create user and oauth provider from repository: %w", err)
	}

	return nil
}

func (s *Service) GetUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepository.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil
}

func (s *Service) GetUserPaginated(
	ctx context.Context,
	pagination domain.Pagination,
	searchConditions []utils.SearchCondition,
	searchJoin string,
	deptId int,
) ([]models.User, int64, error) {

	users, total, err := s.userRepository.GetUserPaginated(
		ctx,
		pagination,
		searchConditions,
		searchJoin,
		deptId,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"get users from repository: %w",
			err,
		)
	}

	return users, total, nil
}

func (s *Service) Create(ctx context.Context, request *requests.CreateUserRequest) error {
	exists, err := s.userRepository.ExistsByEmail(ctx, request.Email)
	if err != nil {
		return fmt.Errorf("check user exists: %w", err)
	}

	if exists {
		return fmt.Errorf("user email already exists")
	}

	existsByUsername, err := s.userRepository.ExistsByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("check user exists: %w", err)
	}

	if existsByUsername {
		return fmt.Errorf("user username already exists")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	user := &models.User{
		Email: request.Email,
		// Name:     request.Name,
		Username:    request.Username,
		Password:    string(encryptedPassword),
		IsSuperUser: request.IsSuperUser,
		IsActive:    request.IsActive,
		DeptId:      request.DeptId,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return fmt.Errorf("create user in repository: %w", err)
	}

	// Assign roles if any
	if len(request.RoleIds) > 0 {
		if err := s.userRepository.AssignRoles(ctx, user.ID, request.RoleIds); err != nil {
			// Optionally: log error but don't fail the user creation?
			// Or return error to rollback the transaction
			slog.Error("assign roles to user: %w", err)
			return fmt.Errorf("assign roles to user: %w", err)
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, request domain.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepository.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("get stored user from repository: %w", err)
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
	// user.IsSuperUser = request.IsSuperUser # Todo -> disable IsSuperUser update
	user.IsActive = request.IsActive
	if request.DeptId != nil {
		user.DeptId = request.DeptId
	}

	if err := s.userRepository.Update(ctx, &user); err != nil {
		return nil, fmt.Errorf("update user in repository: %w", err)
	}

	if err := s.updateUserRoles(ctx, user.ID, request.RoleIds); err != nil {
		return nil, fmt.Errorf("update user roles: %w", err)
	}

	return &user, nil
}

func (s *Service) Patch(ctx context.Context, request domain.PatchUserRequest) (*models.User, error) {
	user, err := s.userRepository.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("get stored User from repository: %w", err)
	}

	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}

	if err := s.userRepository.Update(ctx, &user); err != nil {
		return nil, fmt.Errorf("update User in repository: %w", err)
	}

	return &user, nil
}

func (s *Service) updateUserRoles(ctx context.Context, userID uint, newRoleIDs []int) error {
	// Get current roles
	currentRoles, err := s.userRepository.GetUserRoles(ctx, userID)
	if err != nil {
		return fmt.Errorf("get current roles: %w", err)
	}

	// Convert to maps for easier comparison
	currentRoleMap := make(map[int]bool)
	for _, role := range currentRoles {
		currentRoleMap[int(role.ID)] = true
	}

	newRoleMap := make(map[int]bool)
	for _, roleID := range newRoleIDs {
		newRoleMap[roleID] = true
	}

	// Find roles to add (in new but not in current)
	rolesToAdd := []int{}
	for _, roleID := range newRoleIDs {
		if !currentRoleMap[roleID] {
			rolesToAdd = append(rolesToAdd, roleID)
		}
	}

	// Find roles to remove (in current but not in new)
	rolesToRemove := []int{}
	for _, role := range currentRoles {
		if !newRoleMap[int(role.ID)] {
			rolesToRemove = append(rolesToRemove, int(role.ID))
		}
	}

	// Remove roles
	if len(rolesToRemove) > 0 {
		if err := s.userRepository.RemoveRoles(ctx, userID, rolesToRemove); err != nil {
			return fmt.Errorf("remove roles: %w", err)
		}
	}

	// Add new roles
	if len(rolesToAdd) > 0 {
		// Validate roles exist
		validRoles, err := s.userRepository.ValidateRolesExist(ctx, rolesToAdd)
		if err != nil {
			return fmt.Errorf("validate roles: %w", err)
		}

		if len(validRoles) != len(rolesToAdd) {
			return fmt.Errorf("some roles do not exist")
		}

		if err := s.userRepository.AssignRoles(ctx, userID, rolesToAdd); err != nil {
			return fmt.Errorf("assign roles: %w", err)
		}
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, request domain.DeleteUserRequest) error {
	user, err := s.userRepository.GetByID(ctx, request.UserID)
	if err != nil {
		return fmt.Errorf("get stored user from repository: %w", err)
	}

	if err := s.userRepository.Delete(ctx, &user); err != nil {
		return fmt.Errorf("delete user in repository: %w", err)
	}

	return nil
}

func (s *Service) UpdateLastLogin(ctx context.Context, user models.User) error {
	if err := s.userRepository.UpdateLastLogin(context.Background(), user.ID); err != nil {
		log.Printf("Warning: Failed to update last_login: %v", err)
	}
	return nil
}
