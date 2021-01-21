package entity

import (
	"time"

	"github.com/masibw/blog-server/domain/dto"

	"github.com/masibw/blog-server/util"
)

type Post struct {
	ID           string
	Title        string
	ThumbnailURL string
	Content      string
	Permalink    string
	IsDraft      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewPost(thumbnailURL string, title string, content string, permalink string, isDraft bool) *Post {
	return &Post{
		ID:           util.GenUULD(),
		Title:        title,
		ThumbnailURL: thumbnailURL,
		Content:      content,
		Permalink:    permalink,
		IsDraft:      isDraft,
	}
}

func (p *Post) ConvertToDTO() *dto.PostDTO {
	return &dto.PostDTO{
		ID:           p.ID,
		Title:        p.Title,
		ThumbnailURL: p.ThumbnailURL,
		Content:      p.Content,
		Permalink:    p.Permalink,
		IsDraft:      &p.IsDraft,
		UpdatedAt:    p.UpdatedAt,
		CreatedAt:    p.CreatedAt,
	}
}
