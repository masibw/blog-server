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

func TestTagHandler_StoreTag(t *testing.T) {
	tests := []struct {
		name                 string
		prepareMockTagRepoFn func(mock *mock_repository.MockTag)
		body                 string
		wantCode             int
	}{
		{
			name: "正常にタグを保存できる",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByName(gomock.Any()).Return(nil, entity.ErrTagNotFound)
				mock.EXPECT().Store(gomock.Any()).Return(nil)
			},
			body: `{
				"name" : "new_tag"
			}`,
			wantCode: http.StatusCreated,
		},
		{
			name: "tagDTOが満たされない時はStatusBadRequestエラーが返る",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
			},
			body:     "",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "保存に失敗した時はStatusInternalServerErrorエラーが返る",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByName("new_tag").Return(nil, nil)
				mock.EXPECT().Store(gomock.Any()).Return(errors.New("dummy error"))
			},
			body: `{
				"name" : "new_tag"
			}`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Repositoryのモック
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			tagUC := usecase.NewTagUseCase(mr)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body := bytes.NewBufferString(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/tags", body)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &TagHandler{
				tagUC: tagUC,
			}
			p.StoreTag(c)
			if w.Code != tt.wantCode {
				t.Errorf("StoreTag() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestTagHandler_GetTags(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	existsTags := []*entity.Tag{{
		ID:        "abcdefghijklmnopqrstuvwxyz",
		Name:      "new_tag",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}, {
		ID:        "abcdefghijklmnopqrstuvwxy2",
		Name:      "new_tag2",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}}

	tests := []struct {
		name                 string
		prepareMockTagRepoFn func(mock *mock_repository.MockTag)
		params               []struct {
			name  string
			value string
		}
		isDraft  string
		wantCode int
	}{
		{
			name: "正常にタグを取得できる",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(existsTags, nil)
			},
			params:   nil,
			wantCode: http.StatusOK,
		},
		{
			name: "タグが0件の時はhttp.StatusNotFoundを返す",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, entity.ErrTagNotFound)
			},
			params:   nil,
			wantCode: http.StatusNotFound,
		},
		{
			name: "タグの取得に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("dummy error"))
			},
			params:   nil,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "ページングを指定した時も正しく取得できる",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(existsTags[1:], nil)
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"page",
					"2",
				}, {
					"page_size",
					"1",
				},
			},
			wantCode: http.StatusOK,
		}, {
			name: "pageにint型に変換できない値が入っていた場合はStatusBadRequestを返す",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
			},
			params: []struct {
				name  string
				value string
			}{
				{
					"page",
					"can't_parse",
				}, {
					"page_size",
					"0",
				}, {
					"is_draft",
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
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			tagUC := usecase.NewTagUseCase(mr)

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

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tags"+queryParam, nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &TagHandler{
				tagUC: tagUC,
			}
			p.GetTags(c)
			if w.Code != tt.wantCode {
				t.Errorf("GetTags() code = %d, want = %d", w.Code, tt.wantCode)
			}

		})
	}
}

func TestTagHandler_GetTag(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	existsTag := &entity.Tag{
		ID:        "abcdefghijklmnopqrstuvwxyz",
		Name:      "new_tag",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}

	tests := []struct {
		name                 string
		prepareMockTagRepoFn func(mock *mock_repository.MockTag)
		ID                   string
		wantCode             int
	}{
		{
			name: "正常にタグを取得できる",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByID(gomock.Any()).Return(existsTag, nil)
			},
			ID:       "abcdefghijklmnopqrstuvwxyz",
			wantCode: http.StatusOK,
		},
		{
			name: "タグがない場合はStatusNotFoundを返す",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByID(gomock.Any()).Return(nil, entity.ErrTagNotFound)
			},
			ID:       "not_found",
			wantCode: http.StatusNotFound,
		},
		{
			name: "タグの取得に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
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
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			tagUC := usecase.NewTagUseCase(mr)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tags/"+tt.ID, nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &TagHandler{
				tagUC: tagUC,
			}
			p.GetTag(c)
			if w.Code != tt.wantCode {
				t.Errorf("GetTag() code = %d, want = %d", w.Code, tt.wantCode)
			}

		})
	}
}

func TestTagHandler_DeleteTag(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name                 string
		prepareMockTagRepoFn func(mock *mock_repository.MockTag)
		ID                   string
		wantCode             int
	}{
		{
			name: "正常にタグを削除できる",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			ID:       "abcdefghijklmnopqrstuvwxyz",
			wantCode: http.StatusOK,
		},
		{
			name: "タグがない場合はStatusNotFoundを返す",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().Delete(gomock.Any()).Return(entity.ErrTagNotFound)
			},
			ID:       "not_found",
			wantCode: http.StatusNotFound,
		},
		{
			name: "タグの削除に失敗した場合はStatusInternalServerErrorエラーが返る",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
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
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			tagUC := usecase.NewTagUseCase(mr)

			// HTTPRequestをテストするために必要な部分
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/tags/"+tt.ID, nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			p := &TagHandler{
				tagUC: tagUC,
			}
			p.DeleteTag(c)
			if w.Code != tt.wantCode {
				t.Errorf("GetTag() code = %d, want = %d", w.Code, tt.wantCode)
			}

		})
	}
}
