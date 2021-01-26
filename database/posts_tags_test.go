package database

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/Songmu/flextime"

	"github.com/masibw/blog-server/domain/entity"
)

func TestPostsTagsRepository_FindByPostIDAndTagName(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	if err := tx.Create(&entity.Post{
		ID:           "abcdefghijklmnopqrstuvwxy1",
		Title:        "new_postTags",
		ThumbnailURL: "new_thumbnail_url",
		Content:      "new_content",
		Permalink:    "new_permalink",
		IsDraft:      false,
		CreatedAt:    flextime.Now(),
		UpdatedAt:    flextime.Now(),
		PublishedAt:  flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	if err := tx.Create(&entity.Tag{
		ID:        "abcdefghijklmnopqrstuvwxy2",
		Name:      "new_tag",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	if err := tx.Create(&entity.PostsTags{
		ID:        "abcdefghijklmnopqrstuvwxy3",
		PostID:    "abcdefghijklmnopqrstuvwxy1",
		TagID:     "abcdefghijklmnopqrstuvwxy2",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		postID  string
		tagName string
		want    *entity.PostsTags
		wantErr error
	}{
		{
			name:    "存在する投稿とタグの関連を正常に取得できる",
			postID:  "abcdefghijklmnopqrstuvwxy1",
			tagName: "new_tag",
			want: &entity.PostsTags{
				ID:        "abcdefghijklmnopqrstuvwxy3",
				PostID:    "abcdefghijklmnopqrstuvwxy1",
				TagID:     "abcdefghijklmnopqrstuvwxy2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			wantErr: nil,
		},
		{
			name:    "存在しないpostIDの場合ErrPostsTagsNotFoundを返す",
			postID:  "not_found",
			tagName: "new_tag",
			want:    nil,
			wantErr: entity.ErrPostsTagsNotFound,
		},
		{
			name:    "存在しないtagNameの場合ErrPostsTagsNotFoundを返す",
			postID:  "abcdefghijklmnopqrstuvwxy1",
			tagName: "not_found",
			want:    nil,
			wantErr: entity.ErrPostsTagsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostsTagsRepository{db: tx}
			got, err := r.FindByPostIDAndTagName(tt.postID, tt.tagName)
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

func TestPostsTagsRepository_Store(t *testing.T) {
	tx := db.Begin()

	if err := tx.Create(&entity.Post{
		ID:           "abcdefghijklmnopqrstuvwxy1",
		Title:        "new_postTags",
		ThumbnailURL: "new_thumbnail_url",
		Content:      "new_content",
		Permalink:    "new_permalink",
		IsDraft:      false,
		CreatedAt:    flextime.Now(),
		UpdatedAt:    flextime.Now(),
		PublishedAt:  flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	if err := tx.Create(&entity.Tag{
		ID:        "abcdefghijklmnopqrstuvwxy2",
		Name:      "new_tag",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		postTags *entity.PostsTags
		wantErr  error
	}{
		{
			name: "新規の投稿とタグの関連を正常に保存できる",
			postTags: &entity.PostsTags{
				ID:        "abcdefghijklmnopqrstuvwxy3",
				PostID:    "abcdefghijklmnopqrstuvwxy1",
				TagID:     "abcdefghijklmnopqrstuvwxy2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			wantErr: nil,
		},
		{
			name: "既に存在するIDの場合ErrPostsTagsAlreadyExistedエラーを返す",
			postTags: &entity.PostsTags{
				ID:        "abcdefghijklmnopqrstuvwxy3",
				PostID:    "abcdefghijklmnopqrstuvwxy1",
				TagID:     "abcdefghijklmnopqrstuvwxy2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			wantErr: entity.ErrPostsTagsAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostsTagsRepository{db: tx}
			if err := r.Store(tt.postTags); !errors.Is(err, tt.wantErr) {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tx.Rollback()
}

func TestPostsTagsRepository_Delete(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	if err := tx.Create(&entity.Post{
		ID:           "abcdefghijklmnopqrstuvwxy1",
		Title:        "new_postTags",
		ThumbnailURL: "new_thumbnail_url",
		Content:      "new_content",
		Permalink:    "new_permalink",
		IsDraft:      false,
		CreatedAt:    flextime.Now(),
		UpdatedAt:    flextime.Now(),
		PublishedAt:  flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	if err := tx.Create(&entity.Tag{
		ID:        "abcdefghijklmnopqrstuvwxy2",
		Name:      "new_tag",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		ID             string
		existPostsTags *entity.PostsTags
		want           *entity.PostsTags
		wantErr        error
	}{
		{
			name: "存在する投稿とタグの関連を正常に削除できる",
			ID:   "abcdefghijklmnopqrstuvwxy3",
			existPostsTags: &entity.PostsTags{
				ID:        "abcdefghijklmnopqrstuvwxy3",
				PostID:    "abcdefghijklmnopqrstuvwxy1",
				TagID:     "abcdefghijklmnopqrstuvwxy2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name:           "存在しないIDの場合ErrPostsTagsNotFoundを返す",
			ID:             "not_found",
			existPostsTags: nil,
			want:           nil,
			wantErr:        entity.ErrPostsTagsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existPostsTags != nil {
				if err := tx.Create(tt.existPostsTags).Error; err != nil {
					t.Fatal(err)
				}
			}

			r := &PostsTagsRepository{db: tx}
			err := r.Delete(tt.ID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			//TODO 削除したことを確かめるテスト

		})
	}

	tx.Rollback()
}
