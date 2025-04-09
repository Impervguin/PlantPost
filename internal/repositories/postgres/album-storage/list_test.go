//go:build integration

package albumstorage_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *AlbumRepositoryTestSuite) TestListAlbums() {
	ctx := context.Background()

	plant1 := s.pushTestPlant()
	plant2 := s.pushTestPlant()

	owner := s.pushTestUser()
	plantIDs := uuid.UUIDs{plant1.ID(), plant2.ID()}
	album1 := s.createTestAlbum(plantIDs, owner.ID())
	_, err := s.albumRepo.Create(ctx, album1)
	require.NoError(s.T(), err)

	plantIDs = uuid.UUIDs{plant1.ID()}
	album2 := s.createTestAlbum(plantIDs, owner.ID())
	_, err = s.albumRepo.Create(ctx, album2)
	require.NoError(s.T(), err)

	owner2 := s.pushTestUser()
	otherAlbum := s.createTestAlbum(plantIDs, owner2.ID())
	_, err = s.albumRepo.Create(ctx, otherAlbum)
	require.NoError(s.T(), err)

	albums, err := s.albumRepo.List(ctx, owner.ID())
	require.NoError(s.T(), err)

	assert.Len(s.T(), albums, 2)

	albumIDs := make([]uuid.UUID, 0, 2)
	for _, a := range albums {
		albumIDs = append(albumIDs, a.ID())
		assert.Equal(s.T(), owner.ID(), a.GetOwnerID())
	}

	assert.Contains(s.T(), albumIDs, album1.ID())
	assert.Contains(s.T(), albumIDs, album2.ID())
	assert.NotContains(s.T(), albumIDs, otherAlbum.ID())
}

func (s *AlbumRepositoryTestSuite) TestListAlbumsForNonExistentOwner() {
	ctx := context.Background()
	nonExistentOwnerID := uuid.New()

	albums, err := s.albumRepo.List(ctx, nonExistentOwnerID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), albums)
}
