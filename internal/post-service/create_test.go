package postservice_test

import (
	"bytes"
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	postservice "PlantSite/internal/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	ctx := context.Background()
	validUserID := uuid.New()

	// Create test data
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

	validPost, err := post.NewPost(
		validData.Title,
		validData.Content,
		validData.Tags,
		validUserID,
		post.NewPostPhotos(), // Will be filled in tests
	)
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		// Setup expectations
		user.On("ID").Return(validUserID)
		user.On("HasAuthorRights").Return(true)

		// Expect file uploads
		for i, file := range validFiles {
			frepo.On("Upload", mock.Anything, &file).Return(validPhotoFiles[i], nil)
		}

		// Expect post creation with any post that has our photos
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Run(func(args mock.Arguments) {
			p := args.Get(1).(*post.Post)
			assert.Equal(t, validData.Title, p.Title())
			assert.Equal(t, validData.Content, p.Content())
			assert.Equal(t, validData.Tags, p.Tags())
			assert.Equal(t, validUserID, p.AuthorID())
			assert.Equal(t, 2, p.Photos().Len())
		}).Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		result, err := svc.CreatePost(ctx, validData, validFiles)
		require.NoError(t, err)
		assert.Equal(t, validPost, result)

		// Verify all expectations were met
		user.AssertExpectations(t)
		frepo.AssertExpectations(t)
		prepo.AssertExpectations(t)
	})

	t.Run("NotAuthorized", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo)

		_, err := svc.CreatePost(ctx, validData, validFiles)
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

		_, err := svc.CreatePost(ctx, validData, validFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthor)
	})

	t.Run("InvalidFileContentType", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		invalidFiles := []models.FileData{
			{Name: "file.txt", ContentType: "text/plain", Reader: bytes.NewReader([]byte("text data"))},
		}

		user.On("HasAuthorRights").Return(true)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.CreatePost(ctx, validData, invalidFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrInvalidFileContentType)
	})

	t.Run("FileUploadError", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("ID").Return(validUserID)
		user.On("HasAuthorRights").Return(true)
		frepo.On("Upload", mock.Anything, &validFiles[0]).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.CreatePost(ctx, validData, validFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("PostCreationError", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("ID").Return(validUserID)
		user.On("HasAuthorRights").Return(true)
		for i, file := range validFiles {
			frepo.On("Upload", mock.Anything, &file).Return(validPhotoFiles[i], nil)
		}
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.CreatePost(ctx, validData, validFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("EmptyFiles", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("ID").Return(validUserID)
		user.On("HasAuthorRights").Return(true)
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		result, err := svc.CreatePost(ctx, validData, []models.FileData{})
		require.NoError(t, err)
		assert.Equal(t, validPost, result)
		assert.Equal(t, 0, result.Photos().Len())
	})
}
