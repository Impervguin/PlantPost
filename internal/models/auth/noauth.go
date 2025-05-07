package auth

import "github.com/google/uuid"

type NoAuthUser struct{}

func (u NoAuthUser) HasAuthorRights() bool {
	return false
}

func (u NoAuthUser) HasMemberRights() bool {
	return false
}

func NewNoAuthUser() User {
	return &NoAuthUser{}
}

func (u NoAuthUser) Auth(_ []byte, _ func(hashPasswd []byte, plainPasswd []byte) (bool, error)) bool {
	return false
}

func (u NoAuthUser) ID() uuid.UUID {
	return uuid.Nil
}

func (u NoAuthUser) IsAuthenticated() bool {
	return false
}
