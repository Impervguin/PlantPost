package postservice_test

import (
	"context"
	"testing"
	"time"

	"PlantSite/internal/models/post"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	postservice "PlantSite/internal/services/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdatePost(t *testing.T) {
	validSessionID := uuid.New()
	ctx := context.Background()
	validUserID := uuid.New()

	// Create test data
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

	validPost, err := post.CreatePost(
		uuid.New(),
		"Test Post",
		*validContent,
		[]string{"tag1", "tag2"},
		validUserID,
		*photos,
		time.Now().Add(-time.Hour),
		time.Now().Add(-time.Hour).Add(time.Minute),
	)
	require.NoError(t, err)

	validPostID := validPost.ID()
	createdAt := validPost.CreatedAt()
	updatedAt := validPost.UpdatedAt()

	newContent, err := post.NewContent("updated content", post.ContentTypePlainText)
	require.NoError(t, err)

	updatedPost, err := post.CreatePost(
		validPost.ID(),
		"Updated Title",
		*newContent,
		[]string{"newtag1", "newtag2"},
		validUserID,
		*photos,
		createdAt,
		updatedAt.Add(time.Minute),
	)
	require.NoError(t, err)
	updateData := postservice.UpdatePostTextData{
		Title:   "Updated Title",
		Content: *newContent,
		Tags:    []string{"newtag1", "newtag2"},
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
		user.On("ID").Return(validUserID)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validUserID).Return(user, nil)

		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		result, err := svc.UpdatePost(ctx, validPostID, updateData)
		require.NoError(t, err)
		assert.Equal(t, result.Tags(), updatedPost.Tags())
		assert.Equal(t, result.Title(), updatedPost.Title())
		assert.Equal(t, result.Content(), updatedPost.Content())
		assert.NotEqual(t, result.UpdatedAt(), updatedAt)

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

		_, err := svc.UpdatePost(ctx, validPostID, updateData)
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

		_, err := svc.UpdatePost(ctx, validPostID, updateData)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthor)
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

		prepo.On("Update", mock.Anything, validPostID, mock.Anything).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.UpdatePost(ctx, validPostID, updateData)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("InvalidContent", func(t *testing.T) {
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

		invalidData := updateData
		invalidData.Content = post.Content{Text: "", ContentType: "invalid_type"}

		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.UpdatePost(ctx, validPostID, invalidData)
		require.Error(t, err)
	})

	t.Run("EmptyTitle", func(t *testing.T) {
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

		invalidData := updateData
		invalidData.Title = ""

		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.UpdatePost(ctx, validPostID, invalidData)
		require.Error(t, err)
	})

	t.Run("NilTags", func(t *testing.T) {
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

		invalidData := updateData
		invalidData.Tags = nil

		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo, asvc)

		_, err := svc.UpdatePost(ctx, validPostID, invalidData)
		require.Error(t, err)
	})
}
