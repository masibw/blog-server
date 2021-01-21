package usecase

import (
	"fmt"

	"github.com/masibw/blog-server/domain/dto"
	"github.com/masibw/blog-server/domain/entity"
	"github.com/masibw/blog-server/domain/repository"
)

type PostUseCase struct {
	postRepository repository.Post
}

func NewPostUseCase(postRepository repository.Post) *PostUseCase {
	return &PostUseCase{postRepository: postRepository}
}

func (p *PostUseCase) StorePost(postDTO *dto.PostDTO) (*dto.PostDTO, error) {
	post := entity.NewPost(
		postDTO.ThumbnailURL,
		postDTO.Title,
		postDTO.Content,
		postDTO.Permalink,
		*postDTO.IsDraft,
	)
	err := p.postRepository.Store(post)
	if err != nil {
		return nil, fmt.Errorf("store post title=%v: %w", postDTO.Title, err)
	}

	return post.ConvertToDTO(), nil
}
