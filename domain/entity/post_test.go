package entity

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPost_ConvertContentToHTML(t *testing.T) {

	tests := []struct {
		name string
		post *Post
		want *Post
	}{
		{
			name: "正常にmarkdownをhtmlに変換できる",
			post: &Post{
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
			want: &Post{
				ID:           "abcdefghijklmnopqrstuvwxyz",
				Title:        "new_post",
				ThumbnailURL: "new_thumbnail_url",
				Content:      "<p>new_content</p>\n",
				Permalink:    "new_permalink",
				IsDraft:      false,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
				PublishedAt:  time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.post.ConvertContentToHTML()
			if diff := cmp.Diff(tt.want, tt.post); diff != "" {
				t.Errorf("ConvertContentToHTML() mismatch (-want +got):\n%s", diff)
			}

		})
	}

}
