package database

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/masibw/blog-server/domain/entity"
	"gorm.io/gorm"
)

type PostsTagsRepository struct {
	db *gorm.DB
}

func NewPostsTagsRepository(db *gorm.DB) *PostsTagsRepository {
	return &PostsTagsRepository{db: db}
}

func (r *PostsTagsRepository) FindByPostIDAndTagName(postID, tagName string) (*entity.PostsTags, error) {
	postsTags := &entity.PostsTags{}
	if err := r.db.Debug().Joins("JOIN tags ON tags.id = posts_tags.tag_id").Where("posts_tags.post_id = ? AND tags.name = ? ", postID, tagName).First(&postsTags).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find posts_tags: %w", entity.ErrPostsTagsNotFound)
		}
		return nil, fmt.Errorf("find posts_tags: %w", err)
	}
	return postsTags, nil
}

func (r *PostsTagsRepository) Store(postsTags []*entity.PostsTags) error {
	if err := r.db.Create(postsTags).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("create posts_tags: %w", entity.ErrPostsTagsAlreadyExisted)
		}
		return fmt.Errorf("create posts_tags: %w", err)
	}
	return nil
}

func (r *PostsTagsRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&entity.PostsTags{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("delete posts_tags: %w", entity.ErrPostsTagsNotFound)
	}
	if err := result.Error; err != nil {
		return fmt.Errorf("delete posts_tags: %w", err)
	}
	return nil
}

func (r *PostsTagsRepository) DeleteByPostID(postID string) error {
	result := r.db.Where("post_id = ?", postID).Delete(&entity.PostsTags{})
	if err := result.Error; err != nil {
		return fmt.Errorf("delete posts_tags: %w", err)
	}
	return nil
}
