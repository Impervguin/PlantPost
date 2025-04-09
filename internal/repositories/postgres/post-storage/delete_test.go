//go:build integration

package poststorage_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PostRepositoryTestSuite) TestDeletePost() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Create post first
	_, err := s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)

	// Verify exists
	_, err = s.repo.Get(ctx, testPost.ID())
	require.NoError(s.T(), err)

	// Verify photos exist
	for _, photo := range testPost.Photos().List() {
		_, err := s.fileRepo.Get(ctx, photo.FileID())
		require.NoError(s.T(), err)
	}

	// Test deletion
	err = s.repo.Delete(ctx, testPost.ID())
	require.NoError(s.T(), err)

	// Verify post deleted
	_, err = s.repo.Get(ctx, testPost.ID())
	require.Error(s.T(), err)

	// Verify photos still exist in storage (file storage cleanup should be handled separately)
	for _, photo := range testPost.Photos().List() {
		_, err := s.fileRepo.Get(ctx, photo.FileID())
		assert.NoError(s.T(), err)
	}
}

func (s *PostRepositoryTestSuite) TestDeleteNonExistentPost() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	err := s.repo.Delete(ctx, nonExistentID)
	require.Error(s.T(), err)
}
