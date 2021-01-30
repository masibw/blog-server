package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/Songmu/flextime"
	"github.com/masibw/blog-server/domain/dto"

	"github.com/golang/mock/gomock"
	"github.com/masibw/blog-server/domain/entity"
	"github.com/masibw/blog-server/domain/mock_repository"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/usecase"
)

func TestAuthMiddleware_Authenticate(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name                  string
		prepareMockUserRepoFn func(mock *mock_repository.MockUser)
		body                  string
		want                  interface{}
		wantCode              int
		wantErr               error
	}{
		{
			name: "正常に認証できる",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().FindByMailAddress("test@example.com").Return(&entity.User{
					ID:          "abcdefghijklmnopqrstuvwxyz",
					MailAddress: "test@example.com",
					Password:    "$2a$12$MdZRSm..1nFoRkBUqb1SE.Epo8J34q1rGDZkT/vv0.VNgDViQNQPi",
					CreatedAt:   flextime.Now(),
					UpdatedAt:   flextime.Now(),
				}, nil)
			},
			body: `{
			  "mailAddress":"test@example.com",
			  "password":"test"
			}`,
			want: &dto.UserDTO{
				ID:          "abcdefghijklmnopqrstuvwxyz",
				MailAddress: "test@example.com",
				Password:    "$2a$12$MdZRSm..1nFoRkBUqb1SE.Epo8J34q1rGDZkT/vv0.VNgDViQNQPi",
				CreatedAt:   flextime.Now(),
				UpdatedAt:   flextime.Now(),
			},
			wantCode: http.StatusCreated,
			wantErr:  nil,
		},
		{
			name: "loginが満たされない時はjwt.ErrMissingLoginValuesエラーが返る",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
			},
			body:     "",
			want:     nil,
			wantCode: http.StatusUnauthorized,
			wantErr:  jwt.ErrMissingLoginValues,
		},
		{
			name: "ユーザー取得に失敗した時はjwt.ErrFailedAuthenticationエラーが返る",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().FindByMailAddress(gomock.Any()).Return(nil, entity.ErrUserNotFound)
			},
			body: `{
			  "mailAddress":"test@example.com",
			  "password":"test"
			}`,
			want:     nil,
			wantCode: http.StatusUnauthorized,
			wantErr:  jwt.ErrFailedAuthentication,
		},
		{
			name: "認証に失敗した場合はjwt.ErrFailedAuthenticationエラーが返る",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUser) {
				mock.EXPECT().FindByMailAddress("test@example.com").Return(&entity.User{
					ID:          "abcdefghijklmnopqrstuvwxyz",
					MailAddress: "test@example.com",
					Password:    "not_found",
					CreatedAt:   flextime.Now(),
					UpdatedAt:   flextime.Now(),
				}, nil)
			},
			body: `{
			  "mailAddress":"test@example.com",
			  "password":"test"
			}`,
			want:     nil,
			wantCode: http.StatusUnauthorized,
			wantErr:  jwt.ErrFailedAuthentication,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Repositoryのモック
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockUser(ctrl)
			tt.prepareMockUserRepoFn(mr)
			userUC := usecase.NewUserUseCase(mr)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body := bytes.NewBufferString(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", body)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			a := &AuthMiddleware{
				userUC: userUC,
			}
			got, err := a.Authenticate(c)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				gotUserDTO, ok := got.(*dto.UserDTO)
				if !ok {
					t.Errorf("Authenticate() return not *dto.UserDTO got= \n%v", got)
				}
				if diff := cmp.Diff(tt.want, gotUserDTO); diff != "" {
					t.Errorf("Authenticate() mismatch (-want +got):\n%s", diff)
				}
			}

		})
	}
}

func TestAuthMiddleware_Authorize(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()
	tests := []struct {
		name string
		data interface{}
		want bool
	}{
		{
			name: "正常に認可できる",
			data: &dto.UserDTO{
				ID:          "abcdefghijklmnopqrstuvwxyz",
				MailAddress: "test@example.com",
				Password:    "$2a$12$MdZRSm..1nFoRkBUqb1SE.Epo8J34q1rGDZkT/vv0.VNgDViQNQPi",
				CreatedAt:   flextime.Now(),
				UpdatedAt:   flextime.Now(),
			},
			want: true,
		},
		{
			name: "*userDTO以外の型であればfalseを返す",
			data: &entity.User{
				ID:          "abcdefghijklmnopqrstuvwxyz",
				MailAddress: "test@example.com",
				Password:    "$2a$12$MdZRSm..1nFoRkBUqb1SE.Epo8J34q1rGDZkT/vv0.VNgDViQNQPi",
				CreatedAt:   flextime.Now(),
				UpdatedAt:   flextime.Now(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Repositoryのモック
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockUser(ctrl)
			userUC := usecase.NewUserUseCase(mr)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			a := &AuthMiddleware{
				userUC: userUC,
			}
			got := a.Authorize(tt.data, c)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Authenticate() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
