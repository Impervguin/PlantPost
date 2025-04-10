package authservice

import (
	"PlantSite/internal/models/auth"
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	sessions   SessionStorage
	repository auth.AuthRepository
	hasher     PasswdHasher
}

func NewAuthService(sessions SessionStorage, repository auth.AuthRepository, hasher PasswdHasher) *AuthService {
	if sessions == nil {
		panic("nil sessions")
	}
	if repository == nil {
		panic("nil repository")
	}
	if hasher == nil {
		panic("nil hasher")
	}
	return &AuthService{
		sessions:   sessions,
		repository: repository,
		hasher:     hasher,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) error {
	hashedPasswd, err := s.hasher.Hash([]byte(password))
	if err != nil {
		return err
	}

	member, err := auth.NewMember(name, email, hashedPasswd)
	if err != nil {
		return err
	}

	_, err = s.repository.Create(ctx, member)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, identifier, password string) (uuid.UUID, error) {
	user, err := s.repository.GetByEmail(ctx, identifier)
	if err != nil {
		user, err = s.repository.GetByName(ctx, identifier)
		if err != nil {
			return uuid.Nil, err
		}
	}

	if !user.Auth([]byte(password), s.hasher.Compare) {
		return uuid.Nil, ErrInvalidCredentials
	}

	sid := uuid.New()
	session := &Session{
		ID:        sid,
		MemberID:  user.ID(),
		ExpiresAt: time.Now().Add(SessionExpireTime),
	}

	err = s.sessions.Store(ctx, sid, session)
	if err != nil {
		return uuid.Nil, err
	}

	return sid, nil
}

func (s *AuthService) Logout(ctx context.Context, sid uuid.UUID) error {
	return s.sessions.Delete(ctx, sid)
}

func (s *AuthService) Authenticate(ctx context.Context, sid uuid.UUID) (context.Context, error) {
	session, err := s.sessions.Get(ctx, sid)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, ErrSessionExpired
	}

	user, err := s.repository.Get(ctx, session.MemberID)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, AuthContextKey, user)

	return newCtx, nil
}
