package post

import "errors"

var (
	ErrPostNotFound        = errors.New("post not found")
	ErrContentParsingError = errors.New("content parsing error")
)
