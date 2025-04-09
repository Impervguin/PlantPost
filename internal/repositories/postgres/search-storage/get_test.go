//go:build integration

package searchstorage_test

import (
	"PlantSite/internal/models/plant"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *SearchRepositoryTestSuite) TestGetPostByID() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Create post first
	_, err := s.postRepo.Create(ctx, testPost)
	require.NoError(s.T(), err)

	// Test get
	retrievedPost, err := s.searchRepo.GetPostByID(ctx, testPost.ID())
	require.NoError(s.T(), err)

	// Verify retrieved post
	assert.Equal(s.T(), testPost.ID(), retrievedPost.ID())
	assert.Equal(s.T(), testPost.Title(), retrievedPost.Title())
	assert.Equal(s.T(), testPost.Content().Text, retrievedPost.Content().Text)
	assert.Equal(s.T(), testPost.AuthorID(), retrievedPost.AuthorID())
	assert.Equal(s.T(), testPost.Photos().Len(), retrievedPost.Photos().Len())
	assert.ElementsMatch(s.T(), testPost.Tags(), retrievedPost.Tags())
}

func (s *SearchRepositoryTestSuite) TestGetPlantByID() {
	ctx := context.Background()
	testPlant := s.createConiferousPlant(ctx, "Test Plant", 1.0, 0.5, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))

	// Create plant first
	_, err := s.plantRepo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Test get
	retrievedPlant, err := s.searchRepo.GetPlantByID(ctx, testPlant.ID())
	require.NoError(s.T(), err)

	// Verify retrieved plant
	assert.Equal(s.T(), testPlant.ID(), retrievedPlant.ID())
	assert.Equal(s.T(), testPlant.GetName(), retrievedPlant.GetName())
	assert.Equal(s.T(), testPlant.GetLatinName(), retrievedPlant.GetLatinName())
	assert.Equal(s.T(), testPlant.GetDescription(), retrievedPlant.GetDescription())
	assert.Equal(s.T(), testPlant.MainPhotoID(), retrievedPlant.MainPhotoID())
	assert.Equal(s.T(), testPlant.GetCategory(), retrievedPlant.GetCategory())
	assert.Equal(s.T(), testPlant.GetPhotos().Len(), retrievedPlant.GetPhotos().Len())
}

func (s *SearchRepositoryTestSuite) TestGetNonExistentPost() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.searchRepo.GetPostByID(ctx, nonExistentID)
	require.Error(s.T(), err)
}

func (s *SearchRepositoryTestSuite) TestGetNonExistentPlant() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.searchRepo.GetPlantByID(ctx, nonExistentID)
	require.Error(s.T(), err)
}
