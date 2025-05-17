//go:build integration

package filestorage_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *FileStorageTestSuite) TestUpload() {
	ctx := context.Background()
	testData := createTestFileData()

	// Test upload
	file, err := s.storage.Upload(ctx, &testData)
	require.NoError(s.T(), err)

	// Verify returned file
	assert.NotEqual(s.T(), uuid.Nil, file.ID)
	assert.Equal(s.T(), testData.Name, file.Name)
	assert.Contains(s.T(), file.URL, s.testBucketName)
	assert.False(s.T(), file.CreatedAt.IsZero())

	// Verify database record
	dbFile, err := s.storage.Get(ctx, file.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), file.ID, dbFile.ID)
	assert.Equal(s.T(), file.Name, dbFile.Name)
	assert.Equal(s.T(), file.URL, dbFile.URL)
}

func (s *FileStorageTestSuite) TestUploadWithEmptyName() {
	ctx := context.Background()
	testData := createTestFileData()
	testData.Name = ""

	_, err := s.storage.Upload(ctx, &testData)
	require.Error(s.T(), err)
}

func (s *FileStorageTestSuite) TestUploadWithNilReader() {
	ctx := context.Background()
	testData := createTestFileData()
	testData.Reader = nil

	_, err := s.storage.Upload(ctx, &testData)
	require.Error(s.T(), err)
}

func (s *FileStorageTestSuite) TestUploadWithEmptyContentType() {
	ctx := context.Background()
	testData := createTestFileData()
	testData.ContentType = ""

	file, err := s.storage.Upload(ctx, &testData)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), file)
}
