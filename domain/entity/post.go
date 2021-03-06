package entity

import (
	"time"
	"unsafe"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	"github.com/masibw/blog-server/constant"
	"github.com/masibw/blog-server/util"

	"github.com/Songmu/flextime"
	"github.com/masibw/blog-server/domain/dto"
)

type Post struct {
	ID           string `gorm:"PRIMARY_KEY"`
	Title        string
	ThumbnailURL string
	Content      string
	Permalink    string
	IsDraft      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PublishedAt  time.Time
}

func NewPost() *Post {
	return &Post{
		ID:           util.Generate(flextime.Now()),
		Title:        "",
		ThumbnailURL: constant.DefaultThumbnailURL,
		Content:      "",
		Permalink:    "",
		IsDraft:      true,
		PublishedAt:  time.Time{},
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
		PublishedAt:  p.PublishedAt,
	}
}

func (p *Post) ConvertFromDTO(postDTO *dto.PostDTO) {
	p.ID = postDTO.ID
	p.Title = postDTO.Title
	p.ThumbnailURL = postDTO.ThumbnailURL
	p.Content = postDTO.Content
	p.Permalink = postDTO.Permalink
	p.IsDraft = *postDTO.IsDraft
	p.UpdatedAt = postDTO.UpdatedAt
	p.CreatedAt = postDTO.CreatedAt
	p.PublishedAt = postDTO.PublishedAt
}

func (p *Post) ConvertContentToHTML() {
	bytesContent := *(*[]byte)(unsafe.Pointer(&p.Content))
	unsafeHTML := blackfriday.Run(bytesContent)
	sanitizedHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)
	content := *(*string)(unsafe.Pointer(&sanitizedHTML))
	p.Content = content
}
