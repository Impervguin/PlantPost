//go:build integration

package filestorage_test

import (
	"context"

	"PlantSite/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *FileStorageTestSuite) TestGetFile() {
	ctx := context.Background()
	testData := createTestFileData()


	// Upload test file first
	uploadedFile, err := s.storage.Upload(ctx, &testData)
	require.NoError(s.T(), err)

	// Test get
	retrievedFile, err := s.storage.Get(ctx, uploadedFile.ID)
	require.NoError(s.T(), err)

	// Verify retrieved file
	assert.Equal(s.T(), uploadedFile.ID, retrievedFile.ID)
	assert.Equal(s.T(), uploadedFile.Name, retrievedFile.Name)
	assert.Equal(s.T(), uploadedFile.URL, retrievedFile.URL)
	assert.Equal(s.T(), uploadedFile.CreatedAt.Unix(), retrievedFile.CreatedAt.Unix())
}

func (s *FileStorageTestSuite) TestGetNonExistentFile() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.storage.Get(ctx, nonExistentID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), models.ErrFileNotFound, err)
}

func (s *FileStorageTestSuite) TestGetFileWithInvalidUUID() {
	ctx := context.Background()

	// This test would need to be adjusted if your implementation accepts string IDs
	// Currently testing with zero UUID since we can't create an invalid UUID type
	zeroUUID := uuid.Nil

	_, err := s.storage.Get(ctx, zeroUUID)
	require.Error(s.T(), err)
}
