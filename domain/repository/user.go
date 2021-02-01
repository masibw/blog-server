package repository

import "github.com/masibw/blog-server/domain/entity"

type User interface {
	FindByID(id string) (*entity.User, error)
	FindAll(offset, pageSize int, condition string, params []interface{}) ([]*entity.User, error)
	FindByMailAddress(mailAddress string) (*entity.User, error)
	Create(user *entity.User) error
	UpdateLastLoggedinAt(user *entity.User) error
	DeleteByMailAddress(id string) error
}
