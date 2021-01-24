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

func (r *PostRepository) FindByID(id string) (*entity.Post, error) {
	post := &entity.Post{}
	if err := r.db.Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("find post: %w", entity.ErrPostNotFound)
			return nil, err
		}
		err = fmt.Errorf("find post: %w", err)
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) FindByPermalink(permalink string) (*entity.Post, error) {
	post := &entity.Post{}
	if err := r.db.Where("permalink = ?", permalink).First(post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find post: %w", entity.ErrPostNotFound)
		}
		return nil, fmt.Errorf("find post: %w", err)
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

func (r *PostRepository) FindAll() (posts []*entity.Post, err error) {
	if err = r.db.Find(&posts).Error; err != nil {
		err = fmt.Errorf("find all posts: %w", err)
	}
	if len(posts) == 0 {
		err = fmt.Errorf("find all posts: %w", entity.ErrPostNotFound)
	}
	return
}

func (r *PostRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Post{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("delete post: %w", entity.ErrPostNotFound)
	}
	if err := result.Error; err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}
