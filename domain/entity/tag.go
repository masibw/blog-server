package entity

import (
	"time"

	"github.com/Songmu/flextime"
	"github.com/masibw/blog-server/domain/dto"
	"github.com/masibw/blog-server/util"
)

type Tag struct {
	ID        string `gorm:"PRIMARY_KEY"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTag(name string) *Tag {
	return &Tag{
		ID:   util.Generate(flextime.Now()),
		Name: name,
	}
}

func (p *Tag) ConvertToDTO() *dto.TagDTO {
	return &dto.TagDTO{
		ID:        p.ID,
		Name:      p.Name,
		UpdatedAt: p.UpdatedAt,
		CreatedAt: p.CreatedAt,
	}
}
