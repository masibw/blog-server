package repository

import "github.com/masibw/blog-server/domain/entity"

type Post interface {
	FindByID(id string) (*entity.Post, error)
	FindAll(offset, pageSize int) ([]*entity.Post, error)
	FindByPermalink(permalink string) (*entity.Post, error)
	Store(post *entity.Post) error
	Delete(id string) error
}
