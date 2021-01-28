package entity

import "errors"

// エラーの分岐でMYSQLのエラーコードを用いている https://dev.mysql.com/doc/refman/5.6/ja/error-messages-server.html
// エラー: 1062 SQLSTATE: 23000 (ER_DUP_ENTRY)

var (
	ErrInternalServerError = errors.New("internal server error")
	// ErrUserNotFound はユーザが存在しないエラーを表します。
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExisted はユーザが既に存在しているエラーを表します。
	ErrUserAlreadyExisted = errors.New("user has already existed")

	// ErrPostNotFound は投稿が存在しないエラーを表します。
	ErrPostNotFound = errors.New("post not found")
	// ErrPostAlreadyExisted は投稿が既に存在しているエラーを表します。
	ErrPostAlreadyExisted = errors.New("post has already existed")
	// ErrPermalinkAlreadyExisted はパーマリンクが既に存在しているエラーを表します。
	ErrPermalinkAlreadyExisted = errors.New("permalink has already existed")
	// ErrPostHasEmptyField は投稿に未入力項目があるエラーを表します。
	ErrPostHasEmptyField = errors.New("some fields that have not been filled")

	// ErrTagNotFound はタグが存在しないエラーを表します。
	ErrTagNotFound = errors.New("tag not found")
	// ErrTagAlreadyExisted はタグが既に存在しているエラーを表します。
	ErrTagAlreadyExisted = errors.New("tag has already existed")
	// ErrTagNameAlreadyExisted はその名前のタグが既に存在しているエラーを表します。
	ErrTagNameAlreadyExisted = errors.New("tag name has already existed")

	// ErrPostsTagsNotFound はタグが存在しないエラーを表します。
	ErrPostsTagsNotFound = errors.New("posts_tags not found")
	// ErrPostsTagsAlreadyExisted はタグが既に存在しているエラーを表します。
	ErrPostsTagsAlreadyExisted = errors.New("posts_tags has already existed")
	// ErrPostsTagsCombinationAlreadyExisted はその投稿に同じタグが既に存在しているエラーを表します。
	ErrPostsTagsCombinationAlreadyExisted = errors.New("posts_tags combination has already existed")
)
