package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/models"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return fmt.Errorf("execute insert post query: %w", err)
	}

	return nil
}

// In repositories/post.go
func (r *PostRepository) GetPosts(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post

	result := r.db.WithContext(ctx).Find(&posts)

	if result.Error != nil {
		return nil, fmt.Errorf("execute select posts query: %w", result.Error)
	}

	return posts, nil
}

func (r *PostRepository) GetPost(ctx context.Context, id uint) (models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Post{}, errors.Join(models.ErrPostNotFound, err)
	} else if err != nil {
		return models.Post{}, fmt.Errorf("execute select post by id query: %w", err)
	}

	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, post *models.Post) error {
	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		return fmt.Errorf("execute update post query: %w", err)
	}

	return nil
}

func (r *PostRepository) Delete(ctx context.Context, post *models.Post) error {
	if err := r.db.WithContext(ctx).Delete(post).Error; err != nil {
		return fmt.Errorf("execute delete post query: %w", err)
	}

	return nil
}
