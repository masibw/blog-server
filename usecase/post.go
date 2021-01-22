package usecase

import (
	"errors"
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
	var post *entity.Post
	var err error
	post, err = p.postRepository.FindByPermalink(postDTO.Permalink)
	if err != nil && !errors.Is(err, entity.ErrPostNotFound) {
		return nil, fmt.Errorf("store post title=%v: %w", postDTO.Title, err)
	}
	if post != nil {
		return nil, fmt.Errorf("store post permalink=%v: %w", postDTO.Permalink, entity.ErrPermalinkAlreadyExisted)
	}

	post = entity.NewPost(
		postDTO.Title,
		postDTO.ThumbnailURL,
		postDTO.Content,
		postDTO.Permalink,
		*postDTO.IsDraft,
	)
	err = p.postRepository.Store(post)
	if err != nil {
		return nil, fmt.Errorf("store post title=%v: %w", postDTO.Title, err)
	}

	return post.ConvertToDTO(), nil
}
