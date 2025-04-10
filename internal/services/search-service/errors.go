package searchservice

import "fmt"

type SearchServiceError struct {
	msg string
	err error
}

func (e SearchServiceError) Error() string {
	return fmt.Sprintf("search service error: %v", e.msg)
}

func (e SearchServiceError) Unwrap() error {
	return e.err
}

func Wrap(e error) SearchServiceError {
	return SearchServiceError{msg: fmt.Sprintf("search service error: %v", e), err: e}
}

var (
	ErrNotAuthor     = SearchServiceError{msg: "does not have author rights"}
	ErrNotAuthorized = SearchServiceError{msg: "not authorized"}
)
