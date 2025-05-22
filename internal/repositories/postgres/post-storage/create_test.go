//go:build integration

package poststorage_test

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PostRepositoryTestSuite) TestCreatePost() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Test creation
	createdPost, err := s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)

	// Verify returned post matches input
	assert.Equal(s.T(), testPost.ID(), createdPost.ID())
	assert.Equal(s.T(), testPost.Title(), createdPost.Title())
	assert.Equal(s.T(), testPost.Content().Text, createdPost.Content().Text)
	assert.Equal(s.T(), testPost.AuthorID(), createdPost.AuthorID())
	assert.Equal(s.T(), testPost.Photos().Len(), createdPost.Photos().Len())
	assert.ElementsMatch(s.T(), testPost.Tags(), createdPost.Tags())

	// Verify can retrieve
	fetchedPost, err := s.repo.Get(ctx, testPost.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), testPost.ID(), fetchedPost.ID())
}

func (s *PostRepositoryTestSuite) TestCreatePostWithNoPhotos() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Remove all photos
	testPost.ClearPhotos()

	createdPost, err := s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, createdPost.Photos().Len())
}

func (s *PostRepositoryTestSuite) TestCreatePostWithInvalidPhoto() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Delete a photo from storage
	photos := testPost.Photos().List()
	if len(photos) > 0 {
		err := s.fileRepo.Delete(ctx, photos[0].FileID())
		require.NoError(s.T(), err)
	}

	// Should fail because photo doesn't exist
	_, err := s.repo.Create(ctx, testPost)
	require.Error(s.T(), err)
}

func (s *PostRepositoryTestSuite) TestCreatePostWithPlantContent() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Create plant
	mainPhotoID := s.uploadTestPhoto(ctx)
	plntSpec, err := plant.NewConiferousSpecification(
		1,
		1,
		10,
		plant.DryMoisture,
		plant.Light,
		plant.HeavySoil,
		9,
	)
	require.NoError(s.T(), err)

	plnt, err := plant.CreatePlant(
		uuid.New(),
		"Testus Plantus",
		"Testus Plantus",
		"Test description",
		mainPhotoID,
		*plant.NewPlantPhotos(),
		plntSpec.Category(),
		plntSpec,
		time.Now(),
		time.Now(),
	)
	require.NoError(s.T(), err)

	_, err = s.plantRepo.Create(ctx, plnt)
	require.NoError(s.T(), err)

	content, err := post.NewContent(
		fmt.Sprintf("Test post content with plant %s", plnt.ID().String()),
		post.ContentTypeWithPlant+"_"+"latex",
	)
	require.NoError(s.T(), err)

	testPost.UpdateContent(*content)

	_, err = s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)
}
