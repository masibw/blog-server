package database

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/masibw/blog-server/domain/entity"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) FindByPermalink(permalink string) (*entity.Post, error) {
	post := &entity.Post{}
	if err := r.db.Where("permalink = ?", permalink).First(post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("select user: %w", entity.ErrPostNotFound)
		}
		return nil, fmt.Errorf("select user: %w", err)
	}
	return post, nil
}

func (r *PostRepository) Store(post *entity.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("create post: %w", entity.ErrPostAlreadyExisted)
		}
		return fmt.Errorf("create post: %w", err)
	}
	return nil
}
