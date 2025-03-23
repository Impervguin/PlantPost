package auth

import (
	"context"

	"github.com/google/uuid"
)

type User interface {
	HasAuthorRights() bool
	HasMemberRights() bool
	Auth(passwd []byte, authFunc func(hashPasswd []byte, plainPasswd []byte) (bool, error)) bool
	ID() uuid.UUID
}

type AuthRepository interface {
	Get(ctx context.Context, id uuid.UUID) (User, error)
	Create(ctx context.Context, mem *Member) (User, error)
	Update(ctx context.Context, id uuid.UUID, updateFn func(*User) (*User, error)) (*User, error)
	GetByName(ctx context.Context, name string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
