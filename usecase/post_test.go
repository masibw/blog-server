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

func TestPostUseCase_StorePost(t *testing.T) { // nolint:gocognit

	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, time.UTC))
	defer flextime.Restore()

	tests := []struct {
		name                  string
		postDTO               *dto.PostDTO
		prepareMockPostRepoFn func(mock *mock_repository.MockPost)
		wantErr               error
	}{
		{
			name: "新規の投稿を保存し、その投稿を返す",
			postDTO: &dto.PostDTO{
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := true; return &b }(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Store(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "IsDraftがfalseであればPublishedAtが登録されていること",
			postDTO: &dto.PostDTO{
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := false; return &b }(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink(gomock.Any()).Return(nil, entity.ErrPostNotFound)
				mock.EXPECT().Store(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "パーマリンクが登録済みの場合ErrPermalinkAlreadyExistedエラーを返す",
			postDTO: &dto.PostDTO{
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      func() *bool { b := true; return &b }(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindByPermalink("new_permalink").Return(&entity.Post{}, nil)
				mock.EXPECT().Store(gomock.Any()).AnyTimes().Return(nil)
			},
			wantErr: entity.ErrPermalinkAlreadyExisted,
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

			got, err := p.StorePost(tt.postDTO)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StorePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if errors.Is(err, entity.ErrPermalinkAlreadyExisted) {
				if got == nil {
					return
				}
				t.Errorf("StorePost() got = %v, want = nil", got)
			}

			if got.ID == "" {
				t.Errorf("StorePost() ID nil want UULD")
			}
			if got.CreatedAt.Unix() == 0 || got.UpdatedAt.Unix() == 0 {
				t.Errorf("StorePost() time.Time field did not filled with value")
			}

			if got.ThumbnailURL != tt.postDTO.ThumbnailURL {
				t.Errorf("StorePost() ThumbnailURL does not match got: %v, want: %v", got, tt.postDTO)
			}

			if got.Title != tt.postDTO.Title {
				t.Errorf("StorePost() Title does not match got: %v, want: %v", got, tt.postDTO)
			}
			if got.Content != tt.postDTO.Content {
				t.Errorf("StorePost() Content does not match got: %v, want: %v", got, tt.postDTO)
			}
			if got.Permalink != tt.postDTO.Permalink {
				t.Errorf("StorePost() Permalink does not match got: %v, want: %v", got, tt.postDTO)
			}

			if *got.IsDraft != *tt.postDTO.IsDraft {
				t.Errorf("StorePost() IsDraft does not match got: %v, want: %v", got, tt.postDTO)
			}

			if !*tt.postDTO.IsDraft && got.PublishedAt.IsZero() {
				t.Errorf("StorePost() PublishedAt does not set got: %v, want: %v", got.PublishedAt, flextime.Now())
			}

		})
	}
}

func TestPostUseCase_GetPosts(t *testing.T) {

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
		want                  []*dto.PostDTO
		wantErr               bool
	}{
		{
			name: "postDTOsを返すこと",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPost) {
				mock.EXPECT().FindAll().Return(existsPosts, nil)
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
				mock.EXPECT().FindAll().Return(nil, errors.New("dummy error"))
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

			got, err := p.GetPosts()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetPosts() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}
