package entity

import "errors"

// エラーの分岐でMYSQLのエラーコードを用いている https://dev.mysql.com/doc/refman/5.6/ja/error-messages-server.html
// エラー: 1062 SQLSTATE: 23000 (ER_DUP_ENTRY)

var (
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
)
