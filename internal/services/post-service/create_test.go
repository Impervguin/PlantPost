package postservice_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	postservice "PlantSite/internal/services/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	validSessionID := uuid.New()
	validUserID := uuid.New()
	ctx := context.Background()

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
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

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

		svc := postservice.NewPostService(prepo, frepo, asvc)

		result, err := svc.CreatePost(ctx, validData, validFiles)
		require.NoError(t, err)
		assert.Equal(t, validPost, result)

		// Verify all expectations were met
		frepo.AssertExpectations(t)
		prepo.AssertExpectations(t)
	})

	t.Run("NotAuthorized", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		user := new(authmock.MockUser)
		user.On("HasAuthorRights").Return(true)
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(nil, assert.AnError)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.CreatePost(ctx, validData, validFiles)
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
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.CreatePost(ctx, validData, validFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthor)
	})

	t.Run("InvalidFileContentType", func(t *testing.T) {
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
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		invalidFiles := []models.FileData{
			{Name: "file.txt", ContentType: "text/plain", Reader: bytes.NewReader([]byte("text data"))},
		}

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.CreatePost(ctx, validData, invalidFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrInvalidFileContentType)
	})

	t.Run("FileUploadError", func(t *testing.T) {
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
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		frepo.On("Upload", mock.Anything, &validFiles[0]).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.CreatePost(ctx, validData, validFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("PostCreationError", func(t *testing.T) {
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
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		for i, file := range validFiles {
			frepo.On("Upload", mock.Anything, &file).Return(validPhotoFiles[i], nil)
		}
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.CreatePost(ctx, validData, validFiles)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("EmptyFiles", func(t *testing.T) {
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
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		prepo.On("Create", mock.Anything, mock.AnythingOfType("*post.Post")).Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		result, err := svc.CreatePost(ctx, validData, []models.FileData{})
		require.NoError(t, err)
		assert.Equal(t, validPost, result)
		assert.Equal(t, 0, result.Photos().Len())
	})
}
