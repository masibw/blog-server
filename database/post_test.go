package database

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Songmu/flextime"

	"github.com/masibw/blog-server/domain/entity"

	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db = NewTestDB()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestPostRepository_Store(t *testing.T) {
	tx := db.Begin()

	tests := []struct {
		name    string
		post    *entity.Post
		wantErr error
	}{
		{
			name: "新規の投稿を正常に保存できる",
			post: &entity.Post{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      false,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
				PublishedAt:  time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "既に存在するIDの場合ErrPostAlreadyExistedエラーを返す",
			post: &entity.Post{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink_2",
				IsDraft:      false,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
				PublishedAt:  time.Time{},
			},
			wantErr: entity.ErrPostAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostRepository{db: tx}
			if err := r.Store(tt.post); !errors.Is(err, tt.wantErr) {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tx.Rollback()
}

func TestPostRepository_FindByPermalink(t *testing.T) {
	tx := db.Begin()
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, time.UTC))
	defer flextime.Restore()

	if err := tx.Create(&entity.Post{
		ID:           "abcdefghijklmnopqrstuvwxyz",
		Title:        "new_post",
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

	tests := []struct {
		name      string
		permalink string
		want      *entity.Post
		wantErr   error
	}{
		{
			name:      "存在する投稿を正常に取得できる",
			permalink: "new_permalink",
			want: &entity.Post{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			},
			wantErr: nil,
		},
		{
			name:      "存在しないpermalinkの場合ErrPostNotFoundを返す",
			permalink: "permalink_not_found",
			want:      nil,
			wantErr:   entity.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostRepository{db: tx}
			got, err := r.FindByPermalink(tt.permalink)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByPermalink() got = %v, want %v", got, tt.want)
			}
		})
	}

	tx.Rollback()
}
