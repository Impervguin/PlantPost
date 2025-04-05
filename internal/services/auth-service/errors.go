package authservice

import "fmt"

type AuthServiceError struct {
	msg string
}

func (e *AuthServiceError) Error() string {
	return fmt.Sprintf("auth service error: %v", e.msg)
}

var (
	ErrInvalidCredentials = &AuthServiceError{msg: "invalid credentials"}
	ErrSessionExpired     = &AuthServiceError{msg: "session expired"}
)
