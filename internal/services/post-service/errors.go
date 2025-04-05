package postservice

import "fmt"

type PostServiceError struct {
	msg string
	err error
}

func (e PostServiceError) Error() string {
	return fmt.Sprintf("album service error: %v", e.msg)
}

func (e PostServiceError) Unwrap() error {
	return e.err
}

func Wrap(e error) PostServiceError {
	return PostServiceError{msg: fmt.Sprintf("album service error: %v", e), err: e}
}

var (
	ErrNotAuthor              = PostServiceError{msg: "does not have author rights"}
	ErrNotAuthorized          = PostServiceError{msg: "not authorized"}
	ErrInvalidFileContentType = PostServiceError{msg: "invalid file content type"}
)
