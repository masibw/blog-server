package dto

import "time"

type TagDTO struct {
	ID        string
	Name      string `json:"name" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
