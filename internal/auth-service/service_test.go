package authservice_test

import (
	"context"
	"testing"
	"time"

	authservice "PlantSite/internal/auth-service"
	"PlantSite/internal/models/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSessionStorage implements SessionStorage interface
type MockSessionStorage struct {
	mock.Mock
}

func (m *MockSessionStorage) Store(ctx context.Context, sid uuid.UUID, session *authservice.Session) error {
	args := m.Called(ctx, sid, session)
	return args.Error(0)
}

func (m *MockSessionStorage) Get(ctx context.Context, sid uuid.UUID) (*authservice.Session, error) {
	args := m.Called(ctx, sid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authservice.Session), args.Error(1)
}

func (m *MockSessionStorage) Delete(ctx context.Context, sid uuid.UUID) error {
	args := m.Called(ctx, sid)
	return args.Error(0)
}

func (m *MockSessionStorage) ClearExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockAuthRepository implements auth.AuthRepository interface
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Get(ctx context.Context, id uuid.UUID) (auth.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepository) Create(ctx context.Context, mem *auth.Member) (auth.User, error) {
	args := m.Called(ctx, mem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*auth.User) (*auth.User, error)) (*auth.User, error) {
	args := m.Called(ctx, id, updateFn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockAuthRepository) GetByName(ctx context.Context, name string) (auth.User, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepository) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

// MockPasswdHasher implements PasswdHasher interface
type MockPasswdHasher struct {
	mock.Mock
}

func (m *MockPasswdHasher) Hash(password []byte) ([]byte, error) {
	args := m.Called(password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPasswdHasher) Compare(hashPasswd, plainPasswd []byte) (bool, error) {
	args := m.Called(hashPasswd, plainPasswd)
	return args.Bool(0), args.Error(1)
}

// MockUser implements auth.User interface
type MockUser struct {
	mock.Mock
}

func (m *MockUser) ID() uuid.UUID {
	return m.Called().Get(0).(uuid.UUID)
}

func (m *MockUser) Name() string {
	return m.Called().String(0)
}

func (m *MockUser) Email() string {
	return m.Called().String(0)
}

func (m *MockUser) HashedPassword() []byte {
	return m.Called().Get(0).([]byte)
}

func (m *MockUser) Auth(password []byte, compareFn func(hashPasswd, plainPasswd []byte) (bool, error)) bool {
	args := m.Called(password, compareFn)
	return args.Bool(0)
}

func (m *MockUser) HasMemberRights() bool {
	return m.Called().Bool(0)
}

func (m *MockUser) HasAuthorRights() bool {
	return m.Called().Bool(0)
}

func TestAuthService(t *testing.T) {
	ctx := context.Background()
	validUserID := uuid.New()
	validSessionID := uuid.New()
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "securepassword"
	hashedPassword := []byte("hashedpassword")

	t.Run("Register", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			hasher.On("Hash", []byte(validPassword)).Return(hashedPassword, nil)

			mockUser := new(MockUser)
			repo.On("Create", ctx, mock.AnythingOfType("*auth.Member")).Return(mockUser, nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.NoError(t, err)

			hasher.AssertExpectations(t)
			repo.AssertExpectations(t)
		})

		t.Run("HashError", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			hasher.On("Hash", []byte(validPassword)).Return(nil, assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("CreateUserError", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			hasher.On("Hash", []byte(validPassword)).Return(hashedPassword, nil)
			repo.On("Create", ctx, mock.AnythingOfType("*auth.Member")).Return(nil, assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Login", func(t *testing.T) {
		t.Run("SuccessWithEmail", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			mockUser := new(MockUser)
			mockUser.On("ID").Return(validUserID)
			mockUser.On("Auth", []byte(validPassword), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(true)

			repo.On("GetByEmail", ctx, validEmail).Return(mockUser, nil)
			sessions.On("Store", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*authservice.Session")).Return(nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			sid, err := svc.Login(ctx, validEmail, validPassword)
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, sid)

			repo.AssertExpectations(t)
			sessions.AssertExpectations(t)
			mockUser.AssertExpectations(t)
		})

		t.Run("SuccessWithName", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			mockUser := new(MockUser)
			mockUser.On("ID").Return(validUserID)
			mockUser.On("Auth", []byte(validPassword), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(true)

			repo.On("GetByEmail", ctx, validName).Return(nil, assert.AnError)
			repo.On("GetByName", ctx, validName).Return(mockUser, nil)
			sessions.On("Store", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*authservice.Session")).Return(nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			sid, err := svc.Login(ctx, validName, validPassword)
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, sid)
		})

		t.Run("InvalidCredentials", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			mockUser := new(MockUser)
			mockUser.On("Auth", []byte("wrongpassword"), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(false)

			repo.On("GetByEmail", ctx, validEmail).Return(mockUser, nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Login(ctx, validEmail, "wrongpassword")
			require.Error(t, err)
			assert.ErrorIs(t, err, authservice.ErrInvalidCredentials)
		})

		t.Run("UserNotFound", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			repo.On("GetByEmail", ctx, validEmail).Return(nil, assert.AnError)
			repo.On("GetByName", ctx, validEmail).Return(nil, assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Login(ctx, validEmail, validPassword)
			require.Error(t, err)
		})

		t.Run("SessionStoreError", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			mockUser := new(MockUser)
			mockUser.On("ID").Return(validUserID)
			mockUser.On("Auth", []byte(validPassword), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(true)

			repo.On("GetByEmail", ctx, validEmail).Return(mockUser, nil)
			sessions.On("Store", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*authservice.Session")).Return(assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Login(ctx, validEmail, validPassword)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Logout", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			sessions.On("Delete", ctx, validSessionID).Return(nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			err := svc.Logout(ctx, validSessionID)
			require.NoError(t, err)

			sessions.AssertExpectations(t)
		})

		t.Run("Error", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			sessions.On("Delete", ctx, validSessionID).Return(assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			err := svc.Logout(ctx, validSessionID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Authenticate", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			validSession := &authservice.Session{
				ID:        validSessionID,
				MemberID:  validUserID,
				ExpiresAt: time.Now().Add(time.Hour),
			}

			mockUser := new(MockUser)

			sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
			repo.On("Get", ctx, validUserID).Return(mockUser, nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			newCtx, err := svc.Authenticate(ctx, validSessionID)
			require.NoError(t, err)
			assert.NotNil(t, newCtx)
			assert.Equal(t, mockUser, newCtx.Value(authservice.AuthContextKey))

			sessions.AssertExpectations(t)
			repo.AssertExpectations(t)
		})

		t.Run("SessionNotFound", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			sessions.On("Get", ctx, validSessionID).Return(nil, assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Authenticate(ctx, validSessionID)
			require.Error(t, err)
		})

		t.Run("SessionExpired", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			expiredSession := &authservice.Session{
				ID:        validSessionID,
				MemberID:  validUserID,
				ExpiresAt: time.Now().Add(-time.Hour),
			}

			sessions.On("Get", ctx, validSessionID).Return(expiredSession, nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Authenticate(ctx, validSessionID)
			require.Error(t, err)
			assert.ErrorIs(t, err, authservice.ErrSessionExpired)
		})

		t.Run("UserNotFound", func(t *testing.T) {
			repo := new(MockAuthRepository)
			sessions := new(MockSessionStorage)
			hasher := new(MockPasswdHasher)

			validSession := &authservice.Session{
				ID:        validSessionID,
				MemberID:  validUserID,
				ExpiresAt: time.Now().Add(time.Hour),
			}

			sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
			repo.On("Get", ctx, validUserID).Return(nil, assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Authenticate(ctx, validSessionID)
			require.Error(t, err)
		})
	})
}
