package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/masibw/blog-server/constant"

	"github.com/google/go-cmp/cmp"

	"github.com/Songmu/flextime"

	"github.com/golang/mock/gomock"
	"github.com/masibw/blog-server/domain/mock_repository"

	"github.com/masibw/blog-server/domain/entity"

	"github.com/masibw/blog-server/domain/dto"
)

func TestPostUseCase_StorePost(t *testing.T) { // nolint:gocognit

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		wantErr               error
	}{
		{
			name: "新規の投稿を保存し、その投稿を返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPost(ctrl)
			tt.prepareMockPostRepoFn(mr)
			p := &PostUseCase{
				postRepository: mr,
			}

			got, err := p.CreatePost()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if errors.Is(err, entity.ErrPermalinkAlreadyExisted) {
				if got == nil {
					return
				}
				t.Errorf("CreatePost() got = %v, want = nil", got)
			}

			if got.ID == "" {
				t.Errorf("CreatePost() ID nil want UULD")
			}
			if got.ThumbnailURL != constant.DefaultThumbnailURL {
				t.Errorf("CreatePost() ThumbnailURL is not default")
			}
			if got.CreatedAt.Unix() == 0 || got.UpdatedAt.Unix() == 0 {
				t.Errorf("CreatePost() time.Time field did not filled with value")
			}

		})
	}
}

func TestPostUseCase_UpdatePost(t *testing.T) { // nolint:gocognit

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name                  string
		postDTO               *dto.PostDTO
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		wantErr               error
	}{
		{
			name: "投稿を更新し、その投稿を返す",
			postDTO: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := true; return &b }(),
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  time.Time{},
			},
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
			wantErr: nil,
		},
		{
			name: "投稿が存在しなければErrPostNotFoundエラーを返す",
			postDTO: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := true; return &b }(),
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  time.Time{},
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(nil, entity.ErrPostNotFound)
			},
			wantErr: entity.ErrPostNotFound,
		},
		{
			name: "既に更新先のPermalinkが別の投稿で使用されている場合はErrPermalinkAlreadyExistedエラーを返す",
			postDTO: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := true; return &b }(),
			},
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
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(&entity.Post{
					ID:           "another_id",
					Title:        "new_post",
					ThumbnailURL: "new_thumbnail_url",
					Content:      "new_content2",
					Permalink:    "new_permalink",
					IsDraft:      false,
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
					PublishedAt:  time.Time{},
				}, nil)
			},
			wantErr: entity.ErrPermalinkAlreadyExisted,
		},
		{
			name: "publishedAtが初期値で，IsDraftがfalseであればPublishedAtが登録されていること",
			postDTO: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := false; return &b }(),
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  time.Time{},
			},
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
					PublishedAt:  time.Time{}}, nil)
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "publishedAtが初期値ではなくIsDraftがfalseであればPublishedAtが更新されないこと",
			postDTO: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := false; return &b }(),
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
					ID:           "abcdefghijklmnopqrstuvwxyz",
					Title:        "new_post",
					ThumbnailURL: "new_thumbnail_url",
					Content:      "new_content",
					Permalink:    "new_permalink",
					IsDraft:      false,
					CreatedAt:    flextime.Now(),
					UpdatedAt:    flextime.Now(),
					PublishedAt:  time.Time{},
				}, nil)
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "IsDraftがtrueでTitle,Content,Permalinkが空であれば，更新できる",
			postDTO: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "",
				Permalink:    "",
				IsDraft:      func() *bool { b := true; return &b }(),
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{
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
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "IsDraftがfalseでTitle,Content,Permalinkが空であれば entity.ErrPostHasEmptyFieldエラーを返す",
			postDTO: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "",
				Permalink:    "",
				IsDraft:      func() *bool { b := false; return &b }(),
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {},
			wantErr:               entity.ErrPostHasEmptyField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPost(ctrl)
			tt.prepareMockPostRepoFn(mr)
			p := &PostUseCase{
				postRepository: mr,
			}

			got, err := p.UpdatePost(tt.postDTO)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}
			if got.ID == "" {
				t.Errorf("UpdatePost() ID nil want UULD")
			}
			if got.CreatedAt.Unix() == 0 || got.UpdatedAt.Unix() == 0 {
				t.Errorf("UpdatePost() time.Time field did not filled with value")
			}

			if got.ThumbnailURL != tt.postDTO.ThumbnailURL {
				t.Errorf("UpdatePost() ThumbnailURL does not match got: %v, want: %v", got, tt.postDTO)
			}

			if got.Title != tt.postDTO.Title {
				t.Errorf("UpdatePost() Title does not match got: %v, want: %v", got, tt.postDTO)
			}
			if got.Content != tt.postDTO.Content {
				t.Errorf("UpdatePost() Content does not match got: %v, want: %v", got, tt.postDTO)
			}
			if got.Permalink != tt.postDTO.Permalink {
				t.Errorf("UpdatePost() Permalink does not match got: %v, want: %v", got, tt.postDTO)
			}

			if *got.IsDraft != *tt.postDTO.IsDraft {
				t.Errorf("UpdatePost() IsDraft does not match got: %v, want: %v", got, tt.postDTO)
			}

			if *tt.postDTO.IsDraft && got.PublishedAt != tt.postDTO.PublishedAt {
				t.Errorf("UpdatePost() PublishedAt updated: %v, want: %v", got.PublishedAt, flextime.Now())
			}

			if !*tt.postDTO.IsDraft && got.PublishedAt.IsZero() {
				t.Errorf("UpdatePost() PublishedAt does not set got: %v, want: %v", got.PublishedAt, flextime.Now())
			}

		})
	}
}

func TestPostUseCase_GetPosts(t *testing.T) {

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
		want                  []*dto.PostDTO
		wantErr               bool
	}{
		{
			name: "postDTOsを返すこと",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(existsPosts, nil)
			},
			want: []*dto.PostDTO{
				{
					ID:           "abcdefghijklmnopqrstuvwxyz",
					Title:        "new_post",
					ThumbnailURL: "new_thumbnail_url",
					Content:      "<p>new_content</p>\n",
					Permalink:    "new_permalink",
					IsDraft:      func() *bool { b := false; return &b }(),
					CreatedAt:    flextime.Now(),
					UpdatedAt:    flextime.Now(),
					PublishedAt:  flextime.Now(),
				},
				{
					ID:           "abcdefghijklmnopqrstuvwxy2",
					Title:        "new_post",
					ThumbnailURL: "new_thumbnail_url",
					Content:      "<p>new_content</p>\n",
					Permalink:    "new_permalink2",
					IsDraft:      func() *bool { b := false; return &b }(),
					CreatedAt:    flextime.Now(),
					UpdatedAt:    flextime.Now(),
					PublishedAt:  flextime.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "FindAllがエラーを返した時はpostDTOsが空であること",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("dummy error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPost(ctrl)
			tt.prepareMockPostRepoFn(mr)
			p := &PostUseCase{
				postRepository: mr,
			}

			// このGetPostsの責務はパラメータを受け取ってpostDTOsを返すだけなのでパラメータの中身はなんでも良い(はず)
			got, err := p.GetPosts(0, 0, "", []interface{}{}, "")

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetPosts() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestPostUseCase_GetPost(t *testing.T) {

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
		want                  *dto.PostDTO
		wantErr               bool
	}{
		{
			name: "postDTOを返すこと",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(existsPost, nil)
			},
			want: &dto.PostDTO{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "<p>new_content</p>\n",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := false; return &b }(),
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			},
			permalink: "new_permalink",
			wantErr:   false,
		},
		{
			name: "FindByPermalinkがエラーを返した時はpostDTOが空であること",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink("not_found").Return(nil, entity.ErrPostNotFound)
			},
			permalink: "not_found",
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPost(ctrl)
			tt.prepareMockPostRepoFn(mr)
			p := &PostUseCase{
				postRepository: mr,
			}

			got, err := p.GetPost(tt.permalink)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetPost() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestPostUseCase_DeletePost(t *testing.T) {

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
		wantErr               bool
	}{
		{
			name: "削除に成功した場合はエラーを返さないこと",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			ID:      "abcdefghijklmnopqrstuvwxyz",
			wantErr: false,
		},
		{
			name: "Deleteがエラーを返した時はエラーを返すこと",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().Delete("not_found").Return(entity.ErrPostNotFound)
			},
			ID:      "not_found",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPost(ctrl)
			tt.prepareMockPostRepoFn(mr)
			p := &PostUseCase{
				postRepository: mr,
			}

			err := p.DeletePost(tt.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
