package usecase

import (
	"fmt"

	"github.com/masibw/blog-server/domain/repository"
)

type PostsTagsUseCase struct {
	postsTagsRepository repository.PostsTags
}

func NewPostsTagsUseCase(postsTagsRepository repository.PostsTags) *PostsTagsUseCase {
	return &PostsTagsUseCase{postsTagsRepository: postsTagsRepository}
}

func (p *PostsTagsUseCase) DeletePostsTags(id string) (err error) {
	err = p.postsTagsRepository.Delete(id)
	if err != nil {
		err = fmt.Errorf("delete posts_tags: %w", err)
		return
	}
	return nil
}
