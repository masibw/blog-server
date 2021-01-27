package service

import (
	"errors"
	"fmt"

	"github.com/masibw/blog-server/domain/entity"
	"github.com/masibw/blog-server/domain/repository"
)

func LinkPostTags(postsTagsRepository repository.PostsTags, postRepository repository.Post, tagRepository repository.Tag, postID, tagName string) error {
	var postsTags *entity.PostsTags
	var err error

	//TODO 変更がなかったときにどうする？

	// 既に同じpostIDとtagIDが紐づいていないかのチェック
	postsTags, err = postsTagsRepository.FindByPostIDAndTagName(postID, tagName)
	if err != nil && !errors.Is(err, entity.ErrPostsTagsNotFound) {
		return fmt.Errorf("store posts_tags post id =%v tag name=%v: %w", postID, tagName, err)
	}
	if postsTags != nil {
		return fmt.Errorf("store posts_tags post id =%v tag name=%v: %w", postID, tagName, entity.ErrPostsTagsAlreadyExisted)
	}

	// 実際に投稿が存在するかのチェックであり結果は使わない
	_, err = postRepository.FindByID(postID)
	if err != nil {
		return fmt.Errorf("get post: %w", err)
	}

	// tagNameからタグを取得する
	tag, err := tagRepository.FindByName(tagName)
	if err != nil && !errors.Is(err, entity.ErrTagNotFound) {
		return fmt.Errorf("store tag name=%v: %w", tagName, err)
	}
	// タグが存在しなければ作成する
	if errors.Is(err, entity.ErrTagNotFound) {
		tag = entity.NewTag(tagName)
		err := tagRepository.Store(tag)
		if err != nil {
			return fmt.Errorf("store posts_tags tag name=%v: %w", tagName, entity.ErrPostsTagsAlreadyExisted)
		}
	}

	postsTags = entity.NewPostsTags(
		postID,
		tag.ID,
	)

	err = postsTagsRepository.Store(postsTags)
	if err != nil {
		return fmt.Errorf("store posts_tags post id =%v tag name=%v: %w", postID, tagName, err)
	}

	return nil
}
