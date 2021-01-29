package entity

import (
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/crypto/bcrypt"

	"github.com/Songmu/flextime"
	"github.com/masibw/blog-server/domain/dto"
	"github.com/masibw/blog-server/util"
)

type User struct {
	ID             string
	MailAddress    string
	Password       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastLoggedinAt time.Time
}

func NewUser(mailAddress, password string) (*User, error) {
	user := &User{
		ID:          util.Generate(flextime.Now()),
		MailAddress: mailAddress,
	}
	if len(password) > 72 {
		return nil, ErrPasswordTooLong
	}

	byteHash, err := bcrypt.GenerateFromPassword(*(*[]byte)(unsafe.Pointer(&password)), 12)
	if err != nil {
		return nil, fmt.Errorf("new user crypt error :%w", err)
	}

	user.Password = *(*string)(unsafe.Pointer(&byteHash))
	return user, nil
}

func (u *User) ConvertToDTO() *dto.UserDTO {
	return &dto.UserDTO{
		ID:          u.ID,
		MailAddress: u.MailAddress,
		Password:    u.Password,
		UpdatedAt:   u.UpdatedAt,
		CreatedAt:   u.CreatedAt,
	}
}
