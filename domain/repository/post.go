package repository

import "github.com/masibw/blog-server/domain/entity"

type Post interface {
	FindByPermalink(permalink string) (*entity.Post, error)
	Store(post *entity.Post) error
}
