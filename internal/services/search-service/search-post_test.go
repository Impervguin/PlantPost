package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/services/search-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSearchPosts(t *testing.T) {
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

	createdAt := validPost.CreatedAt()
	updatedAt := validPost.UpdatedAt()

	photoFile1 := &models.File{ID: photo1.FileID(), Name: "photo1.jpg"}
	photoFile2 := &models.File{ID: photo2.FileID(), Name: "photo2.jpg"}

	expectedResult := &searchservice.SearchPost{
		ID:        validPost.ID(),
		Title:     "Test Post",
		Content:   *validContent,
		Tags:      []string{"tag1", "tag2"},
		AuthorID:  validUserID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Photos: []searchservice.SearchPostPhoto{
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

		srch := search.NewPostSearch()
		srch.AddFilter(search.NewPostAuthorFilter(validUserID))
		srepo.On("SearchPosts", mock.Anything, srch).Return([]*post.Post{validPost}, nil)
		ptfrepo.On("Get", mock.Anything, photo1.FileID()).Return(photoFile1, nil)
		ptfrepo.On("Get", mock.Anything, photo2.FileID()).Return(photoFile2, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPosts(ctx, srch)
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, expectedResult, results[0])

		srepo.AssertExpectations(t)
		ptfrepo.AssertExpectations(t)
	})

	t.Run("EmptyResults", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srch := search.NewPostSearch()
		srch.AddFilter(search.NewPostAuthorFilter(uuid.New()))
		srepo.On("SearchPosts", mock.Anything, srch).Return([]*post.Post{}, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPosts(ctx, srch)
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srch := search.NewPostSearch()
		srepo.On("SearchPosts", ctx, srch).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.SearchPosts(ctx, srch)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("PhotoFileNotFound", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srch := search.NewPostSearch()
		srepo.On("SearchPosts", ctx, srch).Return([]*post.Post{validPost}, nil)
		ptfrepo.On("Get", ctx, photo1.FileID()).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.SearchPosts(ctx, srch)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
