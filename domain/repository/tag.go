package repository

import "github.com/masibw/blog-server/domain/entity"

type Tag interface {
	FindByID(id string) (*entity.Tag, error)
	FindAll(offset, pageSize int, condition string, params []interface{}) ([]*entity.Tag, error)
	FindByName(name string) (*entity.Tag, error)
	Store(tag *entity.Tag) error
	Delete(id string) error
	Count() (int, error)
}
