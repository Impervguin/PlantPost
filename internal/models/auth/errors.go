package auth

import (
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

var (
	ErrBaseNoRights   = errors.New("user has no rights")
	ErrNoAuthorRights = fmt.Errorf("%w: author", ErrBaseNoRights)
	ErrNotAuthorized  = errors.New("not authorized")
	ErrNoMemberRights = fmt.Errorf("%w: %w", ErrBaseNoRights, ErrNotAuthorized)
)
