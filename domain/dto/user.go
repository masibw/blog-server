package dto

import "time"

type UserDTO struct {
	ID             string
	MailAddress    string
	Password       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastLoggedinAt time.Time
}
