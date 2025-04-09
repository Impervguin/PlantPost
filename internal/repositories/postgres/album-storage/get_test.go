//go:build integration

package albumstorage_test

import (
	"context"
	"errors"

	"PlantSite/internal/models/album"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *AlbumRepositoryTestSuite) TestGetAlbum() {
	ctx := context.Background()

	plant1 := s.pushTestPlant()
	plant2 := s.pushTestPlant()

	owner := s.pushTestUser()
	plantIDs := uuid.UUIDs{plant1.ID(), plant2.ID()}
	testAlbum := s.createTestAlbum(plantIDs, owner.ID())

	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	fetchedAlbum, err := s.albumRepo.Get(ctx, testAlbum.ID())
	require.NoError(s.T(), err)

	assert.Equal(s.T(), testAlbum.ID(), fetchedAlbum.ID())
	assert.Equal(s.T(), testAlbum.Name(), fetchedAlbum.Name())
	assert.Equal(s.T(), testAlbum.Description(), fetchedAlbum.Description())
	assert.Equal(s.T(), testAlbum.GetOwnerID(), fetchedAlbum.GetOwnerID())
	assert.ElementsMatch(s.T(), testAlbum.PlantIDs(), fetchedAlbum.PlantIDs())
}

func (s *AlbumRepositoryTestSuite) TestGetNonExistentAlbum() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.albumRepo.Get(ctx, nonExistentID)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, album.ErrAlbumNotFound))
}
