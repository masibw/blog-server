package usecase

import (
	"errors"
	"fmt"

	"github.com/Songmu/flextime"
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
	if !post.IsDraft {
		post.PublishedAt = flextime.Now()
	}
	err = p.postRepository.Store(post)
	if err != nil {
		return nil, fmt.Errorf("store post title=%v: %w", postDTO.Title, err)
	}

	return post.ConvertToDTO(), nil
}

func (p *PostUseCase) GetPosts(offset, pageSize int, condition string, params []interface{}) (postDTOs []*dto.PostDTO, err error) {
	var posts []*entity.Post
	posts, err = p.postRepository.FindAll(offset, pageSize, condition, params)
	if err != nil {
		err = fmt.Errorf("get posts: %w", err)
		return
	}

	for _, post := range posts {

		// Markdownをhtmlへパースしている
		post.ConvertContentToHTML()

		postDTOs = append(postDTOs, post.ConvertToDTO())
	}

	return
}

func (p *PostUseCase) GetPost(id string) (postDTO *dto.PostDTO, err error) {
	var post *entity.Post
	post, err = p.postRepository.FindByID(id)
	if err != nil {
		err = fmt.Errorf("get post: %w", err)
		return
	}
	post.ConvertContentToHTML()
	postDTO = post.ConvertToDTO()
	return
}

func (p *PostUseCase) DeletePost(id string) (err error) {
	err = p.postRepository.Delete(id)
	if err != nil {
		err = fmt.Errorf("delete post: %w", err)
		return
	}
	return nil
}
