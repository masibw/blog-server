package entity

import (
	"time"

	"github.com/Songmu/flextime"
	"github.com/masibw/blog-server/domain/dto"
	"github.com/masibw/blog-server/util"
)

type PostsTags struct {
	ID        string `gorm:"PRIMARY_KEY"`
	PostID    string
	TagID     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPostsTags(postID, tagID string) *PostsTags {
	return &PostsTags{
		ID:     util.Generate(flextime.Now()),
		PostID: postID,
		TagID:  tagID,
	}
}

func (p *PostsTags) ConvertToDTO() *dto.PostsTagsDTO {
	return &dto.PostsTagsDTO{
		ID:        p.ID,
		PostID:    p.PostID,
		TagID:     p.TagID,
		UpdatedAt: p.UpdatedAt,
		CreatedAt: p.CreatedAt,
	}
}
