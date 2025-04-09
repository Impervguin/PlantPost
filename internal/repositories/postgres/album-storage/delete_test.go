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

func (s *AlbumRepositoryTestSuite) TestDeleteAlbum() {
	ctx := context.Background()

	plant1 := s.pushTestPlant()
	plant2 := s.pushTestPlant()

	owner := s.pushTestUser()
	plantIDs := uuid.UUIDs{plant1.ID(), plant2.ID()}
	testAlbum := s.createTestAlbum(plantIDs, owner.ID())

	_, err := s.albumRepo.Create(ctx, testAlbum)
	require.NoError(s.T(), err)

	_, err = s.albumRepo.Get(ctx, testAlbum.ID())
	require.NoError(s.T(), err)

	err = s.albumRepo.Delete(ctx, testAlbum.ID())
	require.NoError(s.T(), err)

	_, err = s.albumRepo.Get(ctx, testAlbum.ID())
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, album.ErrAlbumNotFound))
}

func (s *AlbumRepositoryTestSuite) TestDeleteNonExistentAlbum() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	err := s.albumRepo.Delete(ctx, nonExistentID)
	require.Error(s.T(), err)
}
