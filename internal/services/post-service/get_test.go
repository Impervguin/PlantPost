package postservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	postservice "PlantSite/internal/services/post-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPost(t *testing.T) {
	ctx := context.Background()
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
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		prepo.On("Get", mock.Anything, validPostID).Return(validPost, nil)
		frepo.On("Get", mock.Anything, photo1.FileID()).Return(photoFile1, nil)
		frepo.On("Get", mock.Anything, photo2.FileID()).Return(photoFile2, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		result, err := svc.GetPost(ctx, validPostID)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		user.AssertExpectations(t)
		prepo.AssertExpectations(t)
		frepo.AssertExpectations(t)
	})

	t.Run("NotAuthorized", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)

		svc := postservice.NewPostService(prepo, frepo)

		_, err := svc.GetPost(ctx, validPostID)
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

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, postservice.ErrNotAuthor)
	})

	t.Run("PostNotFound", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		prepo.On("Get", mock.Anything, validPostID).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("PhotoFileNotFound", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		prepo.On("Get", mock.Anything, validPostID).Return(validPost, nil)
		frepo.On("Get", mock.Anything, photo1.FileID()).Return(nil, assert.AnError)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("NoPhotos", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		noPhotos := post.NewPostPhotos()
		postNoPhotos, err := post.NewPost(
			"Test Post",
			*validContent,
			[]string{"tag1", "tag2"},
			validUserID,
			noPhotos,
		)
		require.NoError(t, err)

		user.On("HasAuthorRights").Return(true)
		prepo.On("Get", mock.Anything, validPostID).Return(postNoPhotos, nil)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		result, err := svc.GetPost(ctx, validPostID)
		require.NoError(t, err)
		assert.Empty(t, result.Photos)
	})

	t.Run("NilPostID", func(t *testing.T) {
		prepo := new(MockPostRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)

		svc := postservice.NewPostService(prepo, frepo)
		ctx := PutUserInContext(ctx, user)

		_, err := svc.GetPost(ctx, uuid.Nil)
		require.Error(t, err)
	})
}
