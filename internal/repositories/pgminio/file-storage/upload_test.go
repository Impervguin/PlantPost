//go:build integration

package filestorage_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *FileStorageTestSuite) TestUploadFile() {
	ctx := context.Background()
	testData := createTestFileData()

	// Test upload
	file, err := s.storage.UploadFile(ctx, testData)
	require.NoError(s.T(), err)

	// Verify returned file
	assert.NotEqual(s.T(), uuid.Nil, file.ID)
	assert.Equal(s.T(), testData.Name, file.Name)
	assert.Contains(s.T(), file.URL, s.minioClient.GetBucket())
	assert.False(s.T(), file.CreatedAt.IsZero())

	// Verify database record
	dbFile, err := s.storage.Get(ctx, file.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), file.ID, dbFile.ID)
	assert.Equal(s.T(), file.Name, dbFile.Name)
	assert.Equal(s.T(), file.URL, dbFile.URL)
}

func (s *FileStorageTestSuite) TestUploadFileWithEmptyName() {
	ctx := context.Background()
	testData := createTestFileData()
	testData.Name = ""

	_, err := s.storage.UploadFile(ctx, testData)
	require.Error(s.T(), err)
}

func (s *FileStorageTestSuite) TestUploadFileWithNilReader() {
	ctx := context.Background()
	testData := createTestFileData()
	testData.Reader = nil

	_, err := s.storage.UploadFile(ctx, testData)
	require.Error(s.T(), err)
}

func (s *FileStorageTestSuite) TestUploadFileWithEmptyContentType() {
	ctx := context.Background()
	testData := createTestFileData()
	testData.ContentType = ""

	file, err := s.storage.UploadFile(ctx, testData)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), file)
}
