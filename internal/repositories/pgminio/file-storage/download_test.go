//go:build integration

package filestorage_test

import (
	"context"
	"errors"
	"io"

	"PlantSite/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *FileStorageTestSuite) TestDownloadFile() {
	ctx := context.Background()
	testData := createTestFileData()

	// Upload test file first
	uploadedFile, err := s.storage.UploadFile(ctx, testData)
	require.NoError(s.T(), err)

	// Test download
	downloadedData, err := s.storage.Download(ctx, uploadedFile.ID)
	require.NoError(s.T(), err)

	// Verify downloaded data
	assert.Equal(s.T(), testData.Name, downloadedData.Name)
	assert.Equal(s.T(), testData.ContentType, downloadedData.ContentType)

	content, err := io.ReadAll(downloadedData.Reader)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "test content", string(content))
}

func (s *FileStorageTestSuite) TestDownloadNonExistentFile() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.storage.Download(ctx, nonExistentID)
	require.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, models.ErrFileNotFound))
}
