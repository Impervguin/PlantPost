//go:build integration

package filestorage_test

import (
	"bytes"
	"context"
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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type FileStorageTestSuite struct {
	suite.Suite
	dbContainer    testcontainers.Container
	minioContainer testcontainers.Container
	db             *sqpgx.SquirrelPgx
	minioClient    *minioclient.MinioClient
	storage        *filestorage.PgMinioStorage
	prevDir        string
}

func TestFileStorageSuite(t *testing.T) {
	suite.Run(t, new(FileStorageTestSuite))
}

func (s *FileStorageTestSuite) SetupSuite() {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(s.T(), err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(s.T(), err)

	// Setup PostgreSQL container
	dbContainer, dbCreds, err := pgtest.NewTestPostgres(ctx)
	require.NoError(s.T(), err)
	s.dbContainer = dbContainer

	// Run migrations
	err = pgtest.Migrate(ctx, &dbCreds)
	require.NoError(s.T(), err)

	// Create database connection
	dbConfig := &sqpgx.SqpgxConfig{
		User:                   dbCreds.User,
		Password:               dbCreds.Password,
		DbName:                 dbCreds.Database,
		Host:                   dbCreds.Host,
		Port:                   dbCreds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}
	s.db, err = sqpgx.NewSquirrelPgx(ctx, dbConfig)
	require.NoError(s.T(), err)

	// Setup MinIO container
	minioContainer, minioCreds, err := miniotest.NewTestMinio(ctx)
	require.NoError(s.T(), err)
	s.minioContainer = minioContainer

	// Run migrations
	err = miniotest.Migrate(ctx, minioCreds)
	require.NoError(s.T(), err)

	// Create MinIO client
	minioConfig, err := minioclient.NewMinioConfig(
		minioCreds.GetEndpoint(),
		minioCreds.User,
		minioCreds.Password,
		minioCreds.Bucket,
	)
	require.NoError(s.T(), err)

	s.minioClient, err = minioclient.NewMinioClient(minioConfig)
	require.NoError(s.T(), err)

	// Create storage instance
	s.storage, err = filestorage.NewPgMinioStorage(ctx, s.db, s.minioClient)
	require.NoError(s.T(), err)
}

func (s *FileStorageTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.minioContainer != nil {
		s.minioContainer.Terminate(ctx)
	}
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	// Restore original directory
	os.Chdir(s.prevDir)
}

func createTestFileData() models.FileData {
	id := uuid.New()
	return models.FileData{
		Name:        id.String() + ".jpeg",
		Reader:      bytes.NewReader([]byte("test content")),
		ContentType: "image/jpeg",
	}
}
