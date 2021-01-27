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
				"thumbnailUrl" : "new_thumbnail_url",
				"content" : "new_content",
				"permalink" : "new_permalink",
				"isDraft" : false
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
				mock.EXPECT().FindByPermalink("new_permalink").Return(nil, nil)
				mock.EXPECT().Store(gomock.Any()).Return(errors.New("dummy error"))
			},
			body: `{
				"title" : "new_post",
				"thumbnailUrl" : "new_thumbnail_url",
				"content" : "new_content",
				"permalink" : "new_permalink",
				"isDraft" : false
			}`,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "保存に失敗した時はStatusBadRequestエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink("new_permalink").Return(&entity.Post{}, nil)
			},
			body: `{
				"title" : "new_post",
				"thumbnailUrl" : "new_thumbnail_url",
				"content" : "new_content",
				"permalink" : "new_permalink",
				"isDraft" : false
			}`,
			wantCode: http.StatusBadRequest,
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

func TestPostHandler_UpdatePost(t *testing.T) {
	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		ID                    string
		body                  string
		wantCode              int
	}{
		{
			name: "正常に投稿を更新できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
					ID:           "abcdefghijklmnopqrstuvwxyz",
					Title:        "new_post",
					ThumbnailURL: "new_thumbnail_url",
					Content:      "new_content",
					Permalink:    "new_permalink",
					IsDraft:      true,
					CreatedAt:    flextime.Now(),
					UpdatedAt:    flextime.Now(),
					PublishedAt:  time.Time{},
				}, nil)
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Update(gomock.Any()).Return(nil)
			},
			ID: "abcdefghijklmnopqrstuvwxyz",
			body: `{
				"ID" : "abcdefghijklmnopqrstuvwxyz",
				"title" : "new_post",
				"thumbnailUrl" : "new_thumbnail_url",
				"content" : "new_content",
				"permalink" : "new_permalink",
				"isDraft" : true,
				"createdAt": "2021-01-24T17:49:01+09:00",
				"updatedAt": "2021-01-27T14:48:55+09:00",
				"publishedAt": "0001-01-01T00:00:00Z"
			}`,
			wantCode: http.StatusOK,
		},
		{
			name: "postDTOが満たされない時はStatusBadRequestエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
			},
			body:     "",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "更新に失敗した時はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
					ID:           "abcdefghijklmnopqrstuvwxyz",
					Title:        "new_post",
					ThumbnailURL: "new_thumbnail_url",
					Content:      "new_content",
					Permalink:    "new_permalink",
					IsDraft:      true,
					CreatedAt:    flextime.Now(),
					UpdatedAt:    flextime.Now(),
					PublishedAt:  time.Time{},
				}, nil)
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Update(gomock.Any()).Return(errors.New("dummy error"))
			},
			ID: "abcdefghijklmnopqrstuvwxyz",
			body: `{
				"ID" : "abcdefghijklmnopqrstuvwxyz",
				"title" : "new_post",
				"thumbnailUrl" : "new_thumbnail_url",
				"content" : "new_content",
				"permalink" : "new_permalink",
				"isDraft" : true,
				"createdAt": "2021-01-24T17:49:01+09:00",
				"updatedAt": "2021-01-27T14:48:55+09:00",
				"publishedAt": "0001-01-01T00:00:00Z"
			}`,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "bodyのbindに失敗した時はStatusBadRequestエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
			},
			ID: "abcdefghijklmnopqrstuvwxyz",
			body: `{
			}`,
			wantCode: http.StatusBadRequest,
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
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/posts/"+tt.ID, body)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &PostHandler{
				postUC: postUC,
			}
			p.UpdatePost(c)
			if w.Code != tt.wantCode {
				t.Errorf("UpdatePost() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestPostHandler_GetPosts(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
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
		params                []struct {
			name  string
			value string
		}
		isDraft  string
		wantCode int
	}{
		{
			name: "正常に投稿を取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts, nil)
			},
			params:   nil,
			wantCode: http.StatusOK,
		},
		{
			name: "投稿が0件の時はhttp.StatusNotFoundを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, entity.ErrPostNotFound)
			},
			params:   nil,
			wantCode: http.StatusNotFound,
		},
		{
			name: "投稿の取得に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("dummy error"))
			},
			params:   nil,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "ページングを指定した時も正しく取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts[1:], nil)
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"page",
					"2",
				}, {
					"page-size",
					"1",
				},
			},
			wantCode: http.StatusOK,
		}, {
			name: "is-draftを指定した時も正しく取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts[1:], nil)
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"page",
					"0",
				}, {
					"page-size",
					"0",
				}, {
					"is-draft",
					"true",
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "is-draftにboolに変換できない値が入っていた場合はStatusBadRequestを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"page",
					"0",
				}, {
					"page-size",
					"0",
				}, {
					"is-draft",
					"can't_parse",
				},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "pageにintに変換できない値が入っていた場合はStatusBadRequestを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"page",
					"can't_parse",
				}, {
					"page-size",
					"0",
				}, {
					"is-draft",
					"false",
				},
			},
			wantCode: http.StatusBadRequest,
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

			var queryParam string
			for _, v := range tt.params {
				queryParam += "&" + v.name + "=" + v.value
			}
			if queryParam != "" {
				queryParam = "?" + queryParam[1:]
			}

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/posts"+queryParam, nil)
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

func TestPostHandler_GetPost(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	existsPost := &entity.Post{
		ID:           "abcdefghijklmnopqrstuvwxyz",
		Title:        "new_post",
		ThumbnailURL: "new_thumbnail_url",
		Content:      "new_content",
		Permalink:    "new_permalink",
		IsDraft:      false,
		CreatedAt:    flextime.Now(),
		UpdatedAt:    flextime.Now(),
		PublishedAt:  flextime.Now(),
	}

	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		ID                    string
		wantCode              int
	}{
		{
			name: "正常に投稿を取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(existsPost, nil)
			},
			ID:       "abcdefghijklmnopqrstuvwxyz",
			wantCode: http.StatusOK,
		},
		{
			name: "投稿がない場合はStatusNotFoundを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(nil, entity.ErrPostNotFound)
			},
			ID:       "not_found",
			wantCode: http.StatusNotFound,
		},
		{
			name: "投稿の取得に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(nil, errors.New("dummy error"))
			},
			ID:       "not_found",
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
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/posts/"+tt.ID, nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &PostHandler{
				postUC: postUC,
			}
			p.GetPost(c)
			if w.Code != tt.wantCode {
				t.Errorf("GetPost() code = %d, want = %d", w.Code, tt.wantCode)
			}

		})
	}
}

func TestPostHandler_DeletePost(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		ID                    string
		wantCode              int
	}{
		{
			name: "正常に投稿を削除できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			ID:       "abcdefghijklmnopqrstuvwxyz",
			wantCode: http.StatusOK,
		},
		{
			name: "投稿がない場合はStatusNotFoundを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().Delete(gomock.Any()).Return(entity.ErrPostNotFound)
			},
			ID:       "not_found",
			wantCode: http.StatusNotFound,
		},
		{
			name: "投稿の削除に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().Delete(gomock.Any()).Return(errors.New("dummy error"))
			},
			ID:       "not_found",
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
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/posts/"+tt.ID, nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &PostHandler{
				postUC: postUC,
			}
			p.DeletePost(c)
			if w.Code != tt.wantCode {
				t.Errorf("GetPost() code = %d, want = %d", w.Code, tt.wantCode)
			}

		})
	}
}
