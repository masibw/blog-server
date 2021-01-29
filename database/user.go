package database

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/masibw/blog-server/domain/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id string) (*entity.User, error) {
	user := &entity.User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("find user: %w", entity.ErrUserNotFound)
			return nil, err
		}
		err = fmt.Errorf("find user: %w", err)
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByMailAddress(mailAddress string) (*entity.User, error) {
	user := &entity.User{}
	if err := r.db.Where("mail_address = ?", mailAddress).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user: %w", entity.ErrUserNotFound)
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Create(user *entity.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("create user: %w", entity.ErrUserAlreadyExisted)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) Update(user *entity.User) error {

	if err := r.db.Select("*").Updates(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindAll(offset, pageSize int, condition string, params []interface{}) (users []*entity.User, err error) {
	if err = r.db.Where(condition, params...).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		err = fmt.Errorf("find all users: %w", err)
		return
	}
	if len(users) == 0 {
		err = fmt.Errorf("find all users: %w", entity.ErrUserNotFound)
		return
	}
	return
}

func (r *UserRepository) DeleteByMailAddress(mailAddress string) error {
	result := r.db.Where("mail_address = ?", mailAddress).Delete(&entity.User{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("delete user: %w", entity.ErrUserNotFound)
	}
	if err := result.Error; err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
