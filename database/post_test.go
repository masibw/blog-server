package database

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/golang-migrate/migrate/v4"
	"github.com/masibw/blog-server/config"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/Songmu/flextime"

	"github.com/masibw/blog-server/domain/entity"

	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	var mig *migrate.Migrate
	mig, err = migrate.New("file://"+os.Getenv("MIGRATION_FILE"), "mysql://"+config.PureDSN())
	if err != nil {
		panic(err)
	}
	if err := mig.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}

	db = NewTestDB()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestPostRepository_FindByID(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
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
		name    string
		ID      string
		want    *entity.Post
		wantErr error
	}{
		{
			name: "存在する投稿を正常に取得できる",
			ID:   "abcdefghijklmnopqrstuvwxyz",
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
			name:    "存在しないIDの場合ErrPostNotFoundを返す",
			ID:      "not_found",
			want:    nil,
			wantErr: entity.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostRepository{db: tx}
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

func TestPostRepository_Update(t *testing.T) {
	tx := db.Begin()

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
		name    string
		post    *entity.Post
		wantErr error
	}{
		{
			name: "投稿を正常に更新できる",
			post: &entity.Post{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content2",
				Permalink:    "new_permalink",
				IsDraft:      false,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
				PublishedAt:  time.Time{},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostRepository{db: tx}
			if err := r.Update(tt.post); !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tx.Rollback()
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
			if err := r.Create(tt.post); !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tx.Rollback()
}

func TestPostRepository_FindByPermalink(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
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
				t.Errorf("FindByPermalink() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindByPermalink() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	tx.Rollback()
}

func TestPostRepository_FindAll(t *testing.T) { // nolint:gocognit
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name           string
		existPosts     []*entity.Post
		existTag       *entity.Tag
		existPostsTags *entity.PostsTags
		offset         int
		pageSize       int
		condition      string
		params         []interface{}
		want           []*entity.Post
		wantErr        error
	}{
		{
			name: "存在する投稿を正常に全件取得できる",
			existPosts: []*entity.Post{{
				ID:           "abcdefghijklmnopqrstuvwxy1",
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
			}, {
				ID:           "abcdefghijklmnopqrstuvwxy3",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink3",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			offset:    0,
			pageSize:  0,
			condition: "",
			params:    []interface{}{},
			want: []*entity.Post{{
				ID:           "abcdefghijklmnopqrstuvwxy1",
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
			}, {
				ID:           "abcdefghijklmnopqrstuvwxy3",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink3",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			wantErr: nil,
		}, {
			name: "ページネーションを適用して取得できる",
			existPosts: []*entity.Post{{
				ID:           "abcdefghijklmnopqrstuvwxy1",
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
			}, {
				ID:           "abcdefghijklmnopqrstuvwxy3",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink3",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			offset:    1,
			pageSize:  2,
			condition: "",
			params:    []interface{}{},
			want: []*entity.Post{{
				ID:           "abcdefghijklmnopqrstuvwxy2",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink2",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}, {
				ID:           "abcdefghijklmnopqrstuvwxy3",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink3",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			wantErr: nil,
		}, {
			name: "is-draftパラメータを適用して取得できる",
			existPosts: []*entity.Post{{
				ID:           "abcdefghijklmnopqrstuvwxy1",
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
			}, {
				ID:           "abcdefghijklmnopqrstuvwxy3",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink3",
				IsDraft:      true,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			offset:    0,
			pageSize:  0,
			condition: "is_draft = ?",
			params:    []interface{}{true},
			want: []*entity.Post{{
				ID:           "abcdefghijklmnopqrstuvwxy3",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink3",
				IsDraft:      true,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			wantErr: nil,
		},
		{
			name: "tagパラメータを適用して取得できる",
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
			}, {
				ID:           "abcdefghijklmnopqrstuvwxy6",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink2",
				IsDraft:      false,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}, {
				ID:           "abcdefghijklmnopqrstuvwxy7",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "new_content",
				Permalink:    "new_permalink3",
				IsDraft:      true,
				CreatedAt:    flextime.Now(),
				UpdatedAt:    flextime.Now(),
				PublishedAt:  flextime.Now(),
			}},
			existTag: &entity.Tag{
				ID:        "abcdefghijklmnopqrstuvwxy2",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			existPostsTags: &entity.PostsTags{
				ID:        "abcdefghijklmnopqrstuvwxy2",
				PostID:    "abcdefghijklmnopqrstuvwxy5",
				TagID:     "abcdefghijklmnopqrstuvwxy2",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			offset:    0,
			pageSize:  0,
			condition: "tags.name = ?",
			params:    []interface{}{"new_tag"},
			want: []*entity.Post{{
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
			wantErr: nil,
		},
		{
			name:       "投稿が存在しない場合はErrPostNotFoundを返す",
			existPosts: nil,
			want:       []*entity.Post{},
			wantErr:    entity.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existPosts != nil {
				if err := tx.Debug().Create(tt.existPosts).Error; err != nil {
					t.Fatal(err)
				}
			}
			if tt.existTag != nil {
				if err := tx.Debug().Create(tt.existTag).Error; err != nil {
					t.Fatal(err)
				}
			}
			if tt.existPostsTags != nil {
				if err := tx.Debug().Create(tt.existPostsTags).Error; err != nil {
					t.Fatal(err)
				}
			}
			var posts []*entity.Post
			tx.Find(&posts)
			for _, v := range posts {
				t.Logf("%v", v)
			}
			r := &PostRepository{db: tx.Debug()}
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

			if tt.existTag != nil {
				if err := tx.Delete(tt.existTag).Error; err != nil {
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

func TestPostRepository_Delete(t *testing.T) {
	tx := db.Begin()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name      string
		ID        string
		existPost *entity.Post
		want      *entity.Post
		wantErr   error
	}{
		{
			name: "存在する投稿を正常に削除できる",
			ID:   "abcdefghijklmnopqrstuvwxyz",
			existPost: &entity.Post{
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
			want:    nil,
			wantErr: nil,
		},
		{
			name:      "存在しないIDの場合ErrPostNotFoundを返す",
			ID:        "not_found",
			existPost: nil,
			want:      nil,
			wantErr:   entity.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.existPost != nil {
				if err := tx.Create(tt.existPost).Error; err != nil {
					t.Fatal(err)
				}
			}

			r := &PostRepository{db: tx}
			err := r.Delete(tt.ID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			//TODO 削除したことを確かめるテスト

		})
	}

	tx.Rollback()
}
