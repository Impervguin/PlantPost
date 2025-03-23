package authservice

import (
	"PlantSite/internal/models/auth"
	"context"
	"time"
)

const (
	SessionExpireTime = time.Hour
)

type authContextKey int

const (
	AuthContextKey authContextKey = iota
)

func UserFromContext(ctx context.Context) auth.User {
	if user, ok := ctx.Value(AuthContextKey).(auth.User); ok {
		return user
	}
	return nil
}
