//go:build integration

package plantstorage_test

import (
	"PlantSite/internal/models/plant"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PlantRepositoryTestSuite) TestUpdatePlant() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)

	// Create initial plant
	_, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Upload new photo for update
	newPhotoID := s.pushTestPhoto(ctx)

	// Test update
	updatedPlant, err := s.repo.Update(ctx, testPlant.ID(), func(p *plant.Plant) (*plant.Plant, error) {
		p.UpdateName("Updated Plant")
		p.UpdateLatinName("Updated Latin")
		p.UpdateDescription("Updated description")
		p.UpdateMainPhotoID(newPhotoID)

		// Update photos
		newAdditionalPhotoID := s.pushTestPhoto(ctx)
		photo, err := plant.NewPlantPhoto(newAdditionalPhotoID, "New photo")
		require.NoError(s.T(), err)
		p.AddPhoto(photo)

		return p, nil
	})
	require.NoError(s.T(), err)

	// Verify updates
	assert.Equal(s.T(), "Updated Plant", updatedPlant.GetName())
	assert.Equal(s.T(), "Updated Latin", updatedPlant.GetLatinName())
	assert.Equal(s.T(), "Updated description", updatedPlant.GetDescription())
	assert.Equal(s.T(), newPhotoID, updatedPlant.MainPhotoID())
	assert.Equal(s.T(), 2, updatedPlant.GetPhotos().Len())
	assert.True(s.T(), updatedPlant.UpdatedAt().After(testPlant.UpdatedAt()))

	// Verify new photos exist in storage
	_, err = s.fileRepo.Get(ctx, newPhotoID)
	assert.NoError(s.T(), err)

	updatedPlant.GetPhotos().Iterate(func(photo plant.PlantPhoto) error {
		_, err := s.fileRepo.Get(ctx, photo.FileID())
		assert.NoError(s.T(), err)
		return nil
	})

}

func (s *PlantRepositoryTestSuite) TestUpdatePlantWithInvalidPhoto() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)

	// Create initial plant
	_, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Try to update with invalid photo ID
	invalidPhotoID := uuid.New()

	_, err = s.repo.Update(ctx, testPlant.ID(), func(p *plant.Plant) (*plant.Plant, error) {
		p.UpdateMainPhotoID(invalidPhotoID)
		return p, nil
	})
	require.Error(s.T(), err)
}
