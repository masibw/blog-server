package dto

import "time"

type TagDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
