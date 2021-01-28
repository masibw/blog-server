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

func (p *PostUseCase) UpdatePost(postDTO *dto.PostDTO) (*dto.PostDTO, error) {

	// 下書きじゃないのにTitleとContent,Permalinkに未入力項目があればエラー
	if !*postDTO.IsDraft {
		errMsg := ""
		if postDTO.Title == "" {
			errMsg += "title is nil "
		}
		if postDTO.Content == "" {
			errMsg += "content is nil "
		}
		if postDTO.Permalink == "" {
			errMsg += "permalink is nil "
		}
		if errMsg != "" {
			return nil, fmt.Errorf("update post some fields that have not been filled %s: %w", errMsg, entity.ErrPostHasEmptyField)
		}
	}

	var post *entity.Post
	var err error

	// 更新対象の投稿が存在するかの確認
	post, err = p.postRepository.FindByID(postDTO.ID)
	if err != nil && !errors.Is(err, entity.ErrPostNotFound) {
		return nil, fmt.Errorf("update post title=%v: %w", postDTO.Title, err)
	}
	if errors.Is(err, entity.ErrPostNotFound) {
		return nil, fmt.Errorf("update post not found ID=%v: %w", postDTO.ID, entity.ErrPostNotFound)
	}

	var permalinkPost *entity.Post
	// 重複確認の処理をDomainServiceに切り出すべきだけど2箇所なので一旦保留
	permalinkPost, err = p.postRepository.FindByPermalink(postDTO.Permalink)
	if err != nil && !errors.Is(err, entity.ErrPostNotFound) {
		return nil, fmt.Errorf("update post title=%v: %w", postDTO.Title, err)
	}

	// 更新する投稿と違うIDを持ち，既に更新先Permalinkを持つ投稿があるとエラー
	if permalinkPost != nil && permalinkPost.ID != postDTO.ID {
		return nil, fmt.Errorf("update post permalink=%v: %w", postDTO.Permalink, entity.ErrPermalinkAlreadyExisted)
	}

	post.ConvertFromDTO(postDTO)

	// 初めて公開するときのみ投稿時間を設定する
	if post.PublishedAt.IsZero() && !post.IsDraft {
		post.PublishedAt = flextime.Now()
	}

	err = p.postRepository.Update(post)

	if err != nil {
		return nil, fmt.Errorf("update post title=%v: %w", postDTO.Title, err)
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
