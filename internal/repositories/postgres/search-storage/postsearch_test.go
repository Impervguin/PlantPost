//go:build integration

package searchstorage_test

import (
	"context"

	"PlantSite/internal/models/search"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *SearchRepositoryTestSuite) TestSearchPostsByTitle() {
	ctx := context.Background()

	// Create test posts
	post1 := s.createTestPost(ctx)
	post1.UpdateTitle("Unique Title 1")

	post2 := s.createAuthorPost(ctx, post1.AuthorID())
	post2.UpdateTitle("Unique Title 2")

	_, err := s.postRepo.Create(ctx, post1)
	require.NoError(s.T(), err)
	_, err = s.postRepo.Create(ctx, post2)
	require.NoError(s.T(), err)

	// Create search with title filter
	srch := search.NewPostSearch()
	srch.AddFilter(search.NewPostTitleContainsFilter("Unique Title 1"))

	// Test search
	posts, err := s.searchRepo.SearchPosts(ctx, srch)
	require.NoError(s.T(), err)

	// Verify results
	assert.Len(s.T(), posts, 1)
	assert.Equal(s.T(), "Unique Title 1", posts[0].Title())
}

func (s *SearchRepositoryTestSuite) TestSearchPostsByTitleContains() {
	ctx := context.Background()

	// Create test posts
	post1 := s.createTestPost(ctx)
	post1.UpdateTitle("Gardening Tips")

	post2 := s.createTestPost(ctx)
	post2.UpdateTitle("Plant Care Guide")

	_, err := s.postRepo.Create(ctx, post1)
	require.NoError(s.T(), err)
	_, err = s.postRepo.Create(ctx, post2)
	require.NoError(s.T(), err)

	// Create search with title contains filter
	srch := search.NewPostSearch()
	srch.AddFilter(search.NewPostTitleContainsFilter("Gardening"))

	// Test search
	posts, err := s.searchRepo.SearchPosts(ctx, srch)
	require.NoError(s.T(), err)

	// Verify results
	assert.Len(s.T(), posts, 1)
	assert.Equal(s.T(), "Gardening Tips", posts[0].Title())
}

func (s *SearchRepositoryTestSuite) TestSearchPostsByTag() {
	ctx := context.Background()

	// Create test posts
	post1 := s.createTestPost(ctx)
	post1.UpdateTags([]string{"gardening", "tips"})

	post2 := s.createTestPost(ctx)
	post2.UpdateTags([]string{"plants", "care"})

	_, err := s.postRepo.Create(ctx, post1)
	require.NoError(s.T(), err)
	_, err = s.postRepo.Create(ctx, post2)
	require.NoError(s.T(), err)

	// Create search with tag filter
	srch := search.NewPostSearch()
	srch.AddFilter(search.NewPostTagFilter([]string{"gardening"}))

	// Test search
	posts, err := s.searchRepo.SearchPosts(ctx, srch)
	require.NoError(s.T(), err)

	// Verify results
	assert.Len(s.T(), posts, 1)
	assert.Contains(s.T(), posts[0].Tags(), "gardening")
}

func (s *SearchRepositoryTestSuite) TestSearchPostsByAuthor() {
	ctx := context.Background()

	// Create test posts
	post1 := s.createTestPost(ctx)

	post2 := s.createAuthorPost(ctx, post1.AuthorID())

	post3 := s.createTestPost(ctx)

	_, err := s.postRepo.Create(ctx, post1)
	require.NoError(s.T(), err)
	_, err = s.postRepo.Create(ctx, post2)
	require.NoError(s.T(), err)
	_, err = s.postRepo.Create(ctx, post3)
	require.NoError(s.T(), err)

	// Create search with author filter
	srch := search.NewPostSearch()
	srch.AddFilter(search.NewPostAuthorFilter(post1.AuthorID()))

	// Test search
	posts, err := s.searchRepo.SearchPosts(ctx, srch)
	require.NoError(s.T(), err)

	// Verify results
	assert.Len(s.T(), posts, 2)
}

func (s *SearchRepositoryTestSuite) TestSearchPostsWithMultipleFilters() {
	ctx := context.Background()

	// Create test posts
	post1 := s.createTestPost(ctx)
	post1.UpdateTitle("Gardening Tips")
	post1.UpdateTags([]string{"gardening", "tips"})

	post2 := s.createTestPost(ctx)
	post2.UpdateTitle("Plant Care")
	post2.UpdateTags([]string{"plants", "tips"})

	post3 := s.createTestPost(ctx)
	post3.UpdateTitle("Gardening Tips 2")
	post3.UpdateTags([]string{"garden", "garden_tips"})

	_, err := s.postRepo.Create(ctx, post1)
	require.NoError(s.T(), err)
	_, err = s.postRepo.Create(ctx, post2)
	require.NoError(s.T(), err)
	_, err = s.postRepo.Create(ctx, post3)
	require.NoError(s.T(), err)

	// Create search with multiple filters
	srch := search.NewPostSearch()
	srch.AddFilter(search.NewPostTitleContainsFilter("Gardening"))
	srch.AddFilter(search.NewPostTagFilter([]string{"tips"}))

	// Test search
	posts, err := s.searchRepo.SearchPosts(ctx, srch)
	require.NoError(s.T(), err)

	// Verify results
	assert.Len(s.T(), posts, 1)
	assert.Equal(s.T(), "Gardening Tips", posts[0].Title())
	assert.Contains(s.T(), posts[0].Tags(), "tips")
}
