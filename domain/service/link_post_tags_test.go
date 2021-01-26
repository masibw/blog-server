package service

import (
	"errors"
	"testing"
	"time"

	"github.com/Songmu/flextime"

	"github.com/golang/mock/gomock"
	"github.com/masibw/blog-server/domain/mock_repository"

	"github.com/masibw/blog-server/domain/entity"
)

func TestLinkPostTags(t *testing.T) { // nolint:gocognit

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name                 string
		postID               string
		tagName              string
		prepareMockTagRepoFn func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags)
		wantErr              error
	}{
		{
			name:    "新規に投稿と存在するタグの関連を保存できる",
			postID:  "abcdefghijklmnopqrstuvwxy1",
			tagName: "new_tag",
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPT.EXPECT().FindByPostIDAndTagName(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(nil, nil)
				mockTags.EXPECT().FindByName(gomock.Any()).Return(&entity.Tag{
					ID:        "abcdefghijklmnopqrstuvwxy2",
					Name:      "new_tag",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				}, nil)
				mockPT.EXPECT().Store(gomock.AssignableToTypeOf(&entity.PostsTags{})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "タグが存在しなければ作成して投稿とタグの関連を保存する",
			postID:  "abcdefghijklmnopqrstuvwxy1",
			tagName: "new_tag_not_found",
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPT.EXPECT().FindByPostIDAndTagName(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(nil, nil)
				mockTags.EXPECT().FindByName(gomock.Any()).Return(nil, entity.ErrTagNotFound)
				mockTags.EXPECT().Store(gomock.Any()).Return(nil)
				mockPT.EXPECT().Store(gomock.AssignableToTypeOf(&entity.PostsTags{})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "既に同じpostIDとtagIDが紐づいていればErrPostsTagsAlreadyExistedエラーを返す",
			postID:  "",
			tagName: "",
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPT.EXPECT().FindByPostIDAndTagName(gomock.Any(), gomock.Any()).Return(nil, entity.ErrPostsTagsAlreadyExisted)
			},
			wantErr: entity.ErrPostsTagsAlreadyExisted,
		},
		{
			name:    "投稿が存在しなければErrPostNotFoundエラーを返す",
			postID:  "",
			tagName: "",
			prepareMockTagRepoFn: func(mockTags *mock_repository.MockTag, mockPosts *mock_repository.MockPost, mockPT *mock_repository.MockPostsTags) {
				mockPT.EXPECT().FindByPostIDAndTagName(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockPosts.EXPECT().FindByID(gomock.Any()).Return(nil, entity.ErrPostNotFound)
			},
			wantErr: entity.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mP := mock_repository.NewMockPost(ctrl)
			mT := mock_repository.NewMockTag(ctrl)
			mPT := mock_repository.NewMockPostsTags(ctrl)
			tt.prepareMockTagRepoFn(mT, mP, mPT)

			err := LinkPostTags(mPT, mP, mT, tt.postID, tt.tagName)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StoreTag() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
