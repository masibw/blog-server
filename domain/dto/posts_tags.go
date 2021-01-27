package dto

import "time"

type PostsTagsDTO struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId" binding:"required"`
	TagID     string    `json:"tagId" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
