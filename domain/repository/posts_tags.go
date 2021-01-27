package repository

import "github.com/masibw/blog-server/domain/entity"

type PostsTags interface {
	FindByPostIDAndTagName(postID, tagName string) (*entity.PostsTags, error)
	Store(postsTags *entity.PostsTags) error
	Delete(id string) error
}
