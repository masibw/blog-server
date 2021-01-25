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

func TestTagUseCase_StoreTag(t *testing.T) { // nolint:gocognit

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	flextime.Fix(time.Date(2021, 1, 22, 0, 0, 0, 0, loc))
	defer flextime.Restore()

	tests := []struct {
		name                 string
		tagDTO               *dto.TagDTO
		prepareMockTagRepoFn func(mock *mock_repository.MockTag)
		wantErr              error
	}{
		{
			name: "新規のタグを保存し、そのタグを返す",
			tagDTO: &dto.TagDTO{
				Name: "new_tag",
			},
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByName(gomock.Any()).Return(nil, entity.ErrTagNotFound)
				mock.EXPECT().Store(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "名前が登録済みの場合ErrNameAlreadyExistedエラーを返す",
			tagDTO: &dto.TagDTO{
				Name: "new_tag",
			},
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByName("new_tag").Return(&entity.Tag{}, nil)
				mock.EXPECT().Store(gomock.Any()).AnyTimes().Return(nil)
			},
			wantErr: entity.ErrTagNameAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			p := &TagUseCase{
				tagRepository: mr,
			}

			got, err := p.StoreTag(tt.tagDTO)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StoreTag() error = %v, wantErr %v", err, tt.wantErr)
			}

			if errors.Is(err, entity.ErrTagNameAlreadyExisted) {
				if got == nil {
					return
				}
				t.Errorf("StoreTag() got = %v, want = nil", got)
			}

			if got.ID == "" {
				t.Errorf("StoreTag() ID nil want UULD")
			}
			if got.CreatedAt.Unix() == 0 || got.UpdatedAt.Unix() == 0 {
				t.Errorf("StoreTag() time.Time field did not filled with value")
			}

			if got.Name != tt.tagDTO.Name {
				t.Errorf("StoreTag() Name does not match got: %v, want: %v", got, tt.tagDTO)
			}
		})
	}
}

func TestTagUseCase_GetTags(t *testing.T) {

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
		Name:      "new_tag",
		CreatedAt: flextime.Now(),
		UpdatedAt: flextime.Now(),
	}}

	tests := []struct {
		name                 string
		prepareMockTagRepoFn func(mock *mock_repository.MockTag)
		want                 []*dto.TagDTO
		wantErr              bool
	}{
		{
			name: "tagDTOsを返すこと",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(existsTags, nil)
			},
			want: []*dto.TagDTO{
				{
					ID:        "abcdefghijklmnopqrstuvwxyz",
					Name:      "new_tag",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				},
				{
					ID:        "abcdefghijklmnopqrstuvwxy2",
					Name:      "new_tag",
					CreatedAt: flextime.Now(),
					UpdatedAt: flextime.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "FindAllがエラーを返した時はtagDTOsが空であること",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("dummy error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			p := &TagUseCase{
				tagRepository: mr,
			}

			// このGetTagsの責務はパラメータを受け取ってtagDTOsを返すだけなのでパラメータの中身はなんでも良い(はず)
			got, err := p.GetTags(0, 0)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTags() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetTags() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestTagUseCase_GetTag(t *testing.T) {

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
		want                 *dto.TagDTO
		wantErr              bool
	}{
		{
			name: "tagDTOを返すこと",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByID(gomock.Any()).Return(existsTag, nil)
			},
			want: &dto.TagDTO{
				ID:        "abcdefghijklmnopqrstuvwxyz",
				Name:      "new_tag",
				CreatedAt: flextime.Now(),
				UpdatedAt: flextime.Now(),
			},
			ID:      "abcdefghijklmnopqrstuvwxyz",
			wantErr: false,
		},
		{
			name: "FindByIDがエラーを返した時はtagDTOが空であること",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().FindByID("not_found").Return(nil, entity.ErrTagNotFound)
			},
			ID:      "not_found",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			p := &TagUseCase{
				tagRepository: mr,
			}

			got, err := p.GetTag(tt.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTag() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetTag() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestTagUseCase_DeleteTag(t *testing.T) {

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
		wantErr              bool
	}{
		{
			name: "削除に成功した場合はエラーを返さないこと",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			ID:      "abcdefghijklmnopqrstuvwxyz",
			wantErr: false,
		},
		{
			name: "Deleteがエラーを返した時はエラーを返すこと",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTag) {
				mock.EXPECT().Delete("not_found").Return(entity.ErrTagNotFound)
			},
			ID:      "not_found",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockTag(ctrl)
			tt.prepareMockTagRepoFn(mr)
			p := &TagUseCase{
				tagRepository: mr,
			}

			err := p.DeleteTag(tt.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
