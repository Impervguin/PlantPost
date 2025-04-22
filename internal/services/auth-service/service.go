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

func (s *AuthService) Logout(ctx context.Context) error {
	sess, err := s.sessions.Get(ctx, uuid.Nil)
	if err != nil {
		return err
	}

	return s.sessions.Delete(ctx, sess.ID)
}

func (s *AuthService) authenticate(ctx context.Context, sid uuid.UUID) (uuid.UUID, error) {
	session, err := s.sessions.Get(ctx, sid)
	if err != nil {
		return uuid.Nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return uuid.Nil, ErrSessionExpired
	}

	return session.MemberID, nil
}

func (s *AuthService) Authenticate(ctx context.Context, sid uuid.UUID) context.Context {
	userID, err := s.authenticate(ctx, sid)
	if err != nil {
		userID = uuid.Nil
	}

	ctx = context.WithValue(ctx, AuthContextKey, userID)

	return ctx
}

func (s *AuthService) UserFromContext(ctx context.Context) auth.User {
	if userID, ok := ctx.Value(AuthContextKey).(uuid.UUID); ok {
		if userID == uuid.Nil {
			return auth.NewNoAuthUser()
		}
		s.repository.Get(ctx, userID)
		user, err := s.repository.Get(ctx, userID)
		if err != nil {
			return auth.NewNoAuthUser()
		}
		return user
	}

	return auth.NewNoAuthUser()
}
