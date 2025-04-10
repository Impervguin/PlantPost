//go:build integration

package albumstorage_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *AlbumRepositoryTestSuite) TestCreateAlbumWithValidPlants() {
	ctx := context.Background()

	// Create test plants first
	plant1 := s.pushTestPlant()
	plant2 := s.pushTestPlant()

	owner := s.pushTestUser()
	plantIDs := uuid.UUIDs{plant1.ID(), plant2.ID()}
	testAlbum := s.createTestAlbum(plantIDs, owner.ID())

	// Test creation
	createdAlbum, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	// Verify plants were associated
	fetchedAlbum, err := s.albumRepo.Get(ctx, createdAlbum.ID())
	require.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), []uuid.UUID{plant1.ID(), plant2.ID()}, fetchedAlbum.PlantIDs())
}

func (s *AlbumRepositoryTestSuite) TestCreateAlbumWithNonExistentPlant() {
	ctx := context.Background()

	owner := s.pushTestUser()
	plantIDs := uuid.UUIDs{uuid.New(), uuid.New()}
	testAlbum := s.createTestAlbum(plantIDs, owner.ID())

	// Test creation - should fail
	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.Error(s.T(), err)

	// Verify no album was created
	_, err = s.albumRepo.Get(ctx, testAlbum.ID())
	require.Error(s.T(), err)
}

func (s *AlbumRepositoryTestSuite) TestCreateAlbumWithSomeInvalidPlants() {
	ctx := context.Background()

	// Create one valid plant
	validPlant := s.pushTestPlant()

	// Create album with mix of valid and invalid plants
	owner := s.pushTestUser()
	plantIDs := uuid.UUIDs{validPlant.ID(), uuid.New()}
	testAlbum := s.createTestAlbum(plantIDs, owner.ID())

	// Test creation - should fail
	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.Error(s.T(), err)

	// Verify no album was created
	_, err = s.albumRepo.Get(ctx, testAlbum.ID())
	require.Error(s.T(), err)
}

func (s *AlbumRepositoryTestSuite) TestCreateAlbumWithEmptyPlants() {
	ctx := context.Background()

	owner := s.pushTestUser()
	plantIDs := uuid.UUIDs{}
	testAlbum := s.createTestAlbum(plantIDs, owner.ID())

	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	alb, err := s.albumRepo.Get(ctx, testAlbum.ID())
	require.NoError(s.T(), err)
	require.Equal(s.T(), len(alb.PlantIDs()), 0)

}
