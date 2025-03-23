package albumservice

import "fmt"

type AlbumServiceError struct {
	msg string
	err error
}

func (e AlbumServiceError) Error() string {
	return fmt.Sprintf("album service error: %v", e.msg)
}

func (e AlbumServiceError) Unwrap() error {
	return e.err
}

func Wrap(e error) AlbumServiceError {
	return AlbumServiceError{msg: fmt.Sprintf("album service error: %v", e), err: e}
}

var (
	ErrNotMember     = AlbumServiceError{msg: "does not have member rights"}
	ErrNotAuthorized = AlbumServiceError{msg: "not authorized"}
	ErrNotOwner      = AlbumServiceError{msg: "not owner"}
)
