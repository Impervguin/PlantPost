package postservice_test

import (
	"context"
	"testing"
	"time"

	"PlantSite/internal/models/post"
	postservice "PlantSite/internal/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdatePost(t *testing.T) {
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
	photos.Add(photo1)
	photos.Add(photo2)

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
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		user.On("ID").Return(validUserID)
		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		result, err := svc.UpdatePost(ctx, validPostID, updateData)
		require.NoError(t, err)
		assert.Equal(t, result.Tags(), updatedPost.Tags())
		assert.Equal(t, result.Title(), updatedPost.Title())
		assert.Equal(t, result.Content(), updatedPost.Content())
		assert.NotEqual(t, result.UpdatedAt(), updatedAt)

		user.AssertExpectations(t)
		prepo.AssertExpectations(t)
	})

	t.Run("NotAuthorized", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo)

		_, err := svc.UpdatePost(ctx, validPostID, updateData)
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

		_, err := svc.UpdatePost(ctx, validPostID, updateData)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthor)
	})

	t.Run("PostNotFound", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		user.On("ID").Return(validUserID)
		prepo.On("Update", mock.Anything, validPostID, mock.Anything).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.UpdatePost(ctx, validPostID, updateData)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("InvalidContent", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		invalidData := updateData
		invalidData.Content = post.Content{Text: "", ContentType: "invalid_type"}

		user.On("HasAuthorRights").Return(true)
		user.On("ID").Return(validUserID)
		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.UpdatePost(ctx, validPostID, invalidData)
		require.Error(t, err)
	})

	t.Run("EmptyTitle", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		invalidData := updateData
		invalidData.Title = ""

		user.On("HasAuthorRights").Return(true)
		user.On("ID").Return(validUserID)
		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.UpdatePost(ctx, validPostID, invalidData)
		require.Error(t, err)
	})

	t.Run("NilTags", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		invalidData := updateData
		invalidData.Tags = nil

		user.On("HasAuthorRights").Return(true)
		user.On("ID").Return(validUserID)
		prepo.On("Update", mock.Anything, validPostID, mock.AnythingOfType("func(*post.Post) (*post.Post, error)")).
			Return(validPost, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.UpdatePost(ctx, validPostID, invalidData)
		require.Error(t, err)
	})
}
