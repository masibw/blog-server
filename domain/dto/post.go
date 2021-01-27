package dto

import (
	"time"
)

type PostDTO struct {
	ID           string    `json:"id"`
	Title        string    `json:"title" binding:"required"`
	ThumbnailURL string    `json:"thumbnailUrl" binding:"required"`
	Content      string    `json:"content" binding:"required"`
	Permalink    string    `json:"permalink" binding:"required"`
	IsDraft      *bool     `json:"isDraft" binding:"required"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	PublishedAt  time.Time `json:"publishedAt"`
}
