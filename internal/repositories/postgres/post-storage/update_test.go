//go:build integration

package poststorage_test

import (
	"PlantSite/internal/models/post"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *PostRepositoryTestSuite) TestUpdatePost() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Create post first
	_, err := s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)

	// Upload new photo for update
	newPhotoID := s.uploadTestPhoto(ctx)

	// Test update
	_, err = s.repo.Update(ctx, testPost.ID(), func(p *post.Post) (*post.Post, error) {
		p.UpdateTitle("Updated Post")
		content, err := post.NewContent("Updated content", post.ContentTypePlainText)
		if err != nil {
			return nil, err
		}
		p.UpdateContent(*content)
		p.UpdateTags([]string{"updated", "tags"})

		// Update photos

		photo, err := post.CreatePostPhoto(uuid.New(), newPhotoID, 0)
		if err != nil {
			return nil, err
		}
		err = p.AddPhoto(photo)
		if err != nil {
			return nil, err
		}

		return p, nil
	})
	require.NoError(s.T(), err)

	fetchedPost, err := s.repo.Get(ctx, testPost.ID())
	require.NoError(s.T(), err)

	// Verify updates
	assert.Equal(s.T(), "Updated Post", fetchedPost.Title())
	assert.Equal(s.T(), "Updated content", fetchedPost.Content().Text)
	assert.ElementsMatch(s.T(), []string{"updated", "tags"}, fetchedPost.Tags())
	assert.Equal(s.T(), 3, fetchedPost.Photos().Len())
	assert.True(s.T(), fetchedPost.UpdatedAt().After(testPost.UpdatedAt()))
}

func (s *PostRepositoryTestSuite) TestUpdatePostWithInvalidPhoto() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Create post first
	_, err := s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)

	// Try to update with invalid photo ID
	invalidPhotoID := uuid.New()

	_, err = s.repo.Update(ctx, testPost.ID(), func(p *post.Post) (*post.Post, error) {
		photo, err := post.CreatePostPhoto(uuid.New(), invalidPhotoID, 0)
		if err != nil {
			return nil, err
		}
		p.AddPhoto(photo)
		return p, nil
	})
	require.Error(s.T(), err)

	fetchedPost, err := s.repo.Get(ctx, testPost.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 2, fetchedPost.Photos().Len())
}

func (s *PostRepositoryTestSuite) TestUpdatePostRemovePhoto() {
	ctx := context.Background()
	testPost := s.createTestPost(ctx)

	// Create post first
	_, err := s.repo.Create(ctx, testPost)
	require.NoError(s.T(), err)

	_, err = s.repo.Update(ctx, testPost.ID(), func(p *post.Post) (*post.Post, error) {
		p.ClearPhotos()
		return p, nil
	})
	require.NoError(s.T(), err)

	fetchedPost, err := s.repo.Get(ctx, testPost.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 0, fetchedPost.Photos().Len())
}
