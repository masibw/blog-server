package usecase

import (
	"errors"
	"fmt"

	"github.com/masibw/blog-server/domain/dto"
	"github.com/masibw/blog-server/domain/entity"
	"github.com/masibw/blog-server/domain/repository"
)

type UserUseCase struct {
	userRepository repository.User
}

func NewUserUseCase(userRepository repository.User) *UserUseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (p *UserUseCase) StoreUser(userDTO *dto.UserDTO) error {
	var user *entity.User
	var err error

	user, err = p.userRepository.FindByMailAddress(userDTO.MailAddress)
	if err != nil && !errors.Is(err, entity.ErrUserNotFound) {
		return fmt.Errorf("store user mailAddress=%v: %w", userDTO.MailAddress, err)
	}
	if user != nil {
		return fmt.Errorf("store user mailAddress=%v: %w", userDTO.MailAddress, entity.ErrUserMailAddressAlreadyExisted)
	}

	user, err = entity.NewUser(userDTO.MailAddress, userDTO.Password)

	if err != nil {
		return fmt.Errorf("store user mailAddress=%v: %w", userDTO.MailAddress, err)
	}

	err = p.userRepository.Create(user)
	if err != nil {
		return fmt.Errorf("store user mailAddress=%v: %w", userDTO.MailAddress, err)
	}

	return nil
}

func (p *UserUseCase) GetUserByMailAddress(mailAddress string) (userDTO *dto.UserDTO, err error) {
	var user *entity.User
	user, err = p.userRepository.FindByMailAddress(mailAddress)
	if err != nil {
		err = fmt.Errorf("get user: %w", err)
		return
	}
	userDTO = user.ConvertToDTO()
	return
}

func (p *UserUseCase) UpdateLastLoggedinAt(userDTO *dto.UserDTO) (err error) {
	user := &entity.User{
		ID:             userDTO.ID,
		LastLoggedinAt: userDTO.LastLoggedinAt}
	err = p.userRepository.UpdateLastLoggedinAt(user)
	if err != nil {
		err = fmt.Errorf("update user: %w", err)
		return
	}
	return
}

func (p *UserUseCase) DeleteUserByMailAddress(mailAddress string) (err error) {
	err = p.userRepository.DeleteByMailAddress(mailAddress)
	if err != nil {
		err = fmt.Errorf("delete user: %w", err)
		return
	}
	return nil
}
