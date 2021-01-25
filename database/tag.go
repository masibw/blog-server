package database

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/masibw/blog-server/domain/entity"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) FindByID(id string) (*entity.Tag, error) {
	tag := &entity.Tag{}
	if err := r.db.Where("id = ?", id).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("find tag: %w", entity.ErrTagNotFound)
			return nil, err
		}
		err = fmt.Errorf("find tag: %w", err)
		return nil, err
	}
	return tag, nil
}

func (r *TagRepository) FindByName(name string) (*entity.Tag, error) {
	tag := &entity.Tag{}
	if err := r.db.Where("name = ?", name).First(tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find tag: %w", entity.ErrTagNotFound)
		}
		return nil, fmt.Errorf("find tag: %w", err)
	}
	return tag, nil
}

func (r *TagRepository) Store(tag *entity.Tag) error {
	if err := r.db.Create(tag).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("create tag: %w", entity.ErrTagAlreadyExisted)
		}
		return fmt.Errorf("create tag: %w", err)
	}
	return nil
}

func (r *TagRepository) FindAll(offset, pageSize int) (tags []*entity.Tag, err error) {
	if err = r.db.Offset(offset).Limit(pageSize).Find(&tags).Error; err != nil {
		err = fmt.Errorf("find all tags: %w", err)
		return
	}
	if len(tags) == 0 {
		err = fmt.Errorf("find all tags: %w", entity.ErrTagNotFound)
		return
	}
	return
}

func (r *TagRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&entity.Tag{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("delete tag: %w", entity.ErrTagNotFound)
	}
	if err := result.Error; err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return nil
}
