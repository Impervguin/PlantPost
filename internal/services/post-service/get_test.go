//go:build unit

package postservice_test

import (
	"context"
	"fmt"
	"testing"

	"PlantSite/internal/models"
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

type PostServiceTestSuite struct {
	suite.Suite
}

func (s *PostServiceTestSuite) BeforeAll(t provider.T) {
	logs.InitNoopLogger()
}

func (s *PostServiceTestSuite) BeforeEach(t provider.T) {
	t.Epic("Post Service")
	t.Feature("Post Management")
}

type GetPostResultBuilder struct {
	post       *post.Post
	photoFiles map[uuid.UUID]*models.File
}

func NewGetPostResultBuilder(post *post.Post) *GetPostResultBuilder {
	return &GetPostResultBuilder{
		post:       post,
		photoFiles: make(map[uuid.UUID]*models.File),
	}
}

func (b *GetPostResultBuilder) WithPhotoFile(fileID uuid.UUID, fileName string) *GetPostResultBuilder {
	b.photoFiles[fileID] = &models.File{ID: fileID, Name: fileName}
	return b
}

func (b *GetPostResultBuilder) Build() *postservice.GetPost {
	photos := make([]postservice.GetPostPhoto, 0)
	postPhotos := b.post.Photos()
	postPhotos.Iterate(func(e post.PostPhoto) error {
		if file, exists := b.photoFiles[e.FileID()]; exists {
			photos = append(photos, postservice.GetPostPhoto{
				ID:          e.ID(),
				PlaceNumber: e.PlaceNumber(),
				File:        *file,
			})
		}
		return nil
	})

	return &postservice.GetPost{
		ID:        b.post.ID(),
		Title:     b.post.Title(),
		Content:   b.post.Content(),
		Tags:      b.post.Tags(),
		AuthorID:  b.post.AuthorID(),
		CreatedAt: b.post.CreatedAt(),
		UpdatedAt: b.post.UpdatedAt(),
		Photos:    photos,
	}
}

func (s *PostServiceTestSuite) TestGetPost(t provider.T) {
	t.Tags("get", "read")
	t.Description("Test post retrieval functionality")
	t.Parallel()

	validUserID := uuid.New()

	postBuilder := NewPostBuilder().
		WithAuthorID(validUserID).
		WithPhoto(uuid.New(), 1).
		WithPhoto(uuid.New(), 2)

	validPost, err := postBuilder.Build()
	require.NoError(t, err)
	validPostID := validPost.ID()
	resultBuilder := NewGetPostResultBuilder(validPost)
	for _, photo := range validPost.Photos().List() {
		resultBuilder.WithPhotoFile(photo.FileID(), fmt.Sprintf("photo%d.jpg", photo.PlaceNumber()))
	}
	expectedResult := resultBuilder.Build()

	t.Run("Successful post retrieval", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post repository", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPostID).Return(validPost, nil)
		})

		frepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			for _, photo := range expectedResult.Photos {
				frepo.On("Get", mock.Anything, photo.File.ID).Return(&photo.File, nil)
			}
		})

		svc := postservice.NewPostService(prepo, frepo, asvc)

		var result *postservice.GetPost
		t.WithNewStep("Get post", func(pctx provider.StepCtx) {
			var err error
			result, err = svc.GetPost(ctx, validPostID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify result", func(pctx provider.StepCtx) {
			assert.Equal(t, expectedResult, result)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			prepo.AssertExpectations(t)
			frepo.AssertExpectations(t)
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

		t.WithNewStep("Attempt get post without authorization", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, validPostID)
			require.Error(t, err)
		})
	})

	t.Run("User without author rights", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, false)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt get post without author rights", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, validPostID)
			require.Error(t, err)
			assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
		})
	})

	t.Run("Post not found", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post not found", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPostID).Return(nil, assert.AnError)
		})

		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt get non-existent post", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, validPostID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Photo file not found", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post repository", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPostID).Return(validPost, nil)
		})

		frepo := new(MockFileRepository)
		t.WithNewStep("Setup file not found", func(pctx provider.StepCtx) {
			for _, photo := range expectedResult.Photos {
				frepo.On("Get", mock.Anything, photo.File.ID).Return(nil, assert.AnError)
			}
		})

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt get post with missing photo file", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, validPostID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Post with no photos", func(t provider.T) {
		t.Parallel()

		noPhotosPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			WithNoPhotos().
			Build()
		require.NoError(t, err)

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		t.WithNewStep("Setup post with no photos", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, noPhotosPost.ID()).Return(noPhotosPost, nil)
		})

		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		var result *postservice.GetPost
		t.WithNewStep("Get post with no photos", func(pctx provider.StepCtx) {
			var err error
			result, err = svc.GetPost(ctx, noPhotosPost.ID())
			require.NoError(t, err)
		})

		t.WithNewStep("Verify empty photos", func(pctx provider.StepCtx) {
			assert.Empty(t, result.Photos)
		})
	})

	t.Run("Nil post ID", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt get post with nil ID", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, uuid.Nil)
			require.Error(t, err)
		})
	})
}

func TestPostService(t *testing.T) {
	suite.RunSuite(t, new(PostServiceTestSuite))
}
