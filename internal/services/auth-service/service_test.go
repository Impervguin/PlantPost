package authservice_test

import (
	"context"
	"testing"

	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService(t *testing.T) {
	ctx := context.Background()
	validUserID := uuid.New()
	// validSessionID := uuid.New()
	validName := "testuser"
	validEmail := "test@example.com"
	validPassword := "securepassword"
	hashedPassword := []byte("hashedpassword")

	t.Run("Register", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

			hasher.On("Hash", []byte(validPassword)).Return(hashedPassword, nil)

			mockUser := new(authmock.MockUser)
			repo.On("Create", ctx, mock.AnythingOfType("*auth.Member")).Return(mockUser, nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.NoError(t, err)

			hasher.AssertExpectations(t)
			repo.AssertExpectations(t)
		})

		t.Run("HashError", func(t *testing.T) {
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

			hasher.On("Hash", []byte(validPassword)).Return(nil, assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			err := svc.Register(ctx, validName, validEmail, validPassword)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("CreateUserError", func(t *testing.T) {
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

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
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

			mockUser := new(authmock.MockUser)
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
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

			mockUser := new(authmock.MockUser)
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
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

			mockUser := new(authmock.MockUser)
			mockUser.On("Auth", []byte("wrongpassword"), mock.AnythingOfType("func([]uint8, []uint8) (bool, error)")).Return(false)

			repo.On("GetByEmail", ctx, validEmail).Return(mockUser, nil)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Login(ctx, validEmail, "wrongpassword")
			require.Error(t, err)
			assert.ErrorIs(t, err, authservice.ErrInvalidCredentials)
		})

		t.Run("UserNotFound", func(t *testing.T) {
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

			repo.On("GetByEmail", ctx, validEmail).Return(nil, assert.AnError)
			repo.On("GetByName", ctx, validEmail).Return(nil, assert.AnError)

			svc := authservice.NewAuthService(sessions, repo, hasher)

			_, err := svc.Login(ctx, validEmail, validPassword)
			require.Error(t, err)
		})

		t.Run("SessionStoreError", func(t *testing.T) {
			repo := new(authmock.MockAuthRepository)
			sessions := new(authmock.MockSessionStorage)
			hasher := new(authmock.MockPasswdHasher)

			mockUser := new(authmock.MockUser)
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

	// t.Run("Logout", func(t *testing.T) {
	// 	t.Run("Success", func(t *testing.T) {
	// 		repo := new(authmock.MockAuthRepository)
	// 		sessions := new(authmock.MockSessionStorage)
	// 		hasher := new(authmock.MockPasswdHasher)

	// 		sessions.On("Delete", mock.Anything, validSessionID).Return(nil)
	// 		validSession := &authservice.Session{
	// 			ID:        validSessionID,
	// 			MemberID:  validUserID,
	// 			ExpiresAt: time.Now().Add(time.Hour),
	// 		}
	// 		sessions.On("Get", mock.Anything, validSessionID).Return(validSession, nil)

	// 		svc := authservice.NewAuthService(sessions, repo, hasher)

	// 		err := svc.Logout(ctx)
	// 		require.NoError(t, err)

	// 		sessions.AssertExpectations(t)
	// 	})

	// 	t.Run("Error", func(t *testing.T) {
	// 		repo := new(authmock.MockAuthRepository)
	// 		sessions := new(authmock.MockSessionStorage)
	// 		hasher := new(authmock.MockPasswdHasher)

	// 		sessions.On("Delete", mock.Anything, validSessionID).Return(assert.AnError)
	// 		validSession := &authservice.Session{
	// 			ID:        validSessionID,
	// 			MemberID:  validUserID,
	// 			ExpiresAt: time.Now().Add(time.Hour),
	// 		}
	// 		sessions.On("Get", mock.Anything, validSessionID).Return(validSession, nil)

	// 		svc := authservice.NewAuthService(sessions, repo, hasher)

	// 		err := svc.Logout(ctx)
	// 		require.Error(t, err)
	// 		assert.ErrorIs(t, err, assert.AnError)
	// 	})
	// })

	// t.Run("Authenticate", func(t *testing.T) {
	// 	t.Run("Success", func(t *testing.T) {
	// 		ctx := context.Background()
	// 		repo := new(authmock.MockAuthRepository)
	// 		sessions := new(authmock.MockSessionStorage)
	// 		hasher := new(authmock.MockPasswdHasher)

	// 		validSession := &authservice.Session{
	// 			ID:        validSessionID,
	// 			MemberID:  validUserID,
	// 			ExpiresAt: time.Now().Add(time.Hour),
	// 		}

	// 		mockUser := new(authmock.MockUser)

	// 		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)

	// 		mockUser.On("ID").Return(validUserID)

	// 		svc := authservice.NewAuthService(sessions, repo, hasher)

	// 		newCtx := svc.Authenticate(ctx, validSessionID)
	// 		repo.On("Get", newCtx, validUserID).Return(mockUser, nil)
	// 		getUser := svc.UserFromContext(newCtx)
	// 		assert.NotNil(t, newCtx)
	// 		assert.Equal(t, mockUser, getUser)

	// 		sessions.AssertExpectations(t)
	// 		repo.AssertExpectations(t)
	// 	})

	// 	t.Run("SessionNotFound", func(t *testing.T) {
	// 		repo := new(authmock.MockAuthRepository)
	// 		sessions := new(authmock.MockSessionStorage)
	// 		hasher := new(authmock.MockPasswdHasher)

	// 		sessions.On("Get", ctx, validSessionID).Return(nil, assert.AnError)

	// 		svc := authservice.NewAuthService(sessions, repo, hasher)
	// 		ctx := svc.Authenticate(ctx, validSessionID)

	// 		user := svc.UserFromContext(ctx)
	// 		_, ok := user.(*auth.NoAuthUser)
	// 		require.True(t, ok)
	// 	})

	// 	t.Run("SessionExpired", func(t *testing.T) {
	// 		repo := new(authmock.MockAuthRepository)
	// 		sessions := new(authmock.MockSessionStorage)
	// 		hasher := new(authmock.MockPasswdHasher)

	// 		expiredSession := &authservice.Session{
	// 			ID:        validSessionID,
	// 			MemberID:  validUserID,
	// 			ExpiresAt: time.Now().Add(-time.Hour),
	// 		}

	// 		sessions.On("Get", ctx, validSessionID).Return(expiredSession, nil)

	// 		svc := authservice.NewAuthService(sessions, repo, hasher)

	// 		ctx := svc.Authenticate(ctx, validSessionID)
	// 		user := svc.UserFromContext(ctx)
	// 		_, ok := user.(*auth.NoAuthUser)
	// 		require.True(t, ok)
	// 	})

	// 	t.Run("UserNotFound", func(t *testing.T) {
	// 		repo := new(authmock.MockAuthRepository)
	// 		sessions := new(authmock.MockSessionStorage)
	// 		hasher := new(authmock.MockPasswdHasher)

	// 		validSession := &authservice.Session{
	// 			ID:        validSessionID,
	// 			MemberID:  validUserID,
	// 			ExpiresAt: time.Now().Add(time.Hour),
	// 		}

	// 		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)

	// 		svc := authservice.NewAuthService(sessions, repo, hasher)

	// 		ctx := svc.Authenticate(ctx, validSessionID)
	// 		repo.On("Get", ctx, validUserID).Return(nil, assert.AnError)
	// 		user := svc.UserFromContext(ctx)
	// 		_, ok := user.(*auth.NoAuthUser)
	// 		require.True(t, ok)
	// 	})
	// })
}
