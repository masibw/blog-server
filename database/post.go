package database

import (
	"errors"
	"fmt"

	"github.com/masibw/blog-server/domain/entity"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Store(post *entity.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		if errors.Is(err, gorm.ErrRegistered) {
			return fmt.Errorf("create post: %w", entity.ErrPostAlreadyExisted)
		}
		return fmt.Errorf("create post: %w", err)
	}
	return nil
}
