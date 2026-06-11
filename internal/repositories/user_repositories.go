package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("execute insert user query: %w", err)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("execute update user query: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, errors.Join(models.ErrUserNotFound, err)
	} else if err != nil {
		return models.User{}, fmt.Errorf("execute select user by id query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, errors.Join(models.ErrUserNotFound, err)
	} else if err != nil {
		return models.User{}, fmt.Errorf("execute select user by email query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) CreateUserAndOAuthProvider(ctx context.Context, user *models.User, oAuthProvider *models.OAuthProviders) error {
	tx := r.db.Begin()

	committed := false

	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("insert user (tx): %w", err)
	}

	oAuthProvider.UserID = user.ID

	if err := tx.WithContext(ctx).Create(oAuthProvider).Error; err != nil {
		return fmt.Errorf("insert oauthprovider (tx): %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	committed = true

	return nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	result := r.db.WithContext(ctx).
		Preload("Roles").
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("executes select users query: %w", result.Error)
	}

	return users, nil
}

func (r *UserRepository) GetUserPaginated(
	ctx context.Context,
	pagination domain.Pagination,
) ([]models.User, int64, error) {

	var users []models.User
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"count users: %w",
			err,
		)
	}

	if err := r.db.WithContext(ctx).
		Limit(pagination.PageSize).
		Offset(pagination.Offset()).
		Order("id DESC").
		Preload("Roles").
		Preload("Department").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"select users: %w",
			err,
		)
	}

	return users, total, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Select("id").
		Where("email = ?", email).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("execute exists by email query: %w", err)
	}

	return true, nil
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Select("id").
		Where("username = ?", username).
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("execute exists by username query: %w", err)
	}

	return true, nil
}

func (r *UserRepository) AssignRoles(ctx context.Context, userID uint, roleIDs []int) error {
	// First, verify the user exists
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// If no roles to assign, return early
	if len(roleIDs) == 0 {
		return nil
	}

	// Find the roles
	var roles []models.Role
	if err := r.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}

	// Check if all requested roles exist
	if len(roles) != len(roleIDs) {
		return fmt.Errorf("some roles were not found")
	}

	// Assign roles to user
	if err := r.db.WithContext(ctx).Model(&user).Association("Roles").Append(&roles); err != nil {
		return fmt.Errorf("failed to assign roles: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Delete(user).Error; err != nil {
		return fmt.Errorf("execute delete user query: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	var user models.User
	var roles []models.Role

	// First, load the user
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Then get the roles association
	err := r.db.WithContext(ctx).Model(&user).Association("Roles").Find(&roles)
	return roles, err
}

// RemoveRoles removes specified roles from a user
func (r *UserRepository) RemoveRoles(ctx context.Context, userID uint, roleIDs []int) error {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	var roles []models.Role
	if err := r.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&user).Association("Roles").Delete(&roles)
}

// ValidateRolesExist checks if the given role IDs exist in the database
func (r *UserRepository) ValidateRolesExist(ctx context.Context, roleIDs []int) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error
	return roles, err
}
