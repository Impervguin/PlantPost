package postservice_test

import (
	"context"
	"testing"
	"time"

	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	postservice "PlantSite/internal/services/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeletePost(t *testing.T) {
	validSessionID := uuid.New()
	validUserID := uuid.New()
	ctx := context.Background()
	validPostID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		prepo.On("Delete", mock.Anything, validPostID).Return(nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		err := svc.Delete(ctx, validPostID)
		require.NoError(t, err)

		// Verify all expectations were met
		prepo.AssertExpectations(t)
	})

	t.Run("NotAuthorized", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		user := new(authmock.MockUser)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(nil, assert.AnError)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		err := svc.Delete(ctx, validPostID)
		require.Error(t, err)
	})

	t.Run("NotAuthor", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("HasAuthorRights").Return(false)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		err := svc.Delete(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthor)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		prepo.On("Delete", mock.Anything, validPostID).Return(assert.AnError)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		err := svc.Delete(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("NilPostID", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validUserID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		err := svc.Delete(ctx, uuid.Nil)
		require.Error(t, err)
	})
}
