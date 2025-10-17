package authservice

import (
	"PlantSite/internal/models/auth"
	"context"

	"github.com/google/uuid"
)

type AuthServiceContract interface {
	Login(ctx context.Context, identifier, password string) (uuid.UUID, error)
	Register(ctx context.Context, name, email, password string) error
	Logout(ctx context.Context) error
	Authenticate(ctx context.Context, sid uuid.UUID) context.Context
	UserFromContext(ctx context.Context) auth.User
}
