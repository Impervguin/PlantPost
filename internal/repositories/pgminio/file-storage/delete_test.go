//go:build integration

package filestorage_test

import (
	"context"
	"errors"

	"PlantSite/internal/models"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *FileStorageTestSuite) TestDelete() {
	ctx := context.Background()
	testData := createTestFileData()

	// Upload test file first
	uploadedFile, err := s.storage.Upload(ctx, &testData)
	require.NoError(s.T(), err)

	// Test delete
	err = s.storage.Delete(ctx, uploadedFile.ID)
	require.NoError(s.T(), err)

	// Verify deleted from database
	_, err = s.storage.Get(ctx, uploadedFile.ID)
	assert.True(s.T(), errors.Is(err, models.ErrFileNotFound))

	// Verify deleted from MinIO
	_, err = s.minioClient.StatObject(ctx, s.minioClient.GetBucket(), uploadedFile.URL, minio.StatObjectOptions{})
	assert.Error(s.T(), err)
}

func (s *FileStorageTestSuite) TestDeleteNonExistentFile() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	err := s.storage.Delete(ctx, nonExistentID)
	require.Error(s.T(), err)
}

func (s *FileStorageTestSuite) TestDeleteMissingInMinio() {
	ctx := context.Background()
	testData := createTestFileData()

	// Upload test file first
	uploadedFile, err := s.storage.Upload(ctx, &testData)
	require.NoError(s.T(), err)

	// Delete from MinIO first
	err = s.minioClient.RemoveObject(ctx, s.minioClient.GetBucket(), uploadedFile.URL, minio.RemoveObjectOptions{})
	require.NoError(s.T(), err)

	// Test delete should still succeed for DB cleanup
	err = s.storage.Delete(ctx, uploadedFile.ID)
	require.NoError(s.T(), err)

	// Verify deleted from database and MinIO
	_, err = s.storage.Get(ctx, uploadedFile.ID)
	assert.True(s.T(), errors.Is(err, models.ErrFileNotFound))
}
