//go:build integration

package plantstorage_test

import (
	"PlantSite/internal/models/plant"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PlantRepositoryTestSuite) TestDeletePlant() {
	ctx := context.Background()
	testPlant := s.createTestPlant(ctx)

	// Create plant first
	_, err := s.repo.Create(ctx, testPlant)
	require.NoError(s.T(), err)

	// Verify exists
	_, err = s.repo.Get(ctx, testPlant.ID())
	require.NoError(s.T(), err)

	// Verify photos exist
	_, err = s.fileRepo.Get(ctx, testPlant.MainPhotoID())
	require.NoError(s.T(), err)

	testPlant.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
		_, err := s.fileRepo.Get(ctx, e.FileID())
		require.NoError(s.T(), err)
		return nil
	})

	// Test deletion
	err = s.repo.Delete(ctx, testPlant.ID())
	require.NoError(s.T(), err)

	// Verify plant deleted
	_, err = s.repo.Get(ctx, testPlant.ID())
	require.Error(s.T(), err)

	// Verify photos still exist in storage (file storage cleanup should be handled separately)
	_, err = s.fileRepo.Get(ctx, testPlant.MainPhotoID())
	assert.NoError(s.T(), err)
}

func (s *PlantRepositoryTestSuite) TestDeleteNonExistentPlant() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	err := s.repo.Delete(ctx, nonExistentID)
	require.Error(s.T(), err)
}
