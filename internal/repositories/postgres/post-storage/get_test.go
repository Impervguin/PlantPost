//go:build integration

package poststorage_test

import (
	"PlantSite/internal/models/post"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PostRepositoryTestSuite) TestGetPost() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Create post first
	_, err := s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)

	// Test get
	retrievedPost, err := s.repo.Get(ctx, testPost.ID())
	require.NoError(s.T(), err)

	// Verify retrieved post
	assert.Equal(s.T(), testPost.ID(), retrievedPost.ID())
	assert.Equal(s.T(), testPost.Title(), retrievedPost.Title())
	assert.Equal(s.T(), testPost.Content().Text, retrievedPost.Content().Text)
	assert.Equal(s.T(), testPost.AuthorID(), retrievedPost.AuthorID())
	assert.Equal(s.T(), testPost.Photos().Len(), retrievedPost.Photos().Len())
	assert.ElementsMatch(s.T(), testPost.Tags(), retrievedPost.Tags())
}

func (s *PostRepositoryTestSuite) TestGetNonExistentPost() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.repo.Get(ctx, nonExistentID)
	require.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, post.ErrPostNotFound)
}
