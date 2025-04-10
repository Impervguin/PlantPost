package authservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	MemberID  uuid.UUID
	ExpiresAt time.Time
}

func NewSession(id uuid.UUID, memberID uuid.UUID, expiresAt time.Time) (*Session, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("session ID cannot be nil")
	}
	if memberID == uuid.Nil {
		return nil, fmt.Errorf("session member ID cannot be nil")
	}
	return &Session{
		ID:        id,
		MemberID:  memberID,
		ExpiresAt: expiresAt,
	}, nil
}

type SessionStorage interface {
	Get(ctx context.Context, sid uuid.UUID) (*Session, error)
	Store(ctx context.Context, sid uuid.UUID, session *Session) error
	Delete(ctx context.Context, sid uuid.UUID) error
	ClearExpired(ctx context.Context) error
}
