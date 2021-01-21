package repository

import "github.com/masibw/blog-server/domain/entity"

type Post interface {
	Store(post *entity.Post) error
}
