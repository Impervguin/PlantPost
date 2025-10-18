package authservice

import (
	"PlantSite/internal/models/auth"
	"PlantSite/internal/utils/logs"
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
	logs.Debugf("AuthService.Register: got name=%s, email=%s, password=***", name, email)
	hashedPasswd, err := s.hasher.Hash([]byte(password))
	if err != nil {
		return err
	}
	logs.Debugf("AuthService.Register: hashed password=%s", hashedPasswd)

	member, err := auth.NewMember(name, email, hashedPasswd)
	if err != nil {
		return err
	}

	logs.Debugf("AuthService.Register: created member %s", member.ID())

	_, err = s.repository.Create(ctx, member)
	if err != nil {
		return err
	}

	logs.Debugf("AuthService.Register: saved member %s", member.ID())
	return nil
}

func (s *AuthService) Login(ctx context.Context, identifier, password string) (uuid.UUID, error) {
	user, err := s.repository.GetByEmail(ctx, identifier)
	if err != nil {
		user, err = s.repository.GetByName(ctx, identifier)
		if err != nil {
			return uuid.Nil, err
		}
		logs.Debugf("AuthService.Login: got user %s by name", user.ID())
	} else {
		logs.Debugf("AuthService.Login: got user %s by email", user.ID())
	}

	if !user.Auth([]byte(password), s.hasher.Compare) {
		return uuid.Nil, ErrInvalidCredentials
	}
	logs.Debugf("AuthService.Login: user %s password is valid", user.ID())

	sid := uuid.New()
	session := &Session{
		ID:        sid,
		MemberID:  user.ID(),
		ExpiresAt: time.Now().Add(SessionExpireTime),
	}
	logs.Debugf("AuthService.Login: created session %s for user %s", sid, user.ID())

	err = s.sessions.Store(ctx, sid, session)
	if err != nil {
		return uuid.Nil, err
	}

	logs.Debugf("AuthService.Login: saved session %s for user %s", sid, user.ID())

	return sid, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	sid := s.sessionFromContext(ctx)
	logs.Debugf("AuthService.Logout: got session %s", sid)
	sess, err := s.sessions.Get(ctx, sid)
	if err != nil {
		return err
	}
	logs.Debugf("AuthService.Logout: got session %s for user %s", sess.ID, sess.MemberID)

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

	logs.Debugf("AuthService.Authenticate: got user %s for session %s", userID, sid)

	ctx = context.WithValue(ctx, AuthContextKey, userID)
	logs.Debugf("AuthService.Authenticate: set user %s for session %s in context", userID, sid)
	ctx = context.WithValue(ctx, sessionContextKey, sid)
	logs.Debugf("AuthService.Authenticate: set session %s in context", sid)

	return ctx
}

func (s *AuthService) UserFromContext(ctx context.Context) auth.User {
	if userID, ok := ctx.Value(AuthContextKey).(uuid.UUID); ok {
		if userID == uuid.Nil {
			return auth.NewNoAuthUser()
		}
		logs.Debugf("AuthService.UserFromContext: got user ID %s", userID)
		user, err := s.repository.Get(ctx, userID)
		if err != nil {
			return auth.NewNoAuthUser()
		}
		logs.Debugf("AuthService.UserFromContext: got user %s", user.ID())
		return user
	}
	logs.Debugf("AuthService.UserFromContext: no user in context")
	return auth.NewNoAuthUser()
}

func (s *AuthService) sessionFromContext(ctx context.Context) uuid.UUID {
	if sid, ok := ctx.Value(sessionContextKey).(uuid.UUID); ok {
		return sid
	}
	return uuid.Nil
}
