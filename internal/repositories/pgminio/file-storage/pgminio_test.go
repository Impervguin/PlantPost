//go:build integration

package filestorage_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	minioclient "PlantSite/internal/infra/minio-client"
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models"
	filestorage "PlantSite/internal/repositories/pgminio/file-storage"
	"PlantSite/internal/repositories/tests"
	"PlantSite/internal/testutils/miniotest"
	"PlantSite/internal/testutils/pgtest"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type FileStorageTestSuite struct {
	suite.Suite
	dbContainer testcontainers.Container
	dbCreds     *pgtest.PostgresCredentials
	fileCnt     testcontainers.Container
	fileCreds   *miniotest.MinioCredentials
	db          *sqpgx.SquirrelPgx
	minioClient *minioclient.MinioClient
	storage     *filestorage.PgMinioStorage
	prevDir     string
}

func TestFileStorageSuite(t *testing.T) {
	suite.RunSuite(t, new(FileStorageTestSuite))
}

func (s *FileStorageTestSuite) BeforeEach(t provider.T) {
	t.Epic("File Storage")
	t.Feature("File Management")
}

func (s *FileStorageTestSuite) BeforeAll(t provider.T) {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(t, err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(t, err)

	pgConfig := pgtest.GetConfig()
	var pgCreds *pgtest.PostgresCredentials
	var container testcontainers.Container
	if pgConfig.External {
		pgCreds, err = pgtest.NewTestExternalPostgres(ctx, pgConfig)
		require.NoError(t, err)
	} else {
		container, pgCreds, err = pgtest.NewTestPostgres(ctx)
		require.NoError(t, err)
	}
	s.dbCreds = pgCreds
	s.dbContainer = container

	if pgConfig.External {
		err = pgtest.CheckMigrationVersion(ctx, pgCreds, pgCreds.Database)
		require.NoError(t, err)
	} else {
		err = pgtest.Migrate(ctx, pgCreds, pgCreds.Database)
		require.NoError(t, err)
	}

	// Create database connection
	dbConfig := &sqpgx.SqpgxConfig{
		User:                   pgCreds.User,
		Password:               pgCreds.Password,
		DBName:                 pgCreds.Database,
		Host:                   pgCreds.Host,
		Port:                   pgCreds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}
	s.db, err = sqpgx.NewSquirrelPgx(ctx, dbConfig)
	require.NoError(t, err)

	var fileCnt testcontainers.Container
	var fileCreds *miniotest.MinioCredentials
	minConfig := miniotest.GetConfig()
	if minConfig.External {
		fileCreds, err = miniotest.NewTestExternalMinio(ctx, minConfig)
		require.NoError(t, err)
	} else {
		fileCnt, fileCreds, err = miniotest.NewTestMinio(ctx)
		require.NoError(t, err)
	}
	s.fileCnt = fileCnt
	s.fileCreds = fileCreds

	err = miniotest.CheckBucketExists(context.Background(), fileCreds)
	if err != nil {
		require.ErrorIs(t, err, miniotest.ErrBucketDoesNotExist)
		err = miniotest.Migrate(context.Background(), fileCreds)
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	minioConfig, err := minioclient.NewMinioConfig(
		fileCreds.GetEndpoint(),
		fileCreds.User,
		fileCreds.Password,
		fileCreds.Bucket,
	)
	require.NoError(t, err)

	minioCl, err := minioclient.NewMinioClient(
		minioConfig,
	)
	require.NoError(t, err)
	s.minioClient = minioCl
	s.storage, err = filestorage.NewPgMinioStorage(ctx, s.db, minioCl)
	require.NoError(t, err)
}

func (s *FileStorageTestSuite) AfterAll(t provider.T) {
	ctx := context.Background()
	if s.fileCnt != nil {
		s.fileCnt.Terminate(ctx)
	}
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	// Restore original directory
	os.Chdir(s.prevDir)
}

func (s *FileStorageTestSuite) AfterEach(t provider.T) {
	// err := pgtest.TruncateTables(context.Background(), s.dbCreds)
	// require.NoError(t, err)
	// err = miniotest.CleanUpBucket(context.Background(), s.fileCreds)
	// require.NoError(t, err)
}

func createTestFileData() models.FileData {
	id := uuid.New()
	return models.FileData{
		Name:        id.String() + ".jpeg",
		Reader:      bytes.NewReader([]byte("test content")),
		ContentType: "image/jpeg",
	}
}

func (s *FileStorageTestSuite) TestUpload(t provider.T) {
	t.Tags("upload", "file")
	t.Description("Test file upload functionality")

	ctx := context.Background()
	testData := createTestFileData()

	t.WithNewStep("Upload file", func(pctx provider.StepCtx) {
		file, err := s.storage.Upload(ctx, &testData)
		require.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, file.ID)
		assert.Equal(t, testData.Name, file.Name)
		assert.Contains(t, file.URL, s.minioClient.GetBucket())
		assert.False(t, file.CreatedAt.IsZero())
	})

	t.WithNewStep("Verify database record", func(pctx provider.StepCtx) {
		// Re-upload to get fresh file
		file, err := s.storage.Upload(ctx, &testData)
		require.NoError(t, err)

		dbFile, err := s.storage.Get(ctx, file.ID)
		require.NoError(t, err)
		assert.Equal(t, file.ID, dbFile.ID)
		assert.Equal(t, file.Name, dbFile.Name)
		assert.Equal(t, file.URL, dbFile.URL)
	})
}

func (s *FileStorageTestSuite) TestGetFile(t provider.T) {
	t.Tags("get", "file")
	t.Description("Test file retrieval functionality")

	ctx := context.Background()
	testData := createTestFileData()

	t.WithNewStep("Upload test file", func(pctx provider.StepCtx) {
		uploadedFile, err := s.storage.Upload(ctx, &testData)
		require.NoError(t, err)

		t.WithNewStep("Get file", func(pctx provider.StepCtx) {
			retrievedFile, err := s.storage.Get(ctx, uploadedFile.ID)
			require.NoError(t, err)

			assert.Equal(t, uploadedFile.ID, retrievedFile.ID)
			assert.Equal(t, uploadedFile.Name, retrievedFile.Name)
			assert.Equal(t, uploadedFile.URL, retrievedFile.URL)
			assert.Equal(t, uploadedFile.CreatedAt.Unix(), retrievedFile.CreatedAt.Unix())
		})
	})
}

func (s *FileStorageTestSuite) TestDownloadFile(t provider.T) {
	t.Tags("download", "file")
	t.Description("Test file download functionality")

	ctx := context.Background()
	testData := createTestFileData()

	t.WithNewStep("Upload test file", func(pctx provider.StepCtx) {
		uploadedFile, err := s.storage.Upload(ctx, &testData)
		require.NoError(t, err)

		t.WithNewStep("Download file", func(pctx provider.StepCtx) {
			downloadedData, err := s.storage.Download(ctx, uploadedFile.ID)
			require.NoError(t, err)

			assert.Equal(t, testData.Name, downloadedData.Name)
			assert.Equal(t, testData.ContentType, downloadedData.ContentType)

			content, err := io.ReadAll(downloadedData.Reader)
			require.NoError(t, err)
			assert.Equal(t, "test content", string(content))
		})
	})
}

func (s *FileStorageTestSuite) TestUpdateFile(t provider.T) {
	t.Tags("update", "file")
	t.Description("Test file update functionality")

	ctx := context.Background()
	testData := createTestFileData()

	t.WithNewStep("Upload test file", func(pctx provider.StepCtx) {
		uploadedFile, err := s.storage.Upload(ctx, &testData)
		require.NoError(t, err)

		t.WithNewStep("Update file", func(pctx provider.StepCtx) {
			updateData := models.FileData{
				Name:        "updated.txt",
				Reader:      bytes.NewReader([]byte("updated content")),
				ContentType: "text/plain",
			}

			updatedFile, err := s.storage.Update(ctx, uploadedFile.ID, &updateData)
			require.NoError(t, err)

			assert.Equal(t, uploadedFile.ID, updatedFile.ID)
			assert.Equal(t, updateData.Name, updatedFile.Name)
			assert.NotEqual(t, uploadedFile.CreatedAt, updatedFile.CreatedAt)
		})

		t.WithNewStep("Verify update in database", func(pctx provider.StepCtx) {
			updateData := models.FileData{
				Name:        "updated_db_check.txt",
				Reader:      bytes.NewReader([]byte("updated content")),
				ContentType: "text/plain",
			}

			updatedFile, err := s.storage.Update(ctx, uploadedFile.ID, &updateData)
			require.NoError(t, err)

			dbFile, err := s.storage.Get(ctx, updatedFile.ID)
			require.NoError(t, err)
			assert.Equal(t, updateData.Name, dbFile.Name)
		})
	})
}

func (s *FileStorageTestSuite) TestDeleteFile(t provider.T) {
	t.Tags("delete", "file")
	t.Description("Test file deletion functionality")

	ctx := context.Background()
	testData := createTestFileData()

	t.WithNewStep("Upload test file", func(pctx provider.StepCtx) {
		uploadedFile, err := s.storage.Upload(ctx, &testData)
		require.NoError(t, err)

		t.WithNewStep("Delete file", func(pctx provider.StepCtx) {
			err = s.storage.Delete(ctx, uploadedFile.ID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify deletion", func(pctx provider.StepCtx) {
			_, err = s.storage.Get(ctx, uploadedFile.ID)
			assert.Error(t, err)
			assert.Equal(t, models.ErrFileNotFound, err)
		})
	})
}

func (s *FileStorageTestSuite) TestGetNonExistentFile(t provider.T) {
	t.Tags("get", "not_found")
	t.Description("Test non-existent file retrieval handling")

	ctx := context.Background()
	nonExistentID := uuid.New()

	t.WithNewStep("Get non-existent file", func(pctx provider.StepCtx) {
		_, err := s.storage.Get(ctx, nonExistentID)
		require.Error(t, err)
		assert.Equal(t, models.ErrFileNotFound, err)
	})
}

func (s *FileStorageTestSuite) TestDownloadNonExistentFile(t provider.T) {
	t.Tags("download", "not_found")
	t.Description("Test non-existent file download handling")

	ctx := context.Background()
	nonExistentID := uuid.New()

	t.WithNewStep("Download non-existent file", func(pctx provider.StepCtx) {
		_, err := s.storage.Download(ctx, nonExistentID)
		require.Error(t, err)
		assert.Equal(t, models.ErrFileNotFound, err)
	})
}

func (s *FileStorageTestSuite) TestUpdateNonExistentFile(t provider.T) {
	t.Tags("update", "not_found")
	t.Description("Test non-existent file update handling")

	ctx := context.Background()
	nonExistentID := uuid.New()
	updateData := createTestFileData()

	t.WithNewStep("Update non-existent file", func(pctx provider.StepCtx) {
		_, err := s.storage.Update(ctx, nonExistentID, &updateData)
		require.Error(t, err)
		assert.Equal(t, models.ErrFileNotFound, err)
	})
}

func (s *FileStorageTestSuite) TestDeleteNonExistentFile(t provider.T) {
	t.Tags("delete", "not_found")
	t.Description("Test non-existent file deletion handling")

	ctx := context.Background()
	nonExistentID := uuid.New()

	t.WithNewStep("Delete non-existent file", func(pctx provider.StepCtx) {
		err := s.storage.Delete(ctx, nonExistentID)
		require.Error(t, err)
	})
}

func (s *FileStorageTestSuite) TestUploadWithEmptyName(t provider.T) {
	t.Tags("upload", "validation")
	t.Description("Test file upload with empty name validation")

	ctx := context.Background()
	testData := createTestFileData()
	testData.Name = ""

	t.WithNewStep("Upload file with empty name", func(pctx provider.StepCtx) {
		_, err := s.storage.Upload(ctx, &testData)
		require.Error(t, err)
	})
}

func (s *FileStorageTestSuite) TestUploadWithNilReader(t provider.T) {
	t.Tags("upload", "validation")
	t.Description("Test file upload with nil reader validation")

	ctx := context.Background()
	testData := createTestFileData()
	testData.Reader = nil

	t.WithNewStep("Upload file with nil reader", func(pctx provider.StepCtx) {
		_, err := s.storage.Upload(ctx, &testData)
		require.Error(t, err)
	})
}
