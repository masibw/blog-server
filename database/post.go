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

func (r *PostRepository) Create(post *entity.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("create post: %w", entity.ErrPostAlreadyExisted)
		}
		return fmt.Errorf("create post: %w", err)
	}
	return nil
}

func (r *PostRepository) Update(post *entity.Post) error {

	if err := r.db.Select("*").Updates(post).Error; err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	return nil
}

func (r *PostRepository) FindAll(offset, pageSize int, condition string, params []interface{}, sortCondition string) (posts []*entity.Post, err error) {
	if err = r.db.Distinct().Where(condition, params...).Order(sortCondition).Limit(pageSize).Offset(offset).Joins("LEFT JOIN posts_tags on posts_tags.post_id = posts.id").Joins("LEFT JOIN tags on posts_tags.tag_id = tags.id").Find(&posts).Error; err != nil {
		err = fmt.Errorf("find all posts: %w", err)
		return
	}
	if len(posts) == 0 {
		err = fmt.Errorf("find all posts: %w", entity.ErrPostNotFound)
		return
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

func (r *PostRepository) Count(condition string, params []interface{}) (count int, err error) {
	var count64 int64
	if err = r.db.Model(&entity.Post{}).Distinct("posts.id").Where(condition, params...).Joins("LEFT JOIN posts_tags on posts_tags.post_id = posts.id").Joins("LEFT JOIN tags on posts_tags.tag_id = tags.id").Count(&count64).Error; err != nil {
		err = fmt.Errorf("find all posts: %w", err)
		return
	}
	// int64を溢れることは運用的にないのでキャストしてしまう
	count = int(count64)
	return
}
