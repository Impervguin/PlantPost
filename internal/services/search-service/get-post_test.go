package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	searchservice "PlantSite/internal/services/search-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	expectedResult := &searchservice.GetPost{
		ID:        validPost.ID(),
		Title:     "Test Post",
		Content:   *validContent,
		Tags:      []string{"tag1", "tag2"},
		AuthorID:  validUserID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Photos: []searchservice.GetPostPhoto{
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
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		// Setup expectations
		srepo.On("GetPostByID", ctx, validPostID).Return(validPost, nil)
		ptfrepo.On("Get", ctx, photo1.FileID()).Return(photoFile1, nil)
		ptfrepo.On("Get", ctx, photo2.FileID()).Return(photoFile2, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		result, err := svc.GetPost(ctx, validPostID)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		// Verify all expectations were met
		srepo.AssertExpectations(t)
		ptfrepo.AssertExpectations(t)
	})

	t.Run("PostNotFound", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srepo.On("GetPostByID", ctx, validPostID).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("PhotoFileNotFound", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srepo.On("GetPostByID", ctx, validPostID).Return(validPost, nil)
		ptfrepo.On("Get", ctx, photo1.FileID()).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.GetPost(ctx, validPostID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("NoPhotos", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		noPhotos := post.NewPostPhotos()
		postNoPhotos, err := post.NewPost(
			"Test Post",
			*validContent,
			[]string{"tag1", "tag2"},
			validUserID,
			noPhotos,
		)
		require.NoError(t, err)

		srepo.On("GetPostByID", ctx, validPostID).Return(postNoPhotos, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		result, err := svc.GetPost(ctx, validPostID)
		require.NoError(t, err)
		assert.Empty(t, result.Photos)
	})

	t.Run("NilPostID", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.GetPost(ctx, uuid.Nil)
		require.Error(t, err)
	})
}
