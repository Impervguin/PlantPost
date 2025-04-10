//go:build integration

package poststorage_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PostRepositoryTestSuite) TestListAuthorPosts() {
	ctx := context.Background()

	firstPost := s.createTestPost(ctx)

	_, err := s.repo.Create(ctx, firstPost)
	require.NoError(s.T(), err)

	// Create test posts
	for i := 0; i < 2; i++ {
		testPost := s.createAuthorPost(ctx, firstPost.AuthorID())

		_, err := s.repo.Create(ctx, testPost)
		require.NoError(s.T(), err)
	}

	// Create a post from different author
	otherPost := s.createTestPost(ctx)

	_, err = s.repo.Create(ctx, otherPost)
	require.NoError(s.T(), err)

	// Test list
	posts, err := s.repo.ListAuthorPosts(ctx, firstPost.AuthorID())
	require.NoError(s.T(), err)

	// Verify only author's posts returned
	assert.Len(s.T(), posts, 3)
	for _, p := range posts {
		assert.Equal(s.T(), firstPost.AuthorID(), p.AuthorID())
	}
}

func (s *PostRepositoryTestSuite) TestListAuthorPostsWithNoPosts() {
	ctx := context.Background()
	authorID := uuid.New()

	posts, err := s.repo.ListAuthorPosts(ctx, authorID)
	require.NoError(s.T(), err)
	assert.Empty(s.T(), posts)
}
