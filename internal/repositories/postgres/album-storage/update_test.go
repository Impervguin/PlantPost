//go:build integration

package albumstorage_test

import (
	"PlantSite/internal/models/album"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *AlbumRepositoryTestSuite) TestUpdateAlbumWithValidPlants() {
	ctx := context.Background()

	// Create initial album with one plant
	initialPlant := s.pushTestPlant()
	owner := s.pushTestUser()
	plantIds := []uuid.UUID{initialPlant.ID()}
	testAlbum := s.createTestAlbum(plantIds, owner.ID())
	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	// Create additional plant for update
	newPlant := s.pushTestPlant()

	// Test update with new valid plants
	updatedAlbum, err := s.albumRepo.Update(ctx, testAlbum.ID(), func(a *album.Album) (*album.Album, error) {
		a.AddPlant(newPlant.ID())
		return a, nil
	})
	require.NoError(s.T(), err)

	// Verify update
	assert.ElementsMatch(s.T(), []uuid.UUID{initialPlant.ID(), newPlant.ID()}, updatedAlbum.PlantIDs())

	// Verify persistence
	fetchedAlbum, err := s.albumRepo.Get(ctx, testAlbum.ID())
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), []uuid.UUID{initialPlant.ID(), newPlant.ID()}, fetchedAlbum.PlantIDs())
}

func (s *AlbumRepositoryTestSuite) TestUpdateAlbumWithInvalidPlants() {
	ctx := context.Background()

	// Create initial album with valid plant
	initialPlant := s.pushTestPlant()
	owner := s.pushTestUser()
	plantIds := []uuid.UUID{initialPlant.ID()}
	testAlbum := s.createTestAlbum(plantIds, owner.ID())
	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	// Test update with invalid plant
	_, err = s.albumRepo.Update(ctx, testAlbum.ID(), func(a *album.Album) (*album.Album, error) {
		a.AddPlant(uuid.New()) // Non-existent plant
		return a, nil
	})
	require.Error(s.T(), err)

	// Verify original plants unchanged
	fetchedAlbum, err := s.albumRepo.Get(ctx, testAlbum.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), testAlbum.PlantIDs(), fetchedAlbum.PlantIDs())
}

func (s *AlbumRepositoryTestSuite) TestUpdateAlbumRemovesPlants() {
	ctx := context.Background()

	// Create album with plants
	plant1 := s.pushTestPlant()
	plant2 := s.pushTestPlant()
	owner := s.pushTestUser()
	plantIds := []uuid.UUID{plant1.ID(), plant2.ID()}
	testAlbum := s.createTestAlbum(plantIds, owner.ID())
	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	// Update to remove all plants
	updatedAlbum, err := s.albumRepo.Update(ctx, testAlbum.ID(), func(a *album.Album) (*album.Album, error) {
		a.RemovePlant(plant1.ID())
		a.RemovePlant(plant2.ID())
		return a, nil
	})
	require.NoError(s.T(), err)

	// Verify no plants
	assert.Empty(s.T(), updatedAlbum.PlantIDs())

	// Verify in DB
	fetchedAlbum, err := s.albumRepo.Get(ctx, testAlbum.ID())
	require.NoError(s.T(), err)
	assert.Empty(s.T(), fetchedAlbum.PlantIDs())
}

func (s *AlbumRepositoryTestSuite) TestUpdateDescription() {
	ctx := context.Background()

	// Create album with plants
	plant1 := s.pushTestPlant()
	plant2 := s.pushTestPlant()
	owner := s.pushTestUser()
	plantIds := []uuid.UUID{plant1.ID(), plant2.ID()}
	testAlbum := s.createTestAlbum(plantIds, owner.ID())
	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	// Update to remove all plants
	_, err = s.albumRepo.Update(ctx, testAlbum.ID(), func(a *album.Album) (*album.Album, error) {
		a.UpdateDescription("New Description")
		a.UpdateName("New Name")
		return a, nil
	})
	require.NoError(s.T(), err)

	// Verify in DB
	fetchedAlbum, err := s.albumRepo.Get(ctx, testAlbum.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "New Name", fetchedAlbum.Name())
	assert.Equal(s.T(), "New Description", fetchedAlbum.Description())
	assert.Equal(s.T(), testAlbum.GetOwnerID(), fetchedAlbum.GetOwnerID())
	assert.True(s.T(), testAlbum.UpdatedAt().Before(fetchedAlbum.UpdatedAt()))
}
