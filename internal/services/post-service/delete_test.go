//go:build unit

package postservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models/auth"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	postservice "PlantSite/internal/services/post-service"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type PostServiceDeleteTestSuite struct {
	suite.Suite
}

func (s *PostServiceDeleteTestSuite) BeforeEach(t provider.T) {
	t.Epic("Post Service")
	t.Feature("Post Management")
}

func (s *PostServiceDeleteTestSuite) TestDeletePost(t provider.T) {
	t.Tags("delete", "write")
	t.Description("Test post deletion functionality")
	t.Parallel()

	validUserID := uuid.New()
	validPostID := uuid.New()

	t.Run("Successful post deletion", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post deletion", func(pctx provider.StepCtx) {
			prepo.On("Delete", mock.Anything, validPostID).Return(nil)
		})

		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Delete post", func(pctx provider.StepCtx) {
			err := svc.Delete(ctx, validPostID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			prepo.AssertExpectations(t)
		})
	})

	t.Run("Not authorized user", func(t provider.T) {
		t.Parallel()

		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)

		user := new(authmock.MockUser)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		ctx := asvc.Authenticate(context.Background(), uuid.New())

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt delete post without authorization", func(pctx provider.StepCtx) {
			err := svc.Delete(ctx, validPostID)
			require.Error(t, err)
		})
	})

	t.Run("User without author rights", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, false)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt delete post without author rights", func(pctx provider.StepCtx) {
			err := svc.Delete(ctx, validPostID)
			require.Error(t, err)
			assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
		})
	})

	t.Run("Repository error during deletion", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup repository error", func(pctx provider.StepCtx) {
			prepo.On("Delete", mock.Anything, validPostID).Return(assert.AnError)
		})

		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt delete post with repository error", func(pctx provider.StepCtx) {
			err := svc.Delete(ctx, validPostID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Nil post ID", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt delete post with nil ID", func(pctx provider.StepCtx) {
			err := svc.Delete(ctx, uuid.Nil)
			require.Error(t, err)
		})
	})
}

func TestPostDeleteService(t *testing.T) {
	suite.RunSuite(t, new(PostServiceDeleteTestSuite))
}
