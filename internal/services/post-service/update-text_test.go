//go:build unit

package postservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	postservice "PlantSite/internal/services/post-service"
	"PlantSite/internal/utils/logs"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type PostServiceUpdateTextTestSuite struct {
	suite.Suite
}

func (s *PostServiceUpdateTextTestSuite) BeforeAll(t provider.T) {
	logs.InitNoopLogger()
}

func (s *PostServiceUpdateTextTestSuite) BeforeEach(t provider.T) {
	t.Epic("Post Service")
	t.Feature("Post Management")
}

func (s *PostServiceUpdateTextTestSuite) TestUpdatePost(t provider.T) {
	t.Tags("update", "write")
	t.Description("Test post update functionality")
	t.Parallel()

	validUserID := uuid.New()

	newContent, err := post.NewContent("updated content", post.ContentTypePlainText)
	require.NoError(t, err)

	updateData := postservice.UpdatePostTextData{
		Title:   "Updated Title",
		Content: *newContent,
		Tags:    []string{"newtag1", "newtag2"},
	}

	t.Run("Successful post update", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			Build()
		beforeUpdateTime := validPost.UpdatedAt()
		require.NoError(t, err)
		validPostID := validPost.ID()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post update", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).Return(validPost, nil).Run(func(args mock.Arguments) {
				fn := args.Get(2).(func(*post.Post) (*post.Post, error))
				fn(validPost)
			})
		})

		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		var result *post.Post
		t.WithNewStep("Update post", func(pctx provider.StepCtx) {
			var err error
			result, err = svc.UpdatePost(ctx, validPostID, updateData)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify updated fields", func(pctx provider.StepCtx) {
			assert.Equal(t, result.Tags(), updateData.Tags)
			assert.Equal(t, result.Title(), updateData.Title)
			assert.Equal(t, result.Content(), updateData.Content)
			assert.NotEqual(t, result.UpdatedAt(), beforeUpdateTime)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			prepo.AssertExpectations(t)
		})
	})

	t.Run("Not authorized user", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

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

		t.WithNewStep("Attempt update post without authorization", func(pctx provider.StepCtx) {
			_, err := svc.UpdatePost(ctx, validPostID, updateData)
			require.Error(t, err)
		})
	})

	t.Run("User without author rights", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

		asvc, ctx := setupAuthService(t, validUserID, false)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt update post without author rights", func(pctx provider.StepCtx) {
			_, err := svc.UpdatePost(ctx, validPostID, updateData)
			require.Error(t, err)
			assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
		})
	})

	t.Run("Post not found", func(t provider.T) {
		t.Parallel()

		validPostID := uuid.New()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post not found", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPostID, mock.Anything).Return(nil, assert.AnError)
		})

		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt update non-existent post", func(pctx provider.StepCtx) {
			_, err := svc.UpdatePost(ctx, validPostID, updateData)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Invalid content", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post update", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
				Return(validPost, nil)
		})

		frepo := new(MockFileRepository)

		invalidData := updateData
		invalidData.Content = post.Content{Text: "", ContentType: "invalid_type"}

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt update with invalid content", func(pctx provider.StepCtx) {
			_, err := svc.UpdatePost(ctx, validPostID, invalidData)
			require.Error(t, err)
		})
	})

	t.Run("Empty title", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post update", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
				Return(validPost, nil)
		})

		frepo := new(MockFileRepository)

		invalidData := updateData
		invalidData.Title = ""

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt update with empty title", func(pctx provider.StepCtx) {
			_, err := svc.UpdatePost(ctx, validPostID, invalidData)
			require.Error(t, err)
		})
	})

	t.Run("Nil tags", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post update", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
				Return(validPost, nil)
		})

		frepo := new(MockFileRepository)

		invalidData := updateData
		invalidData.Tags = nil

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt update with nil tags", func(pctx provider.StepCtx) {
			_, err := svc.UpdatePost(ctx, validPostID, invalidData)
			require.Error(t, err)
		})
	})
}

func TestPostUpdateTextService(t *testing.T) {
	suite.RunSuite(t, new(PostServiceUpdateTextTestSuite))
}
