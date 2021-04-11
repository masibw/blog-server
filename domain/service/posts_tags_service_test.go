package service

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/Songmu/flextime"

	"github.com/golang/mock/gomock"
	"github.com/masibw/blog-server/domain/mock_repository"

	"github.com/masibw/blog-server/domain/entity"
)

// getTagEntityを内部的に呼んでいるが非公開なメソッドなので同時にテストする
func TestLinkPostTags(t *testing.T) { // nolint:gocognit

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	t.Cleanup(func() {
		flextime.Restore()
	})

	tests := []struct {
		name                 string
		postID               string
		tagNames             []string
		prepareMockTagRepoFn func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags)
		wantTagsLen          int
		wantErr              error
	}{
		{
			name:     "新規に投稿と存在するタグの関連を保存できる",
			postID:   "abcdefghijklmnopqrstuvwxy1",
			tagNames: []string{"a", "b"},
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
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
			wantTagsLen: 2,
			wantErr:     nil,
		},
		{
			name:     "タグが存在しなければ作成して投稿とタグの関連を保存する",
			postID:   "abcdefghijklmnopqrstuvwxy1",
			tagNames: []string{"a", "b"},
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{}, nil)
				mockPT.EXPECT().DeleteByPostID(gomock.Any()).Return(nil)
				mockTags.EXPECT().FindByName("a").Return(nil, entity.ErrTagNotFound)
				mockTags.EXPECT().Store(gomock.Any()).Return(nil)
				mockTags.EXPECT().FindByName("b").Return(&entity.Tag{
					ID:        "abcdefghijklmnopqrstuvwxy3",
					Name:      "new_tag2",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				}, nil)
				mockPT.EXPECT().Store(gomock.AssignableToTypeOf([]*entity.PostsTags{})).Return(nil)
			},
			wantTagsLen: 2,
			wantErr:     nil,
		},
		{
			name:     "重複したタグがあればUniqueな分のみ作成される",
			postID:   "abcdefghijklmnopqrstuvwxy1",
			tagNames: []string{"a", "a"},
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(&entity.Post{}, nil)
				mockPT.EXPECT().DeleteByPostID(gomock.Any()).Return(nil)
				mockTags.EXPECT().FindByName("a").Return(&entity.Tag{
					ID:        "abcdefghijklmnopqrstuvwxy2",
					Name:      "new_tag",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				}, nil)
				mockPT.EXPECT().Store(gomock.AssignableToTypeOf([]*entity.PostsTags{})).Return(nil)
			},
			wantTagsLen: 1,
			wantErr:     nil,
		},
		{
			name:     "投稿が存在しなければErrPostNotFoundエラーを返す",
			postID:   "",
			tagNames: []string{"a", "b"},
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(nil, entity.ErrPostNotFound)
			},
			wantTagsLen: 0,
			wantErr:     entity.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mP := mock_repository.NewMockPost(ctrl)
			mT := mock_repository.NewMockTag(ctrl)
			mPT := mock_repository.NewMockPostsTags(ctrl)
			tt.prepareMockTagRepoFn(mT, mP, mPT)

			p := &PostsTagsService{
				postsTagsRepository: mPT,
				postRepository:      mP,
				tagRepository:       mT,
			}

			got, err := p.LinkPostTags(tt.postID, tt.tagNames)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StoreTag() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.wantTagsLen, len(got)); diff != "" {
				t.Errorf("ConvertContentToHTML() mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
