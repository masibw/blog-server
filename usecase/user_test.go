package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/Songmu/flextime"

	"github.com/golang/mock/gomock"
	"github.com/masibw/blog-server/domain/mock_repository"

	"github.com/masibw/blog-server/domain/entity"

	"github.com/masibw/blog-server/domain/dto"
)

func TestUserUseCase_StoreUser(t *testing.T) { // nolint:gocognit

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	t.Cleanup(func() {
		flextime.Restore()
	})

	tests := []struct {
		name                  string
		userDTO               *dto.UserDTO
		prepareMockUserRepoFn func(mock *mock_repository.MockUser)
		wantErr               error
	}{
		{
			name: "新規のユーザーを保存し、そのユーザーを返す",
			userDTO: &dto.UserDTO{
				MailAddress: "new_user@example.com",
				Password:    "new_password",
			},
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().FindByMailAddress(gomock.Any()).Return(nil, entity.ErrUserNotFound)
				mock.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "メールアドレスが登録済みの場合ErrUserMailAddressAlreadyExistedエラーを返す",
			userDTO: &dto.UserDTO{
				MailAddress: "new_user@example.com",
				Password:    "new_password",
			},
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().FindByMailAddress("new_user@example.com").Return(&entity.User{}, nil)
				mock.EXPECT().Create(gomock.Any()).AnyTimes().Return(nil)
			},
			wantErr: entity.ErrUserMailAddressAlreadyExisted,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockUser(ctrl)
			tt.prepareMockUserRepoFn(mr)
			u := &UserUseCase{
				userRepository: mr,
			}

			err := u.StoreUser(tt.userDTO)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StoreUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUseCase_GetUserByMailAddress(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	t.Cleanup(func() {
		flextime.Restore()
	})

	existsUser := &entity.User{
		ID:          "abcdefghijklmnopqrstuvwxyz",
		MailAddress: "new_user@example.com",
		Password:    "new_password",
		CreatedAt:   flextime.Now(),
		UpdatedAt:   flextime.Now(),
	}

	tests := []struct {
		name                  string
		prepareMockUserRepoFn func(mock *mock_repository.MockUser)
		mailAddress           string
		want                  *dto.UserDTO
		wantErr               bool
	}{
		{
			name: "userDTOを返すこと",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().FindByMailAddress(gomock.Any()).Return(existsUser, nil)
			},
			want: &dto.UserDTO{
				ID:          "abcdefghijklmnopqrstuvwxyz",
				MailAddress: "new_user@example.com",
				Password:    "new_password",
				CreatedAt:   flextime.Now(),
				UpdatedAt:   flextime.Now(),
			},
			mailAddress: "new_user@example.com",
			wantErr:     false,
		},
		{
			name: "FindByMailAddressがエラーを返した時はuserDTOが空であること",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().FindByMailAddress("not_found").Return(nil, entity.ErrUserNotFound)
			},
			mailAddress: "not_found",
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockUser(ctrl)
			tt.prepareMockUserRepoFn(mr)
			u := &UserUseCase{
				userRepository: mr,
			}

			got, err := u.GetUserByMailAddress(tt.mailAddress)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUser() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestUserUseCase_UpdateLastLoggedinAt(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	t.Cleanup(func() {
		flextime.Restore()
	})

	tests := []struct {
		name                  string
		prepareMockUserRepoFn func(mock *mock_repository.MockUser)
		userDTO               *dto.UserDTO
		wantErr               bool
	}{
		{
			name: "更新に成功した場合はエラーを返さないこと",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().UpdateLastLoggedinAt(gomock.Any()).Return(nil)
			},
			userDTO: &dto.UserDTO{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "test@example.com",
				Password:       "$2a$12$MdZRSm..1nFoRkBUqb1SE.Epo8J34q1rGDZkT/vv0.VNgDViQNQPi",
				CreatedAt:      flextime.Now().Add(-time.Second),
				UpdatedAt:      flextime.Now().Add(-time.Second),
				LastLoggedinAt: flextime.Now(),
			},
			wantErr: false,
		},
		{
			name: "UpdateLastLoggedinAtがエラーを返した時はエラーを返すこと",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().UpdateLastLoggedinAt(gomock.Any()).Return(entity.ErrUserNotFound)
			},
			userDTO: &dto.UserDTO{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "test@example.com",
				Password:       "$2a$12$MdZRSm..1nFoRkBUqb1SE.Epo8J34q1rGDZkT/vv0.VNgDViQNQPi",
				CreatedAt:      flextime.Now().Add(-time.Second),
				UpdatedAt:      flextime.Now().Add(-time.Second),
				LastLoggedinAt: flextime.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockUser(ctrl)
			tt.prepareMockUserRepoFn(mr)
			u := &UserUseCase{
				userRepository: mr,
			}

			err := u.UpdateLastLoggedinAt(tt.userDTO)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUseCase_DeleteUserByMailAddress(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	t.Cleanup(func() {
		flextime.Restore()
	})

	tests := []struct {
		name                  string
		prepareMockUserRepoFn func(mock *mock_repository.MockUser)
		mailAddress           string
		wantErr               bool
	}{
		{
			name: "削除に成功した場合はエラーを返さないこと",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().DeleteByMailAddress(gomock.Any()).Return(nil)
			},
			mailAddress: "new_user@example.com",
			wantErr:     false,
		},
		{
			name: "Deleteがエラーを返した時はエラーを返すこと",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().DeleteByMailAddress("not_found").Return(entity.ErrUserNotFound)
			},
			mailAddress: "not_found",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockUser(ctrl)
			tt.prepareMockUserRepoFn(mr)
			u := &UserUseCase{
				userRepository: mr,
			}

			err := u.DeleteUserByMailAddress(tt.mailAddress)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
