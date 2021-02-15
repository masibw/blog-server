package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/masibw/blog-server/domain/service"

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
				mock.EXPECT().Create(gomock.Any()).Return(nil)
			},
			body:     ``,
			wantCode: http.StatusCreated,
		},
		{
			name: "保存に失敗した時はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().Create(gomock.Any()).Return(errors.New("dummy error"))
			},
			body:     ``,
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
				t.Errorf("CreatePost() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestPostHandler_UpdatePost(t *testing.T) {
	tests := []struct {
		name              string
		prepareMockRepoFn func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags)
		ID                string
		body              string
		wantCode          int
	}{
		{
			name: "正常に投稿を更新できる",
			prepareMockRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
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
				mockPosts.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mockPosts.EXPECT().Update(gomock.Any()).Return(nil)
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{}, nil)
				mockPT.EXPECT().DeleteByPostID(gomock.Any()).Return(nil)
				mockTags.EXPECT().FindByName("a").Return(&entity.Tag{
					ID:        "abcdefghijklmnopqrstuvwxy2",
					Name:      "new_tag",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				}, nil)
				mockTags.EXPECT().FindByName("b").Return(&entity.Tag{
					ID:        "abcdefghijklmnopqrstuvwxy3",
					Name:      "new_tag2",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				}, nil)
				mockPT.EXPECT().Store(gomock.AssignableToTypeOf([]*entity.PostsTags{})).Return(nil)

			},
			ID: "abcdefghijklmnopqrstuvwxyz",
			body: `{
				"post": {
					"id": "abcdefghijklmnopqrstuvwxyz",
					"title": "new_post",
					"thumbnailUrl": "new_thumbnail_url",
					"content": "new_content",
					"permalink": "new_permalink",
					"isDraft": false,
					"createdAt": "2021-01-24T17:49:01+09:00",
					"updatedAt": "2021-01-27T14:48:55+09:00",
					"publishedAt": "0001-01-01T00:00:00Z"
				},
				"tags": [
				"a",
				"b"
			]
			}`,
			wantCode: http.StatusOK,
		},
		{
			name: "関連するタグがなくても正常に投稿を更新できる",
			prepareMockRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
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
				mockPosts.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mockPosts.EXPECT().Update(gomock.Any()).Return(nil)
			},
			ID: "abcdefghijklmnopqrstuvwxyz",
			body: `{
				"post": {
					"id": "abcdefghijklmnopqrstuvwxyz",
					"title": "new_post",
					"thumbnailUrl": "new_thumbnail_url",
					"content": "new_content",
					"permalink": "new_permalink",
					"isDraft": false,
					"createdAt": "2021-01-24T17:49:01+09:00",
					"updatedAt": "2021-01-27T14:48:55+09:00",
					"publishedAt": "0001-01-01T00:00:00Z"
				},
				"tags": []
			}`,
			wantCode: http.StatusOK,
		}, {
			name: "NULLが許容されるフィールドが空でも正常に投稿を更新できる",
			prepareMockRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
					ID:           "abcdefghijklmnopqrstuvwxyz",
					Title:        "",
					ThumbnailURL: "new_thumbnail_url",
					Content:      "",
					Permalink:    "",
					IsDraft:      true,
					CreatedAt:    flextime.Now(),
					UpdatedAt:    flextime.Now(),
					PublishedAt:  time.Time{},
				}, nil)
				mockPosts.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mockPosts.EXPECT().Update(gomock.Any()).Return(nil)
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{}, nil)
				mockPT.EXPECT().DeleteByPostID(gomock.Any()).Return(nil)
				mockTags.EXPECT().FindByName("a").Return(&entity.Tag{
					ID:        "abcdefghijklmnopqrstuvwxy2",
					Name:      "new_tag",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				}, nil)
				mockTags.EXPECT().FindByName("b").Return(&entity.Tag{
					ID:        "abcdefghijklmnopqrstuvwxy3",
					Name:      "new_tag2",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				}, nil)
				mockPT.EXPECT().Store(gomock.AssignableToTypeOf([]*entity.PostsTags{})).Return(nil)

			},
			ID: "abcdefghijklmnopqrstuvwxyz",
			body: `{
				"post": {
					"id": "abcdefghijklmnopqrstuvwxyz",
					"title": "",
					"thumbnailUrl": "new_thumbnail_url",
					"content": "",
					"permalink": "",
					"isDraft": true,
					"createdAt": "2021-01-24T17:49:01+09:00",
					"updatedAt": "2021-01-27T14:48:55+09:00",
					"publishedAt": "0001-01-01T00:00:00Z"
				},
				"tags": [
				"a",
				"b"
			]
			}`,
			wantCode: http.StatusOK,
		},
		{
			name: "更新に失敗した時はStatusInternalServerErrorエラーが返る",
			prepareMockRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
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
				mockPosts.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mockPosts.EXPECT().Update(gomock.Any()).Return(errors.New("dummy error"))
			},
			ID: "abcdefghijklmnopqrstuvwxyz",
			body: `{
				"post": {
					"id": "abcdefghijklmnopqrstuvwxyz",
					"title": "new_post",
					"thumbnailUrl": "new_thumbnail_url",
					"content": "new_content",
					"permalink": "new_permalink",
					"isDraft": true,
					"createdAt": "2021-01-24T17:49:01+09:00",
					"updatedAt": "2021-01-27T14:48:55+09:00",
					"publishedAt": "0001-01-01T00:00:00Z"
				},
				"tags": [
				"a",
				"b"
			]
			}`,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "bodyのbindに失敗した時はStatusBadRequestエラーが返る",
			prepareMockRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
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
			mP := mock_repository.NewMockPost(ctrl)
			mT := mock_repository.NewMockTag(ctrl)
			mPT := mock_repository.NewMockPostsTags(ctrl)
			tt.prepareMockRepoFn(mT, mP, mPT)

			pTS := service.NewPostsTagsService(mPT, mP, mT)
			postUC := usecase.NewPostUseCase(mP)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body := bytes.NewBufferString(tt.body)
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/posts/"+tt.ID, body)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &PostHandler{
				postUC:           postUC,
				postsTagsService: pTS,
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
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts, nil)
				mock.EXPECT().Count(gomock.Any(), gomock.Any()).Return(len(existsPosts), nil)
			},
			params:   nil,
			wantCode: http.StatusOK,
		},
		{
			name: "投稿が0件の時はhttp.StatusNotFoundを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, entity.ErrPostNotFound)
			},
			params:   nil,
			wantCode: http.StatusNotFound,
		},
		{
			name: "投稿の取得に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("dummy error"))
			},
			params:   nil,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "ページングを指定した時も正しく取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts[1:], nil)
				mock.EXPECT().Count(gomock.Any(), gomock.Any()).Return(len(existsPosts), nil)
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
		},
		{
			name: "tagを指定した時も正しく取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts[1:], nil)
				mock.EXPECT().Count(gomock.Any(), gomock.Any()).Return(len(existsPosts), nil)
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"tag",
					"a",
				},
			},
			wantCode: http.StatusOK,
		}, {
			name: "is-draftを指定した時も正しく取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts[1:], nil)
				mock.EXPECT().Count(gomock.Any(), gomock.Any()).Return(len(existsPosts), nil)
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
		{
			name: "sortを指定した時も正しく取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				sort.Slice(existsPosts, func(i, j int) bool {
					return existsPosts[i].CreatedAt.After(existsPosts[j].CreatedAt)
				})
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), "created_at desc").Return(existsPosts, nil)
				mock.EXPECT().Count(gomock.Any(), gomock.Any()).Return(len(existsPosts), nil)
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"sort",
					"-createdAt",
				},
			},
			wantCode: http.StatusOK,
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
		permalink             string
		params                []struct {
			name  string
			value string
		}
		wantCode int
	}{
		{
			name: "正常に投稿を取得できる",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(existsPost, nil)
			},
			permalink: "new_permalink",
			wantCode:  http.StatusOK,
		},
		{
			name: "投稿がない場合はStatusNotFoundを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
			},
			permalink: "not_found",
			wantCode:  http.StatusNotFound,
		},
		{
			name: "投稿の取得に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, errors.New("dummy error"))
			},
			permalink: "not_found",
			wantCode:  http.StatusInternalServerError,
		}, {
			name: "is-markdownにboolに変換できない値が入っていた場合はStatusBadRequestを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
			},
			params: []struct {
				name  string
				value string
			}{{
				"is-markdown",
				"can't_parse",
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
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/posts/"+tt.permalink+queryParam, nil)
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
