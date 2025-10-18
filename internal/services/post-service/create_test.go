//go:build unit

package postservice_test

import (
	"bytes"
	"context"
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

type PostServiceCreateTestSuite struct {
	suite.Suite
}

func (s *PostServiceCreateTestSuite) BeforeAll(t provider.T) {
	logs.InitNoopLogger()
}

func (s *PostServiceCreateTestSuite) BeforeEach(t provider.T) {
	t.Epic("Post Service")
	t.Feature("Post Management")
}

func (s *PostServiceCreateTestSuite) TestCreatePost(t provider.T) {
	t.Tags("create", "write")
	t.Description("Test post creation functionality")
	t.Parallel()

	validUserID := uuid.New()

	validContent, err := post.NewContent("Test content", post.ContentTypePlainText)
	require.NoError(t, err)

	validData := postservice.CreatePostTextData{
		Title:   "Test Post",
		Content: *validContent,
		Tags:    []string{"tag1", "tag2"},
	}

	validFiles := []models.FileData{
		{Name: "photo1.jpg", ContentType: "image/jpeg", Reader: bytes.NewReader([]byte("image data"))},
		{Name: "photo2.png", ContentType: "image/png", Reader: bytes.NewReader([]byte("image data"))},
	}

	validPhotoFiles := []*models.File{
		{ID: uuid.New(), Name: "photo1.jpg"},
		{ID: uuid.New(), Name: "photo2.png"},
	}

	t.Run("Successful post creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup file uploads", func(pctx provider.StepCtx) {
			for i, file := range validFiles {
				frepo.On("Upload", mock.Anything, &file).Return(validPhotoFiles[i], nil)
			}
		})

		t.WithNewStep("Setup post creation", func(pctx provider.StepCtx) {
			prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Run(func(args mock.Arguments) {
				p := args.Get(1).(*post.Post)
				assert.Equal(t, validData.Title, p.Title())
				assert.Equal(t, validData.Content, p.Content())
				assert.Equal(t, validData.Tags, p.Tags())
				assert.Equal(t, validUserID, p.AuthorID())
				assert.Equal(t, 2, p.Photos().Len())
			}).Return(&post.Post{}, nil)
		})

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Create post", func(pctx provider.StepCtx) {
			var err error
			_, err = svc.CreatePost(ctx, validData, validFiles)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			frepo.AssertExpectations(t)
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

		t.WithNewStep("Attempt create post without authorization", func(pctx provider.StepCtx) {
			_, err := svc.CreatePost(ctx, validData, validFiles)
			require.Error(t, err)
		})
	})

	t.Run("User without author rights", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, false)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt create post without author rights", func(pctx provider.StepCtx) {
			_, err := svc.CreatePost(ctx, validData, validFiles)
			require.Error(t, err)
			assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
		})
	})

	t.Run("Invalid file content type", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		invalidFiles := []models.FileData{
			{Name: "file.txt", ContentType: "text/plain", Reader: bytes.NewReader([]byte("text data"))},
		}

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt create post with invalid file type", func(pctx provider.StepCtx) {
			_, err := svc.CreatePost(ctx, validData, invalidFiles)
			require.Error(t, err)
			assert.ErrorIs(t, err, postservice.ErrInvalidFileContentType)
		})
	})

	t.Run("File upload error", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup file upload error", func(pctx provider.StepCtx) {
			frepo.On("Upload", mock.Anything, &validFiles[0]).Return(nil, assert.AnError)
		})

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt create post with file upload error", func(pctx provider.StepCtx) {
			_, err := svc.CreatePost(ctx, validData, validFiles)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Post creation error", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup successful file uploads", func(pctx provider.StepCtx) {
			for i, file := range validFiles {
				frepo.On("Upload", mock.Anything, &file).Return(validPhotoFiles[i], nil)
			}
		})

		t.WithNewStep("Setup post creation error", func(pctx provider.StepCtx) {
			prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Return(nil, assert.AnError)
		})

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Attempt create post with repository error", func(pctx provider.StepCtx) {
			_, err := svc.CreatePost(ctx, validData, validFiles)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Post creation with empty files", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validUserID, true)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup post creation without files", func(pctx provider.StepCtx) {
			prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Return(&post.Post{}, nil)
		})

		svc := postservice.NewPostService(prepo, frepo, asvc)

		t.WithNewStep("Create post without files", func(pctx provider.StepCtx) {
			result, err := svc.CreatePost(ctx, validData, []models.FileData{})
			require.NoError(t, err)
			assert.Equal(t, 0, result.Photos().Len())
		})
	})
}

func TestPostServiceCreate(t *testing.T) {
	suite.RunSuite(t, new(PostServiceCreateTestSuite))
}
