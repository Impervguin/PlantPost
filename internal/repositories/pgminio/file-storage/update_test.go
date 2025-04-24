//go:build integration

package filestorage_test

import (
	"bytes"
	"context"
	"io"

	"PlantSite/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *FileStorageTestSuite) TestUpdateFile() {
	ctx := context.Background()
	testData := createTestFileData()

	// Upload test file first
	uploadedFile, err := s.storage.Upload(ctx, &testData)
	require.NoError(s.T(), err)

	// Prepare update data
	updateData := models.FileData{
		Name:        "updated.txt",
		Reader:      bytes.NewReader([]byte("updated content")),
		ContentType: "text/plain",
	}

	// Test update
	updatedFile, err := s.storage.Update(ctx, uploadedFile.ID, &updateData)
	require.NoError(s.T(), err)

	// Verify updated fields
	assert.Equal(s.T(), uploadedFile.ID, updatedFile.ID)
	assert.Equal(s.T(), updateData.Name, updatedFile.Name)
	assert.NotEqual(s.T(), uploadedFile.CreatedAt, updatedFile.CreatedAt)

	// Verify database record
	dbFile, err := s.storage.Get(ctx, updatedFile.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), updateData.Name, dbFile.Name)

	// Verify MinIO object was updated
	downloadedData, err := s.storage.Download(ctx, updatedFile.ID)
	require.NoError(s.T(), err)
	content, err := io.ReadAll(downloadedData.Reader)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "updated content", string(content))
}

func (s *FileStorageTestSuite) TestUpdateNonExistentFile() {
	ctx := context.Background()
	nonExistentID := uuid.New()
	updateData := createTestFileData()

	_, err := s.storage.Update(ctx, nonExistentID, &updateData)
	require.Error(s.T(), err)
	assert.Equal(s.T(), models.ErrFileNotFound, err)
}

func (s *FileStorageTestSuite) TestUpdateFileWithPartialData() {
	ctx := context.Background()
	testData := createTestFileData()

	// Upload test file first
	uploadedFile, err := s.storage.Upload(ctx, &testData)
	require.NoError(s.T(), err)

	// Prepare update data with only name changed
	updateData := models.FileData{
		Name: "updated.txt",
		// Reader and ContentType not set
	}

	// Test update
	updatedFile, err := s.storage.Update(ctx, uploadedFile.ID, &updateData)
	require.NoError(s.T(), err)

	// Verify name was updated
	assert.Equal(s.T(), updateData.Name, updatedFile.Name)

	// Verify content wasn't changed
	downloadedData, err := s.storage.Download(ctx, updatedFile.ID)
	require.NoError(s.T(), err)
	content, err := io.ReadAll(downloadedData.Reader)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "test content", string(content))
}
