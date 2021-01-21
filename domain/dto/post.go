package dto

import "time"

type PostDTO struct {
	ID           string
	Title        string `json:"title" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url" binding:"required"`
	Content      string `json:"content" binding:"required"`
	Permalink    string `json:"permalink" binding:"required"`
	IsDraft      *bool  `json:"is_draft" binding:"required"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
