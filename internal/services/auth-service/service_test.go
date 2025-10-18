//go:build unit

package authservice_test

import (
	"context"
	"testing"
	"time"

	"PlantSite/internal/models/auth"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	"PlantSite/internal/utils/logs"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AuthServiceTestSuite struct {
	suite.Suite
}

func (s *AuthServiceTestSuite) BeforeAll(t provider.T) {
	logs.InitNoopLogger()
}

func (s *AuthServiceTestSuite) BeforeEach(t provider.T) {
	t.Epic("Authentication")
	t.Feature("Auth Service")
}

func (s *AuthServiceTestSuite) TestRegister(t provider.T) {
	t.Tags("registration", "positive")
	t.Description("Test user registration functionality")
	t.Parallel()

	ctx := context.Background()
	// validUserID := uuid.New()
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "securepassword"
	hashedPassword := []byte("hashedpassword")

	t.Run("Successful user registration", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		t.WithNewStep("Setup password hashing", func(pctx provider.StepCtx) {
			hasher.On("Hash", []byte(validPassword)).Return(hashedPassword, nil)
		})

		mockUser := new(authmock.MockUser)
		t.WithNewStep("Setup user creation", func(pctx provider.StepCtx) {
			repo.On("Create", ctx, mock.AnythingOfType("*auth.Member")).Return(mockUser, nil)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		t.WithNewStep("Perform registration", func(pctx provider.StepCtx) {
			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			hasher.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	})

	t.Run("Password hashing error during registration", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		t.WithNewStep("Setup hashing error", func(pctx provider.StepCtx) {
			hasher.On("Hash", []byte(validPassword)).Return(nil, assert.AnError)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		t.WithNewStep("Attempt registration with hashing error", func(pctx provider.StepCtx) {
			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("User creation error during registration", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		t.WithNewStep("Setup hashing and creation error", func(pctx provider.StepCtx) {
			hasher.On("Hash", []byte(validPassword)).Return(hashedPassword, nil)
			repo.On("Create", ctx, mock.AnythingOfType("*auth.Member")).Return(nil, assert.AnError)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		t.WithNewStep("Attempt registration with creation error", func(pctx provider.StepCtx) {
			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})
}

func (s *AuthServiceTestSuite) TestLogin(t provider.T) {
	t.Tags("login", "positive")
	t.Description("Test user login functionality")
	t.Parallel()

	ctx := context.Background()
	validUserID := uuid.New()
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "securepassword"

	t.Run("Successful login with email", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		mockUser := new(authmock.MockUser)
		t.WithNewStep("Setup user authentication", func(pctx provider.StepCtx) {
			mockUser.On("ID").Return(validUserID)
			mockUser.On("Auth", []byte(validPassword), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(true)
		})

		t.WithNewStep("Setup repository and session storage", func(pctx provider.StepCtx) {
			repo.On("GetByEmail", ctx, validEmail).Return(mockUser, nil)
			sessions.On("Store", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*authservice.Session")).Return(nil)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		var sid uuid.UUID
		var err error
		t.WithNewStep("Perform login with email", func(pctx provider.StepCtx) {
			sid, err = svc.Login(ctx, validEmail, validPassword)
		})

		t.WithNewStep("Verify login success", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, sid)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			repo.AssertExpectations(t)
			sessions.AssertExpectations(t)
			mockUser.AssertExpectations(t)
		})
	})

	t.Run("Successful login with username", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		mockUser := new(authmock.MockUser)
		t.WithNewStep("Setup user authentication", func(pctx provider.StepCtx) {
			mockUser.On("ID").Return(validUserID)
			mockUser.On("Auth", []byte(validPassword), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(true)
		})

		t.WithNewStep("Setup repository fallback to username", func(pctx provider.StepCtx) {
			repo.On("GetByEmail", ctx, validName).Return(nil, assert.AnError)
			repo.On("GetByName", ctx, validName).Return(mockUser, nil)
			sessions.On("Store", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*authservice.Session")).Return(nil)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		var sid uuid.UUID
		var err error
		t.WithNewStep("Perform login with username", func(pctx provider.StepCtx) {
			sid, err = svc.Login(ctx, validName, validPassword)
		})

		t.WithNewStep("Verify login success", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, sid)
		})
	})

	t.Run("Invalid credentials during login", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		mockUser := new(authmock.MockUser)
		t.WithNewStep("Setup failed authentication", func(pctx provider.StepCtx) {
			mockUser.On("ID").Return(validUserID)
			mockUser.On("Auth", []byte("wrongpassword"), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(false)
		})

		t.WithNewStep("Setup repository lookup", func(pctx provider.StepCtx) {
			repo.On("GetByEmail", ctx, validEmail).Return(mockUser, nil)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		t.WithNewStep("Attempt login with wrong password", func(pctx provider.StepCtx) {
			_, err := svc.Login(ctx, validEmail, "wrongpassword")
			require.Error(t, err)
			assert.ErrorIs(t, err, authservice.ErrInvalidCredentials)
		})
	})

	t.Run("User not found during login", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		t.WithNewStep("Setup user not found", func(pctx provider.StepCtx) {
			repo.On("GetByEmail", ctx, validEmail).Return(nil, assert.AnError)
			repo.On("GetByName", ctx, validEmail).Return(nil, assert.AnError)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		t.WithNewStep("Attempt login with non-existent user", func(pctx provider.StepCtx) {
			_, err := svc.Login(ctx, validEmail, validPassword)
			require.Error(t, err)
		})
	})

	t.Run("Session storage error during login", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		mockUser := new(authmock.MockUser)
		t.WithNewStep("Setup successful authentication", func(pctx provider.StepCtx) {
			mockUser.On("ID").Return(validUserID)
			mockUser.On("Auth", []byte(validPassword), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(true)
		})

		t.WithNewStep("Setup session storage error", func(pctx provider.StepCtx) {
			repo.On("GetByEmail", ctx, validEmail).Return(mockUser, nil)
			sessions.On("Store", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*authservice.Session")).Return(assert.AnError)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		t.WithNewStep("Attempt login with session error", func(pctx provider.StepCtx) {
			_, err := svc.Login(ctx, validEmail, validPassword)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})
}

func (s *AuthServiceTestSuite) TestLogout(t provider.T) {
	t.Tags("logout", "positive")
	t.Description("Test user logout functionality")
	t.Parallel()

	validSessionID := uuid.New()
	validUserID := uuid.New()

	t.Run("Successful logout", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		t.WithNewStep("Setup session deletion", func(pctx provider.StepCtx) {
			sessions.On("Delete", mock.Anything, validSessionID).Return(nil)
		})

		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		t.WithNewStep("Setup session retrieval", func(pctx provider.StepCtx) {
			sessions.On("Get", mock.Anything, validSessionID).Return(validSession, nil)
		})

		var svc *authservice.AuthService
		ctx := context.Background()

		t.WithNewStep("Setup service and authenticate", func(pctx provider.StepCtx) {
			svc = authservice.NewAuthService(sessions, repo, hasher)
			ctx = svc.Authenticate(ctx, validSessionID)
		})

		t.WithNewStep("Perform logout", func(pctx provider.StepCtx) {
			err := svc.Logout(ctx)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			sessions.AssertExpectations(t)
		})
	})

	t.Run("Logout error", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		t.WithNewStep("Setup session deletion error", func(pctx provider.StepCtx) {
			sessions.On("Delete", mock.Anything, validSessionID).Return(assert.AnError)
		})

		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		t.WithNewStep("Setup session retrieval", func(pctx provider.StepCtx) {
			sessions.On("Get", mock.Anything, validSessionID).Return(validSession, nil)
		})

		var svc *authservice.AuthService
		ctx := context.Background()

		t.WithNewStep("Setup service and authenticate", func(pctx provider.StepCtx) {
			svc = authservice.NewAuthService(sessions, repo, hasher)
			ctx = svc.Authenticate(ctx, validSessionID)
		})

		t.WithNewStep("Attempt logout with error", func(pctx provider.StepCtx) {
			err := svc.Logout(ctx)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})
}

func (s *AuthServiceTestSuite) TestAuthenticate(t provider.T) {
	t.Tags("authentication", "positive")
	t.Description("Test user authentication functionality")
	t.Parallel()

	ctx := context.Background()
	validSessionID := uuid.New()
	validUserID := uuid.New()

	t.Run("Successful authentication", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		mockUser := new(authmock.MockUser)

		t.WithNewStep("Setup session retrieval", func(pctx provider.StepCtx) {
			sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		})

		t.WithNewStep("Setup user retrieval", func(pctx provider.StepCtx) {
			mockUser.On("ID").Return(validUserID)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		var newCtx context.Context
		t.WithNewStep("Perform authentication", func(pctx provider.StepCtx) {
			newCtx = svc.Authenticate(ctx, validSessionID)
		})

		t.WithNewStep("Setup user lookup in new context", func(pctx provider.StepCtx) {
			repo.On("Get", newCtx, validUserID).Return(mockUser, nil)
		})

		t.WithNewStep("Verify authenticated user", func(pctx provider.StepCtx) {
			getUser := svc.UserFromContext(newCtx)
			assert.NotNil(t, newCtx)
			assert.Equal(t, mockUser, getUser)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			sessions.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	})

	t.Run("Session not found during authentication", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		t.WithNewStep("Setup session not found", func(pctx provider.StepCtx) {
			sessions.On("Get", ctx, validSessionID).Return(nil, assert.AnError)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		var newCtx context.Context
		t.WithNewStep("Attempt authentication with missing session", func(pctx provider.StepCtx) {
			newCtx = svc.Authenticate(ctx, validSessionID)
		})

		t.WithNewStep("Verify no authentication", func(pctx provider.StepCtx) {
			user := svc.UserFromContext(newCtx)
			_, ok := user.(*auth.NoAuthUser)
			require.True(t, ok)
		})
	})

	t.Run("Session expired during authentication", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		expiredSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(-time.Hour),
		}

		t.WithNewStep("Setup expired session", func(pctx provider.StepCtx) {
			sessions.On("Get", ctx, validSessionID).Return(expiredSession, nil)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		var newCtx context.Context
		t.WithNewStep("Attempt authentication with expired session", func(pctx provider.StepCtx) {
			newCtx = svc.Authenticate(ctx, validSessionID)
		})

		t.WithNewStep("Verify no authentication", func(pctx provider.StepCtx) {
			user := svc.UserFromContext(newCtx)
			_, ok := user.(*auth.NoAuthUser)
			require.True(t, ok)
		})
	})

	t.Run("User not found during authentication", func(t provider.T) {
		t.Parallel()

		repo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)

		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		t.WithNewStep("Setup session retrieval", func(pctx provider.StepCtx) {
			sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		})

		svc := authservice.NewAuthService(sessions, repo, hasher)

		var newCtx context.Context
		t.WithNewStep("Perform authentication", func(pctx provider.StepCtx) {
			newCtx = svc.Authenticate(ctx, validSessionID)
		})

		t.WithNewStep("Setup user not found", func(pctx provider.StepCtx) {
			repo.On("Get", newCtx, validUserID).Return(nil, assert.AnError)
		})

		t.WithNewStep("Verify no authentication", func(pctx provider.StepCtx) {
			user := svc.UserFromContext(newCtx)
			_, ok := user.(*auth.NoAuthUser)
			require.True(t, ok)
		})
	})
}

func TestAuthService(t *testing.T) {
	suite.RunSuite(t, new(AuthServiceTestSuite))
}
