package postservice_test

import (
	"context"
	"testing"

	postservice "PlantSite/internal/services/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeletePost(t *testing.T) {
	ctx := context.Background()
	validPostID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		// Setup expectations
		user.On("HasAuthorRights").Return(true)
		prepo.On("Delete", mock.Anything, validPostID).Return(nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.Delete(ctx, validPostID)
		require.NoError(t, err)

		// Verify all expectations were met
		user.AssertExpectations(t)
		prepo.AssertExpectations(t)
	})

	t.Run("NotAuthorized", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo)

		err := svc.Delete(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthorized)
	})

	t.Run("NotAuthor", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(false)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.Delete(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthor)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		prepo.On("Delete", mock.Anything, validPostID).Return(assert.AnError)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.Delete(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("NilPostID", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.Delete(ctx, uuid.Nil)
		require.Error(t, err)
	})
}
