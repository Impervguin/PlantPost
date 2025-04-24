package postservice_test

import (
	"context"
	"testing"
	"time"

	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	postservice "PlantSite/internal/services/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPost(t *testing.T) {
	ctx := context.Background()
	validSessionID := uuid.New()
	validUserID := uuid.New()

	validContent, err := post.NewContent("Test content", post.ContentTypePlainText)
	require.NoError(t, err)

	photo1, err := post.NewPostPhoto(uuid.New(), 1)
	require.NoError(t, err)
	photo2, err := post.NewPostPhoto(uuid.New(), 2)
	require.NoError(t, err)

	photos := post.NewPostPhotos()
	err = photos.Add(photo1)
	require.NoError(t, err)
	err = photos.Add(photo2)
	require.NoError(t, err)

	validPost, err := post.NewPost(
		"Test Post",
		*validContent,
		[]string{"tag1", "tag2"},
		validUserID,
		photos,
	)
	require.NoError(t, err)

	validPostID := validPost.ID()
	createdAt := validPost.CreatedAt()
	updatedAt := validPost.UpdatedAt()

	photoFile1 := &models.File{ID: photo1.FileID(), Name: "photo1.jpg"}
	photoFile2 := &models.File{ID: photo2.FileID(), Name: "photo2.jpg"}

	expectedResult := &postservice.GetPost{
		ID:        validPost.ID(),
		Title:     "Test Post",
		Content:   *validContent,
		Tags:      []string{"tag1", "tag2"},
		AuthorID:  validUserID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Photos: []postservice.GetPostPhoto{
			{
				ID:          photo1.ID(),
				PlaceNumber: 1,
				File:        *photoFile1,
			},
			{
				ID:          photo2.ID(),
				PlaceNumber: 2,
				File:        *photoFile2,
			},
		},
	}

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

		prepo.On("Get", mock.Anything, validPostID).Return(validPost, nil)
		frepo.On("Get", mock.Anything, photo1.FileID()).Return(photoFile1, nil)
		frepo.On("Get", mock.Anything, photo2.FileID()).Return(photoFile2, nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		result, err := svc.GetPost(ctx, validPostID)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		prepo.AssertExpectations(t)
		frepo.AssertExpectations(t)
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

		_, err := svc.GetPost(ctx, validPostID)
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

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
	})

	t.Run("PostNotFound", func(t *testing.T) {
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

		prepo.On("Get", mock.Anything, validPostID).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("PhotoFileNotFound", func(t *testing.T) {
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

		prepo.On("Get", mock.Anything, validPostID).Return(validPost, nil)
		frepo.On("Get", mock.Anything, photo1.FileID()).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("NoPhotos", func(t *testing.T) {
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

		noPhotos := post.NewPostPhotos()
		postNoPhotos, err := post.NewPost(
			"Test Post",
			*validContent,
			[]string{"tag1", "tag2"},
			validUserID,
			noPhotos,
		)
		require.NoError(t, err)

		prepo.On("Get", mock.Anything, validPostID).Return(postNoPhotos, nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		result, err := svc.GetPost(ctx, validPostID)
		require.NoError(t, err)
		assert.Empty(t, result.Photos)
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

		_, err := svc.GetPost(ctx, uuid.Nil)
		require.Error(t, err)
	})
}
