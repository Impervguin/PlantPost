package plantservice

import "fmt"

type PlantServiceError struct {
	msg string
	err error
}

func (e PlantServiceError) Error() string {
	return fmt.Sprintf("plant service error: %v", e.msg)
}

func (e PlantServiceError) Unwrap() error {
	return e.err
}

func Wrap(e error) PlantServiceError {
	return PlantServiceError{msg: fmt.Sprintf("plant service error: %v", e), err: e}
}

var (
	ErrNotAuthor     = PlantServiceError{msg: "does not have author rights"}
	ErrNotAuthorized = PlantServiceError{msg: "not authorized"}
)
