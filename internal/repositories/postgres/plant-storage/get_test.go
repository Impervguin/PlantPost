//go:build integration

package plantstorage_test

import (
	"PlantSite/internal/models/plant"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PlantRepositoryTestSuite) TestGetPlant() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)

	// Create plant first
	_, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Test get
	retrievedPlant, err := s.repo.Get(ctx, testPlant.ID())
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

func (s *PlantRepositoryTestSuite) TestGetNonExistentPlant() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.repo.Get(ctx, nonExistentID)
	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, plant.ErrPlantNotFound)
}
