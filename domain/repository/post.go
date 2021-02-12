package repository

import "github.com/masibw/blog-server/domain/entity"

type Post interface {
	FindByID(id string) (*entity.Post, error)
	FindAll(offset, pageSize int, condition string, params []interface{}, sortCondition string) ([]*entity.Post, error)
	FindByPermalink(permalink string) (*entity.Post, error)
	Create(post *entity.Post) error
	Update(post *entity.Post) error
	Delete(id string) error
	Count(condition string, params []interface{}) (int, error)
}
