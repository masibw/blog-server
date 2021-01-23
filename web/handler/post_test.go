package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Songmu/flextime"

	"github.com/masibw/blog-server/domain/entity"

	"github.com/golang/mock/gomock"
	"github.com/masibw/blog-server/domain/mock_repository"

	"github.com/gin-gonic/gin"
	"github.com/masibw/blog-server/usecase"
)

func TestPostHandler_StorePost(t *testing.T) {
	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		body                  string
		wantCode              int
	}{
		{
			name: "正常に投稿を保存できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Store(gomock.Any()).Return(nil)
			},
			body: `{
				"title" : "new_post",
				"thumbnail_url" : "new_thumbnail_url",
				"content" : "new_content",
				"permalink" : "new_permalink",
				"is_draft" : false
			}`,
			wantCode: http.StatusCreated,
		},
		{
			name: "postDTOが満たされない時はStatusBadRequestエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
			},
			body:     "",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "保存に失敗した時はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink("new_permalink").Return(&entity.Post{}, nil)
			},
			body: `{
				"title" : "new_post",
				"thumbnail_url" : "new_thumbnail_url",
				"content" : "new_content",
				"permalink" : "new_permalink",
				"is_draft" : false
			}`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Repositoryのモック
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPost(ctrl)
			tt.prepareMockPostRepoFn(mr)
			postUC := usecase.NewPostUseCase(mr)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body := bytes.NewBufferString(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/posts", body)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &PostHandler{
				postUC: postUC,
			}
			p.StorePost(c)
			if w.Code != tt.wantCode {
				t.Errorf("StorePost() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestPostHandler_GetPosts(t *testing.T) {

	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, time.UTC))
	defer flextime.Restore()

	existsPosts := []*entity.Post{{
		ID:           "abcdefghijklmnopqrstuvwxyz",
		Title:        "new_post",
		ThumbnailURL: "new_thumbnail_url",
		Content:      "new_content",
		Permalink:    "new_permalink",
		IsDraft:      false,
		CreatedAt:    flextime.Now(),
		UpdatedAt:    flextime.Now(),
		PublishedAt:  flextime.Now(),
	}, {
		ID:           "abcdefghijklmnopqrstuvwxy2",
		Title:        "new_post",
		ThumbnailURL: "new_thumbnail_url",
		Content:      "new_content",
		Permalink:    "new_permalink2",
		IsDraft:      false,
		CreatedAt:    flextime.Now(),
		UpdatedAt:    flextime.Now(),
		PublishedAt:  flextime.Now(),
	}}

	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		wantCode              int
	}{
		{
			name: "正常に投稿を取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll().Return(existsPosts, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "投稿が0件でもエラーにならない",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll().Return(existsPosts, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "投稿の取得に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll().Return(nil, errors.New("dummy error"))
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Repositoryのモック
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPost(ctrl)
			tt.prepareMockPostRepoFn(mr)
			postUC := usecase.NewPostUseCase(mr)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/posts", nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &PostHandler{
				postUC: postUC,
			}
			p.GetPosts(c)
			if w.Code != tt.wantCode {
				t.Errorf("GetPosts() code = %d, want = %d", w.Code, tt.wantCode)
			}

		})
	}
}
