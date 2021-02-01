package database

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/Songmu/flextime"

	"github.com/masibw/blog-server/domain/entity"
)

func TestUserRepository_FindByID(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	if err := tx.Create(&entity.User{
		ID:             "abcdefghijklmnopqrstuvwxyz",
		MailAddress:    "new_mailAddress",
		Password:       "new_password",
		CreatedAt:      flextime.Now(),
		UpdatedAt:      flextime.Now(),
		LastLoggedinAt: flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		ID      string
		want    *entity.User
		wantErr error
	}{
		{
			name: "存在するユーザーを正常に取得できる",
			ID:   "abcdefghijklmnopqrstuvwxyz",
			want: &entity.User{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			},
			wantErr: nil,
		},
		{
			name:    "存在しないIDの場合ErrUserNotFoundを返す",
			ID:      "not_found",
			want:    nil,
			wantErr: entity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{db: tx}
			got, err := r.FindByID(tt.ID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	tx.Rollback()
}

func TestUserRepository_UpdateLastLoggedinAt(t *testing.T) {
	tx := db.Begin()

	existUser := &entity.User{
		ID:             "abcdefghijklmnopqrstuvwxyz",
		MailAddress:    "new_mailAddress",
		CreatedAt:      flextime.Now(),
		UpdatedAt:      flextime.Now(),
		LastLoggedinAt: flextime.Now(),
	}

	if err := tx.Create(existUser).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		user    *entity.User
		want    *entity.User
		wantErr error
	}{
		{
			name: "ユーザーのLastLoggedinAtだけを正常に更新できる",
			user: &entity.User{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
				LastLoggedinAt: time.Time{},
			},
			want: &entity.User{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "new_mailAddress",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: time.Time{},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{db: tx}
			if err := r.UpdateLastLoggedinAt(tt.user); !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateLastLoggedinAt() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := &entity.User{ID: tt.want.ID}
			err := tx.First(&got).Error
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateApproxTime(time.Second)); diff != "" {
				t.Errorf("Authenticate() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	tx.Rollback()
}
func TestUserRepository_Store(t *testing.T) {
	tx := db.Begin()

	tests := []struct {
		name    string
		user    *entity.User
		wantErr error
	}{
		{
			name: "新規のユーザーを正常に保存できる",
			user: &entity.User{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
				LastLoggedinAt: time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "既に存在するIDの場合ErrUserAlreadyExistedエラーを返す",
			user: &entity.User{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "new_mailAddress_2",
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
				LastLoggedinAt: time.Time{},
			},
			wantErr: entity.ErrUserAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{db: tx}
			if err := r.Create(tt.user); !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tx.Rollback()
}

func TestUserRepository_FindByMailAddress(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	if err := tx.Create(&entity.User{
		ID:             "abcdefghijklmnopqrstuvwxyz",
		MailAddress:    "new_mailAddress",
		Password:       "new_password",
		CreatedAt:      flextime.Now(),
		UpdatedAt:      flextime.Now(),
		LastLoggedinAt: flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		mailAddress string
		want        *entity.User
		wantErr     error
	}{
		{
			name:        "存在するユーザーを正常に取得できる",
			mailAddress: "new_mailAddress",
			want: &entity.User{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			},
			wantErr: nil,
		},
		{
			name:        "存在しないmailAddressの場合ErrUserNotFoundを返す",
			mailAddress: "mailAddress_not_found",
			want:        nil,
			wantErr:     entity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{db: tx}
			got, err := r.FindByMailAddress(tt.mailAddress)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindByMailAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindByMailAddress() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	tx.Rollback()
}

func TestUserRepository_FindAll(t *testing.T) { // nolint:gocognit
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name       string
		existUsers []*entity.User
		offset     int
		pageSize   int
		condition  string
		params     []interface{}
		want       []*entity.User
		wantErr    error
	}{
		{
			name: "存在するユーザーを正常に全件取得できる",
			existUsers: []*entity.User{{
				ID:             "abcdefghijklmnopqrstuvwxy1",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}, {
				ID:             "abcdefghijklmnopqrstuvwxy2",
				MailAddress:    "new_mailAddress2",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}, {
				ID:             "abcdefghijklmnopqrstuvwxy3",
				MailAddress:    "new_mailAddress3",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}},
			offset:    0,
			pageSize:  0,
			condition: "",
			params:    []interface{}{},
			want: []*entity.User{{
				ID:             "abcdefghijklmnopqrstuvwxy1",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}, {
				ID:             "abcdefghijklmnopqrstuvwxy2",
				MailAddress:    "new_mailAddress2",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}, {
				ID:             "abcdefghijklmnopqrstuvwxy3",
				MailAddress:    "new_mailAddress3",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}},
			wantErr: nil,
		}, {
			name: "ページネーションを適用して取得できる",
			existUsers: []*entity.User{{
				ID:             "abcdefghijklmnopqrstuvwxy1",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}, {
				ID:             "abcdefghijklmnopqrstuvwxy2",
				MailAddress:    "new_mailAddress2",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}, {
				ID:             "abcdefghijklmnopqrstuvwxy3",
				MailAddress:    "new_mailAddress3",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}},
			offset:    1,
			pageSize:  2,
			condition: "",
			params:    []interface{}{},
			want: []*entity.User{{
				ID:             "abcdefghijklmnopqrstuvwxy2",
				MailAddress:    "new_mailAddress2",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}, {
				ID:             "abcdefghijklmnopqrstuvwxy3",
				MailAddress:    "new_mailAddress3",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			}},
			wantErr: nil,
		},
		{
			name:       "ユーザーが存在しない場合はErrUserNotFoundを返す",
			existUsers: nil,
			want:       []*entity.User{},
			wantErr:    entity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existUsers != nil {
				if err := tx.Debug().Create(tt.existUsers).Error; err != nil {
					t.Fatal(err)
				}
			}
			var users []*entity.User
			tx.Find(&users)
			for _, v := range users {
				t.Logf("%v", v)
			}
			r := &UserRepository{db: tx.Debug()}
			got, err := r.FindAll(tt.offset, tt.pageSize, tt.condition, tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindAll()  error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindByID() mismatch (-want +got):\n%s", diff)
			}

			if tt.existUsers != nil {
				if err := tx.Delete(tt.existUsers).Error; err != nil {
					t.Fatal(err)
				}
			}
		})
	}

	tx.Rollback()
}

func TestUserRepository_DeleteByMailAddress(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name        string
		mailAddress string
		existUser   *entity.User
		want        *entity.User
		wantErr     error
	}{
		{
			name:        "存在するユーザーを正常に削除できる",
			mailAddress: "new_mailAddress",
			existUser: &entity.User{
				ID:             "abcdefghijklmnopqrstuvwxyz",
				MailAddress:    "new_mailAddress",
				Password:       "new_password",
				CreatedAt:      flextime.Now(),
				UpdatedAt:      flextime.Now(),
				LastLoggedinAt: flextime.Now(),
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name:        "存在しないメールアドレスの場合ErrUserNotFoundを返す",
			mailAddress: "not_found",
			existUser:   nil,
			want:        nil,
			wantErr:     entity.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existUser != nil {
				if err := tx.Create(tt.existUser).Error; err != nil {
					t.Fatal(err)
				}
			}

			r := &UserRepository{db: tx}
			err := r.DeleteByMailAddress(tt.mailAddress)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			//TODO 削除したことを確かめるテスト

		})
	}

	tx.Rollback()
}
