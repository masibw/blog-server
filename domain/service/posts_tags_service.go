package service

import (
	"errors"
	"fmt"

	"github.com/masibw/blog-server/domain/entity"
	"github.com/masibw/blog-server/domain/repository"
)

type PostsTagsService struct {
	postsTagsRepository repository.PostsTags
	postRepository      repository.Post
	tagRepository       repository.Tag
}

func NewPostsTagsService(postsTagsRepository repository.PostsTags, postRepository repository.Post, tagRepository repository.Tag) *PostsTagsService {
	return &PostsTagsService{
		postsTagsRepository: postsTagsRepository,
		postRepository:      postRepository,
		tagRepository:       tagRepository,
	}
}

func (p *PostsTagsService) LinkPostTags(postID string, tagNames []string) ([]*entity.Tag, error) {
	var err error

	// 実際に投稿が存在するかのチェックであり結果は使わない
	_, err = p.postRepository.FindByID(postID)
	if err != nil {
		return nil, fmt.Errorf("LinkPostTags() get post: %w", err)
	}

	err = p.postsTagsRepository.DeleteByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("LinkPostTags() delete : %w", err)
	}

	// タグの重複を削除する
	m := make(map[string]bool)
	var uniqTagNames []string
	for _, tagName := range tagNames {
		if !m[tagName] {
			m[tagName] = true
			uniqTagNames = append(uniqTagNames, tagName)
		}
	}

	postsTagsSlice := make([]*entity.PostsTags, 0)
	tags := make([]*entity.Tag, 0)
	for _, tagName := range uniqTagNames {
		var postsTags *entity.PostsTags
		var tag *entity.Tag
		postsTags, tag, err = p.getTagEntity(postID, tagName)
		if err != nil {
			return nil, fmt.Errorf("LinkPostTags() tagName =%s :%w", tagName, err)
		}
		postsTagsSlice = append(postsTagsSlice, postsTags)
		tags = append(tags, tag)
	}
	err = p.postsTagsRepository.Store(postsTagsSlice)
	if err != nil {
		return nil, fmt.Errorf("LinkPostTags() store posts_tags post id =%v tagNames =%v: %w", postID, tagNames, err)
	}

	return tags, nil
}

func (p *PostsTagsService) getTagEntity(postID, tagName string) (postsTags *entity.PostsTags, tag *entity.Tag, err error) {

	// tagNameからタグを取得する
	tag, err = p.tagRepository.FindByName(tagName)
	if err != nil && !errors.Is(err, entity.ErrTagNotFound) {
		err = fmt.Errorf("getTagEntity() store tag name=%v: %w", tagName, err)
		return
	}
	// タグが存在しなければ作成する
	if errors.Is(err, entity.ErrTagNotFound) {
		tag = entity.NewTag(tagName)
		err = p.tagRepository.Store(tag)
		if err != nil {
			err = fmt.Errorf("getTagEntity() store posts_tags tag name=%v: %w", tagName, entity.ErrPostsTagsAlreadyExisted)
			return
		}
	}

	postsTags = entity.NewPostsTags(
		postID,
		tag.ID,
	)
	err = nil
	return
}
