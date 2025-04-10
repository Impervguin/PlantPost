//go:build integration

package plantstorage_test

import (
	"PlantSite/internal/models/plant"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PlantRepositoryTestSuite) TestCreatePlant() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)

	// Test creation
	createdPlant, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Verify returned plant matches input
	assert.Equal(s.T(), testPlant.ID(), createdPlant.ID())
	assert.Equal(s.T(), testPlant.GetName(), createdPlant.GetName())
	assert.Equal(s.T(), testPlant.GetLatinName(), createdPlant.GetLatinName())
	assert.Equal(s.T(), testPlant.GetDescription(), createdPlant.GetDescription())
	assert.Equal(s.T(), testPlant.MainPhotoID(), createdPlant.MainPhotoID())
	assert.Equal(s.T(), testPlant.GetCategory(), createdPlant.GetCategory())
	assert.Equal(s.T(), testPlant.GetPhotos().Len(), createdPlant.GetPhotos().Len())
}

func (s *PlantRepositoryTestSuite) TestCreatePlantWithInvalidPhoto() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)

	// Delete a photo before creating plant
	photos := testPlant.GetPhotos()
	if photos.Len() > 0 {
		photos.Iterate(func(e plant.PlantPhoto) error {
			return s.fileRepo.DeleteFile(ctx, e.FileID())
		})
	}

	// Should fail because photo doesn't exist
	_, err := s.repo.Create(ctx, testPlant)
	require.Error(s.T(), err)
}

func (s *PlantRepositoryTestSuite) TestCreatePlantWithNoPhotos() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)

	photoIDs := make([]uuid.UUID, 0, testPlant.GetPhotos().Len())

	testPlant.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
		photoIDs = append(photoIDs, e.ID())
		return nil
	})
	for _, photoID := range photoIDs {
		err := testPlant.DeletePhoto(photoID)
		require.NoError(s.T(), err)
	}
	require.NoError(s.T(), testPlant.Validate())

	createdPlant, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, createdPlant.GetPhotos().Len())
}

func (s *PlantRepositoryTestSuite) TestCreateConiferousPlant() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)
	conSpec, err := plant.NewConiferousSpecification(
		1.5,                  // heightM
		0.5,                  // diameterM
		10,                   // soilAcidity
		plant.MediumMoisture, // soilMoisture
		plant.Light,          // lightRelation
		plant.MediumSoil,     // soilType
		10,                   // winterHardiness
	)
	require.NoError(s.T(), err)

	testPlant.UpdateSpec(conSpec)
	require.NoError(s.T(), testPlant.Validate())

	// Test creation
	createdPlant, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Verify returned plant matches input
	assert.Equal(s.T(), testPlant.ID(), createdPlant.ID())
	assert.Equal(s.T(), testPlant.GetName(), createdPlant.GetName())
	assert.Equal(s.T(), testPlant.GetLatinName(), createdPlant.GetLatinName())
	assert.Equal(s.T(), testPlant.GetDescription(), createdPlant.GetDescription())
	assert.Equal(s.T(), testPlant.MainPhotoID(), createdPlant.MainPhotoID())
	assert.Equal(s.T(), testPlant.GetCategory(), createdPlant.GetCategory())
	assert.Equal(s.T(), testPlant.GetPhotos().Len(), createdPlant.GetPhotos().Len())
}

func (s *PlantRepositoryTestSuite) TestCreateDeciduousPlant() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)
	decSpec, err := plant.NewDeciduousSpecification(
		1.5,                  // heightM
		0.5,                  // diameterM
		plant.April,          // floweringPeriod
		10,                   // soilAcidity
		plant.MediumMoisture, // soilMoisture
		plant.Light,          // lightRelation
		plant.MediumSoil,     // soilType
		10,                   // winterHardiness
	)
	require.NoError(s.T(), err)

	testPlant.UpdateSpec(decSpec)
	require.NoError(s.T(), testPlant.Validate())

	// Test creation
	createdPlant, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Verify returned plant matches input
	createdPlant.UpdateSpec(decSpec)
	assert.Equal(s.T(), testPlant.ID(), createdPlant.ID())
	assert.Equal(s.T(), testPlant.GetName(), createdPlant.GetName())
	assert.Equal(s.T(), testPlant.GetLatinName(), createdPlant.GetLatinName())
	assert.Equal(s.T(), testPlant.GetDescription(), createdPlant.GetDescription())
	assert.Equal(s.T(), testPlant.MainPhotoID(), createdPlant.MainPhotoID())
	assert.Equal(s.T(), testPlant.GetCategory(), createdPlant.GetCategory())
	assert.Equal(s.T(), testPlant.GetPhotos().Len(), createdPlant.GetPhotos().Len())
}
