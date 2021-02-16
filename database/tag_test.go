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

func TestTagRepository_FindByID(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	if err := tx.Create(&entity.Tag{
		ID:        "abcdefghijklmnopqrstuvwxyz",
		Name:      "new_tag",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		ID      string
		want    *entity.Tag
		wantErr error
	}{
		{
			name: "存在するタグを正常に取得できる",
			ID:   "abcdefghijklmnopqrstuvwxyz",
			want: &entity.Tag{
				ID:        "abcdefghijklmnopqrstuvwxyz",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			wantErr: nil,
		},
		{
			name:    "存在しないIDの場合ErrTagNotFoundを返す",
			ID:      "not_found",
			want:    nil,
			wantErr: entity.ErrTagNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TagRepository{db: tx}
			got, err := r.FindByID(tt.ID)
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

func TestTagRepository_Store(t *testing.T) {
	tx := db.Begin()

	tests := []struct {
		name    string
		tag     *entity.Tag
		wantErr error
	}{
		{
			name: "新規のタグを正常に保存できる",
			tag: &entity.Tag{
				ID:        "abcdefghijklmnopqrstuvwxyz",
				Name:      "new_tag",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			wantErr: nil,
		},
		{
			name: "既に存在するIDの場合ErrTagAlreadyExistedエラーを返す",
			tag: &entity.Tag{
				ID:        "abcdefghijklmnopqrstuvwxyz",
				Name:      "new_tag",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			wantErr: entity.ErrTagAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TagRepository{db: tx}
			if err := r.Store(tt.tag); !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tx.Rollback()
}

func TestTagRepository_FindAll(t *testing.T) { // nolint:gocognit
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name           string
		existTags      []*entity.Tag
		existPosts     []*entity.Post
		existPostsTags *entity.PostsTags
		offset         int
		pageSize       int
		condition      string
		params         []interface{}
		want           []*entity.Tag
		wantErr        error
	}{
		{
			name: "存在するタグを正常に全件取得できる",
			existTags: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy1",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy3",
				Name:      "new_tag3",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			offset:    0,
			pageSize:  0,
			condition: "",
			params:    []interface{}{},
			want: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy1",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy3",
				Name:      "new_tag3",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			wantErr: nil,
		}, {
			name: "ページネーションを適用して取得できる",
			existTags: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy1",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy3",
				Name:      "new_tag3",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			offset:    1,
			pageSize:  2,
			condition: "",
			params:    []interface{}{},
			want: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy3",
				Name:      "new_tag3",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			wantErr: nil,
		},
		{
			name: "postフィルターを適用して取得できる",
			existTags: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy1",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy3",
				Name:      "new_tag3",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			existPostsTags: &entity.PostsTags{
				ID:        "abcdefghijklmnopqrstuvwxy2",
				PostID:    "abcdefghijklmnopqrstuvwxy5",
				TagID:     "abcdefghijklmnopqrstuvwxy2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			existPosts: []*entity.Post{{
				ID:           "abcdefghijklmnopqrstuvwxy5",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			offset:    0,
			pageSize:  0,
			condition: "posts.id = ?",
			params:    []interface{}{"abcdefghijklmnopqrstuvwxy5"},
			want: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			wantErr: nil,
		}, {
			name:      "タグが存在しない場合はErrTagNotFoundを返す",
			existTags: nil,
			want:      []*entity.Tag{},
			wantErr:   entity.ErrTagNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existTags != nil {
				if err := tx.Create(tt.existTags).Error; err != nil {
					t.Fatal(err)
				}
			}
			if tt.existPosts != nil {
				if err := tx.Create(tt.existPosts).Error; err != nil {
					t.Fatal(err)
				}
			}
			if tt.existPostsTags != nil {
				if err := tx.Create(tt.existPostsTags).Error; err != nil {
					t.Fatal(err)
				}
			}

			r := &TagRepository{db: tx}
			got, err := r.FindAll(tt.offset, tt.pageSize, tt.condition, tt.params)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindAll()  error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindByID() mismatch (-want +got):\n%s", diff)
			}
			if tt.existPostsTags != nil {
				if err := tx.Delete(tt.existPostsTags).Error; err != nil {
					t.Fatal(err)
				}
			}

			if tt.existTags != nil {
				if err := tx.Delete(tt.existTags).Error; err != nil {
					t.Fatal(err)
				}
			}
			if tt.existPosts != nil {
				if err := tx.Delete(tt.existPosts).Error; err != nil {
					t.Fatal(err)
				}
			}

		})
	}

	tx.Rollback()
}

func TestTagRepository_Delete(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name     string
		ID       string
		existTag *entity.Tag
		want     *entity.Tag
		wantErr  error
	}{
		{
			name: "存在するタグを正常に削除できる",
			ID:   "abcdefghijklmnopqrstuvwxyz",
			existTag: &entity.Tag{
				ID:        "abcdefghijklmnopqrstuvwxyz",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name:     "存在しないIDの場合ErrTagNotFoundを返す",
			ID:       "not_found",
			existTag: nil,
			want:     nil,
			wantErr:  entity.ErrTagNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existTag != nil {
				if err := tx.Create(tt.existTag).Error; err != nil {
					t.Fatal(err)
				}
			}

			r := &TagRepository{db: tx}
			err := r.Delete(tt.ID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			//TODO 削除したことを確かめるテスト

		})
	}

	tx.Rollback()
}

func TestTagRepository_Count(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name      string
		existTags []*entity.Tag
		condition string
		params    []interface{}
		want      int
		wantErr   error
	}{
		{
			name: "存在するタグを正常に全件取得できる",
			existTags: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy1",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy3",
				Name:      "new_tag3",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			condition: "",
			params:    []interface{}{},
			want:      3,
			wantErr:   nil,
		}, {
			name: "ページネーションを適用して取得できる",
			existTags: []*entity.Tag{{
				ID:        "abcdefghijklmnopqrstuvwxy1",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}, {
				ID:        "abcdefghijklmnopqrstuvwxy3",
				Name:      "new_tag3",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			}},
			condition: "",
			params:    []interface{}{},
			want:      3,
			wantErr:   nil,
		}, {
			name:      "タグが存在しない場合は0を返す",
			existTags: nil,
			want:      0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existTags != nil {
				if err := tx.Create(tt.existTags).Error; err != nil {
					t.Fatal(err)
				}
			}

			r := &TagRepository{db: tx}
			got, err := r.Count()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Count()  error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Count() mismatch (-want +got):\n%s", diff)
			}

			if tt.existTags != nil {
				if err := tx.Delete(tt.existTags).Error; err != nil {
					t.Fatal(err)
				}
			}
		})
	}

	tx.Rollback()
}
